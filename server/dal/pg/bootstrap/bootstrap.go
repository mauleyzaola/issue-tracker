package bootstrap

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecweb/setup"
	"github.com/rubenv/sql-migrate"
)

type BootstrapDb struct {
	Base       *pg.Db
	priorityDb *database.Priority
	statusDb   *database.Status
	userDb     *database.User
}

func New(db database.Db) *BootstrapDb {
	base := db.(*pg.Db)
	return &BootstrapDb{Base: base}
}

func (db *BootstrapDb) PriorityDb() database.Priority {
	return *db.priorityDb
}

func (db *BootstrapDb) SetPriorityDb(item *database.Priority) {
	db.priorityDb = item
}

func (db *BootstrapDb) StatusDb() database.Status {
	return *db.statusDb
}

func (db *BootstrapDb) SetStatusDb(item *database.Status) {
	db.statusDb = item
}

func (db *BootstrapDb) UserDb() database.User {
	return *db.userDb
}

func (db *BootstrapDb) SetUserDb(item *database.User) {
	db.userDb = item
}

func (t *BootstrapDb) UpgradeDbScripts(conn *setup.ConfigurationDatabase, environment, chdir string) error {
	connStr, err := conn.ConnectionString()
	if err != nil {
		return err
	}

	db, err := sql.Open(conn.Driver, connStr)
	if err != nil {
		return err
	}

	migrations := &migrate.FileMigrationSource{Dir: filepath.Join(chdir, "migrations", conn.Driver)}
	migrate.SetTable("migration_issue_tracker")

	log.Printf("connecting to: %s\t%s\n", conn.Driver, conn.DatabaseName)
	counter, err := migrate.Exec(db, conn.Driver, migrations, migrate.Up)
	if err != nil {
		log.Println("error trying to execute migration:", err)
		return err
	}
	log.Printf("applied %v db scripts in issue-tracker migration\n", counter)

	return t.ResetApplicationPath(chdir)
}

func (db *BootstrapDb) BootstrapAll(tx interface{}, conn *setup.ConfigurationDatabase, environment string, chdir string) error {
	err := db.UpgradeDbScripts(conn, environment, chdir)
	if err != nil {
		return err
	}

	//create default user if it doesn't exist
	if _, _, err = db.BootstrapAdminUser(tx); err != nil {
		return err
	}

	//add missing permission names
	if err = db.CreatePermissionNames(tx); err != nil {
		return err
	}

	//create default workflows
	if err = db.BootstrapWorkflows(tx); err != nil {
		return err
	}

	if err = db.BootstrapPriorities(tx); err != nil {
		return err
	}

	//point to the root directory of the application again
	return db.ResetApplicationPath(chdir)
}

func (t *BootstrapDb) BootstrapAdminUser(tx interface{}) (bool, *domain.User, error) {
	var users []domain.User
	_, err := t.Base.Executor(tx).Select(&users, "select * from users where issystemadministrator=$1 and isactive=$1", true)
	if err != nil {
		return false, nil, err
	}
	if len(users) != 0 {
		return true, &users[0], nil
	}

	admin := &domain.User{}
	admin.DateCreated = time.Now()
	admin.Email = "admin@admin.com"
	admin.IsActive = true
	admin.IsSystemAdministrator = true
	admin.Name = "System"
	admin.LastName = "Administrator"
	admin.Password = "admin"
	err = t.UserDb().Create(tx, admin)

	return false, admin, err
}

func (db *BootstrapDb) ResetApplicationPath(chdir string) error {
	return os.Chdir(chdir)
}

func (t *BootstrapDb) CreatePermissionNames(tx interface{}) error {
	item := &domain.PermissionName{}
	items := item.Permissions()
	for i := range items {
		p := items[i]
		if count, err := t.Base.Executor(tx).SelectInt("select count(*) from permission_name where lower(name)=$1", strings.ToLower(p)); err != nil {
			return err
		} else if count == 0 {
			per := &domain.PermissionName{}
			per.Name = p
			err = t.Base.Executor(tx).Insert(per)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *BootstrapDb) BootstrapWorkflows(tx interface{}) error {
	oldItems, err := t.StatusDb().WorkflowList(tx)
	if err != nil {
		return err
	}

	for i := range oldItems {
		it := &oldItems[i]
		if it.Name == "Approval" || it.Name == "Issue" {
			return nil
		}
	}

	wf := &domain.Workflow{Name: "Approval"}
	err = t.StatusDb().WorkflowCreate(tx, wf)
	if err != nil {
		return err
	}

	opened := &domain.Status{Name: "Open", Workflow: wf}
	approved := &domain.Status{Name: "Approved", Workflow: wf}
	cancelled := &domain.Status{Name: "Cancelled", Workflow: wf}
	err = t.StatusDb().Create(tx, opened)
	if err != nil {
		return err
	}
	err = t.StatusDb().Create(tx, approved)
	if err != nil {
		return err
	}
	err = t.StatusDb().Create(tx, cancelled)
	if err != nil {
		return err
	}

	open := &domain.WorkflowStep{Workflow: wf, Name: "Open"}
	open.NextStatus = opened
	err = t.StatusDb().WorkflowStepCreate(tx, open)
	if err != nil {
		return err
	}

	approve := &domain.WorkflowStep{Workflow: wf, Name: "Approve"}
	approve.PrevStatus = opened
	approve.NextStatus = approved
	approve.Resolves = true
	err = t.StatusDb().WorkflowStepCreate(tx, approve)
	if err != nil {
		return err
	}

	cancel := &domain.WorkflowStep{Workflow: wf, Name: "Cancel"}
	cancel.PrevStatus = opened
	cancel.NextStatus = cancelled
	cancel.Cancels = true
	err = t.StatusDb().WorkflowStepCreate(tx, cancel)
	if err != nil {
		return err
	}

	//default issue workflow
	wt := &domain.Workflow{Name: "Issue"}

	stNew := &domain.Status{Name: "New", Workflow: wt, Description: "The issue has not yet started"}
	stInProgress := &domain.Status{Name: "In Progress", Workflow: wt, Description: "Working has begun on the issue"}
	stResolved := &domain.Status{Name: "Resolved", Workflow: wt, Description: "The issue has been resolved and is waiting to be validated"}
	stClosed := &domain.Status{Name: "Closed", Workflow: wt, Description: "The issue has been successfully validated"}
	stCancelled := &domain.Status{Name: "Cancelled", Workflow: wt, Description: "The issue has been cancelled"}
	stReopened := &domain.Status{Name: "Reopened", Workflow: wt, Description: "The issue was resolved at some point, but was not successfully validated"}

	pNew := &domain.WorkflowStep{Workflow: wt}
	pNew.NextStatus = stNew
	pNew.Name = "Create Issue"

	pBeginProgress := &domain.WorkflowStep{Workflow: wt}
	pBeginProgress.PrevStatus = stNew
	pBeginProgress.NextStatus = stInProgress
	pBeginProgress.Name = "Start Progress"

	pResolve := &domain.WorkflowStep{Workflow: wt}
	pResolve.PrevStatus = stInProgress
	pResolve.NextStatus = stResolved
	pResolve.Name = "Resolve"

	pStopProgress := &domain.WorkflowStep{Workflow: wt}
	pStopProgress.PrevStatus = stInProgress
	pStopProgress.NextStatus = stNew
	pStopProgress.Name = "Stop Progress"

	pClose := &domain.WorkflowStep{Workflow: wt}
	pClose.PrevStatus = stResolved
	pClose.NextStatus = stClosed
	pClose.Resolves = true
	pClose.Name = "Close"

	pReopen := &domain.WorkflowStep{Workflow: wt}
	pReopen.PrevStatus = stResolved
	pReopen.NextStatus = stReopened
	pReopen.Name = "Reopen"

	pResolve2 := &domain.WorkflowStep{Workflow: wt}
	pResolve2.PrevStatus = stReopened
	pResolve2.NextStatus = stResolved
	pResolve2.Name = "Resolve"

	pCancel := &domain.WorkflowStep{Workflow: wt}
	pCancel.PrevStatus = stNew
	pCancel.NextStatus = stCancelled
	pCancel.Cancels = true
	pCancel.Name = "Cancel"

	err = t.StatusDb().WorkflowCreate(tx, wt)
	if err != nil {
		return err
	}

	err = t.StatusDb().Create(tx, stCancelled)
	if err != nil {
		return err
	}

	err = t.StatusDb().Create(tx, stClosed)
	if err != nil {
		return err
	}

	err = t.StatusDb().Create(tx, stInProgress)
	if err != nil {
		return err
	}

	err = t.StatusDb().Create(tx, stNew)
	if err != nil {
		return err
	}

	err = t.StatusDb().Create(tx, stReopened)
	if err != nil {
		return err
	}

	err = t.StatusDb().Create(tx, stResolved)
	if err != nil {
		return err
	}

	err = t.StatusDb().WorkflowStepCreate(tx, pNew)
	if err != nil {
		return err
	}

	err = t.StatusDb().WorkflowStepCreate(tx, pBeginProgress)
	if err != nil {
		return err
	}

	err = t.StatusDb().WorkflowStepCreate(tx, pResolve)
	if err != nil {
		return err
	}

	err = t.StatusDb().WorkflowStepCreate(tx, pStopProgress)
	if err != nil {
		return err
	}

	err = t.StatusDb().WorkflowStepCreate(tx, pClose)
	if err != nil {
		return err
	}

	err = t.StatusDb().WorkflowStepCreate(tx, pReopen)
	if err != nil {
		return err
	}

	err = t.StatusDb().WorkflowStepCreate(tx, pCancel)
	if err != nil {
		return err
	}

	err = t.StatusDb().WorkflowStepCreate(tx, pResolve2)
	if err != nil {
		return err
	}

	return nil
}

func (t *BootstrapDb) BootstrapPriorities(tx interface{}) error {
	createIfNotExist := func(name string) error {
		if count, err := t.Base.Executor(tx).SelectInt("select count(*) from priority where lower(name)=$1", strings.ToLower(name)); err != nil {
			return err
		} else if count == 0 {
			pr := &domain.Priority{Name: name}
			return t.PriorityDb().Create(tx, pr)
		} else {
			return nil
		}
	}
	values := []string{"Blocker", "Critical", "Major", "Minor", "Trivial"}
	for _, v := range values {
		if err := createIfNotExist(v); err != nil {
			return err
		}
	}
	return nil
}
