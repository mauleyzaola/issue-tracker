package issue

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
)

type IssueDb struct {
	Base         *pg.Db
	fileItemDb   *database.FileItem
	permissionDb *database.Permission
	priorityDb   *database.Priority
	projectDb    *database.Project
	statusDb     *database.Status
	userDb       *database.User
}

func New(db database.Db) *IssueDb {
	base := db.(*pg.Db)
	return &IssueDb{Base: base}
}

func (t *IssueDb) FileItemDb() database.FileItem {
	return *t.fileItemDb
}

func (t *IssueDb) SetFileItemDb(item *database.FileItem) {
	t.fileItemDb = item
}

func (t *IssueDb) PermissionDb() database.Permission {
	return *t.permissionDb
}

func (t *IssueDb) SetPermissionDb(item *database.Permission) {
	t.permissionDb = item
}

func (t *IssueDb) PriorityDb() database.Priority {
	return *t.priorityDb
}

func (t *IssueDb) SetPriorityDb(item *database.Priority) {
	t.priorityDb = item
}

func (t *IssueDb) ProjectDb() database.Project {
	return *t.projectDb
}

func (t *IssueDb) SetProjectDb(item *database.Project) {
	t.projectDb = item
}

func (t *IssueDb) StatusDb() database.Status {
	return *t.statusDb
}

func (t *IssueDb) SetStatusDb(item *database.Status) {
	t.statusDb = item
}

func (t *IssueDb) UserDb() database.User {
	return *t.userDb
}

func (t *IssueDb) SetUserDb(item *database.User) {
	t.userDb = item
}

func (t *IssueDb) StatusChange(tx interface{}, issue *domain.Issue, status *domain.Status, fn database.IssueStatusFn) error {
	issue, err := t.Load(tx, issue.Id, issue.Pkey)
	if err != nil {
		return err
	}

	oldStatus, err := t.StatusDb().Load(tx, issue.Status.Id)
	if err != nil {
		return err
	}

	availableSteps, err := t.StatusDb().WorkflowStepAvailableUser(tx, issue.Workflow, issue.Status)
	if err != nil {
		return err
	}
	found := false
	var nextStep *domain.WorkflowStep
	for i := range availableSteps {
		st := &availableSteps[i]
		if st.NextStatus.Id == status.Id {
			found = true
			nextStep = st
			break
		}
	}

	if !found {
		return fmt.Errorf("invalid status change")
	}

	if nextStep.Resolves {
		issue.ResolvedDate = &time.Time{}
		*issue.ResolvedDate = time.Now()
		issue.CancelledDate = nil
	} else if nextStep.Cancels {
		issue.CancelledDate = &time.Time{}
		*issue.CancelledDate = time.Now()
		issue.ResolvedDate = nil
	} else {
		issue.ResolvedDate = nil
		issue.CancelledDate = nil
	}

	if nextStep.Resolves || nextStep.Cancels {
		children, err := t.Children(tx, issue)
		if err != nil {
			return err
		}
		for i := range children {
			child := &children[i]
			if child.ResolvedDate == nil && child.CancelledDate == nil {
				return fmt.Errorf("cannot resolve or cancel the issue until all their subtasks are resolved or cancelled")
			}
		}
	}

	issue.IdStatus = nextStep.NextStatus.Id

	_, err = t.Base.Executor(tx).Update(issue)
	if err != nil {
		return err
	}

	if fn != nil {
		if err = fn(tx, issue, nextStep, oldStatus, nextStep.NextStatus); err != nil {
			return err
		}
	}

	//TODO: generate notifications to subscribers
	return nil
}

func (t *IssueDb) Create(tx interface{}, item *domain.Issue, parent string) error {
	var err error
	item.Initialize()

	if item.Reporter == nil || len(item.Reporter.Id) == 0 {
		if item.Project != nil && len(item.Project.Id) != 0 {
			item.Project, err = t.ProjectDb().Load(tx, item.Project.Id)
			if err != nil {
				return err
			}
			if item.Project.ProjectLead != nil && len(item.Project.ProjectLead.Id) != 0 {
				item.Reporter = item.Project.ProjectLead
			}
		}
		if item.Reporter == nil || len(item.Reporter.Id) == 0 {
			item.Reporter = t.Base.CurrentSession().User
		}
	}

	if item.Priority == nil || len(item.Priority.Id) == 0 {
		item.Priority, err = t.PriorityDb().GetFirst(tx)
		if err != nil {
			return err
		}
	}

	err = item.Validate()
	if err != nil {
		return err
	}

	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_CREATE_ISSUE
	ok, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, item, permission)

	if !ok {
		return errors.New("db.Base.AccessDenied()")
	}

	if len(item.IdProject.String) != 0 {
		project, e := t.ProjectDb().NextNumber(tx, item.Project.Id)
		if e != nil {
			return e
		}
		item.Pkey = project.IssueKey()
	} else {
		next, e := t.Base.SequenceNextValue(tx, operations.ISSUE_SEQUENCE_NAME)
		if e != nil {
			return e
		}
		item.Pkey = fmt.Sprintf(operations.ISSUE_NO_PROJECT_MASK, next)
	}

	statuses, err := t.StatusDb().WorkflowStepAvailableStatus(tx, item.Workflow, nil)
	if err != nil {
		return err
	}

	if len(statuses) != 1 {
		return fmt.Errorf("cannot find the first status for the given workflow")
	}

	item.IdStatus = statuses[0].NextStatus.Id

	if len(parent) != 0 {
		item.IdParent.Valid = true
		item.IdParent.String = parent
	}

	err = t.Base.Executor(tx).Insert(item)
	if err != nil {
		return err
	}

	err = t.SubscriptionAdd(tx, item, item.Reporter)
	if err != nil {
		return err
	}
	if item.Assignee != nil && len(item.Assignee.Id) != 0 {
		err = t.SubscriptionAdd(tx, item, item.Assignee)
		if err != nil {
			return err
		}
	}
	return err
}

func (t *IssueDb) Load(tx interface{}, id string, pkey string) (*domain.Issue, error) {
	var err error
	item := &domain.Issue{}
	if len(id) != 0 {
		err = t.Base.Executor(tx).SelectOne(item, "select * from issue where id=$1", id)
	} else {
		err = t.Base.Executor(tx).SelectOne(item, "select * from issue where pkey=$1", pkey)
	}

	if err != nil {
		return nil, err
	}

	item.Initialize()

	permission := &domain.PermissionName{Name: domain.PERMISSION_BROWSE_PROJECT}
	ok, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, item, permission)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("db.Base.AccessDenied()")
	}

	if item.IdAssignee.Valid {
		item.Assignee, err = t.UserDb().Load(tx, item.IdAssignee.String)
		if err != nil {
			return nil, err
		}
	}

	item.Priority, err = t.PriorityDb().Load(tx, item.IdPriority)
	if err != nil {
		return nil, err
	}
	item.Reporter, err = t.UserDb().Load(tx, item.IdReporter)
	if err != nil {
		return nil, err
	}
	item.Status, err = t.StatusDb().Load(tx, item.IdStatus)
	if err != nil {
		return nil, err
	}
	item.Workflow, err = t.StatusDb().WorkflowLoad(tx, item.IdWorkflow)
	if err != nil {
		return nil, err
	}

	if item.IdProject.Valid {
		item.Project, err = t.ProjectDb().Load(tx, item.IdProject.String)
		if err != nil {
			return nil, err
		}
	}

	return item, nil
}

func (t *IssueDb) Remove(tx interface{}, id string) (*domain.Issue, error) {
	item, err := t.Load(tx, id, "")
	if err != nil {
		return nil, err
	}

	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_DELETE_ISSUE
	ok, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, item, permission)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("db.Base.AccessDenied()")
	}

	children, err := t.Base.Executor(tx).SelectInt("select count(*) from issue where idparent=$1", item.Id)
	if children != 0 {
		return nil, fmt.Errorf("cannot remove the issue because there are %s subtasks", children)
	}

	query := "delete from issue_subscription where idissue=$1"
	_, err = t.Base.Executor(tx).Exec(query, item.Id)
	if err != nil {
		return nil, err
	}

	err = t.CommentRemoveAll(tx, item)
	if err != nil {
		return nil, err
	}

	err = t.AttachmentRemoveAll(tx, item)
	if err != nil {
		return nil, err
	}

	_, err = t.Base.Executor(tx).Delete(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

type IssueHierrachy struct {
	Id       string
	Pkey     string
	IdParent sql.NullString
	Level    int
	Path     string
}

func (t *IssueDb) FindRoot(tx interface{}, id string) (root string, err error) {
	currItem := &IssueHierrachy{}
	err = t.Base.Executor(tx).SelectOne(currItem, "select * from view_issue_parents where id=$1", id)

	if err != nil {
		return
	}

	if !currItem.IdParent.Valid {
		root = id
		return
	}
	paths := strings.Split(currItem.Path, ":")
	root = paths[0]
	return
}

func (t *IssueDb) Update(tx interface{}, item *domain.Issue) error {
	err := item.Validate()
	if err != nil {
		return err
	}

	oldItem, err := t.Load(tx, item.Id, "")
	if err != nil {
		return err
	}

	if !oldItem.AcceptUpdates() {
		return fmt.Errorf("cannot update a resolved issue")
	}

	//validar permisos sobre esta accion
	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_EDIT_ISSUE
	ok, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, oldItem, permission)
	if !ok {
		return errors.New("db.Base.AccessDenied()")
	}
	if oldItem.IdAssignee.String != item.IdAssignee.String {
		permission.Name = domain.PERMISSION_ASSIGN_USER
		ok, err = t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, oldItem, permission)
		if !ok {
			return errors.New("db.Base.AccessDenied()")
		}
	}
	if oldItem.IdReporter != item.IdReporter {
		permission.Name = domain.PERMISSION_CHANGE_REPORTER
		ok, err = t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, oldItem, permission)
		if !ok {
			return errors.New("db.Base.AccessDenied()")
		}
	}
	if oldItem.DueDate != item.DueDate {
		permission.Name = domain.PERMISSION_CHANGE_DUEDATE
		ok, err = t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, oldItem, permission)
		if !ok {
			item.DueDate = oldItem.DueDate
		}
	}

	item.DateCreated = oldItem.DateCreated
	item.IdStatus = oldItem.IdStatus
	item.IdProject = oldItem.IdProject
	item.IdWorkflow = oldItem.IdWorkflow
	item.Pkey = oldItem.Pkey
	item.IdParent = oldItem.IdParent

	_, err = t.Base.Executor(tx).Update(item)

	if err != nil {
		return err
	}

	if item.Assignee != nil && len(item.Assignee.Id) != 0 && (oldItem.Assignee == nil || oldItem.Assignee.Id != item.Assignee.Id) {
		err = t.SubscriptionAdd(tx, item, item.Assignee)
		if err != nil {
			return err
		}
	}

	if item.Reporter.Id != oldItem.Reporter.Id {
		err = t.SubscriptionAdd(tx, item, item.Reporter)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *IssueDb) Children(tx interface{}, issue *domain.Issue) ([]database.IssueGrid, error) {
	var items []database.IssueGrid
	_, err := t.Base.Executor(tx).Select(&items, "select * from view_issues where idparent=$1 order by datecreated", issue.Id)
	if err != nil {
		return nil, err
	}
	for i := range items {
		item := &items[i]
		item.Initialize()
	}
	return items, nil
}

func (t *IssueDb) MoveProject(tx interface{}, issue *domain.Issue, target *domain.Project) (*domain.Issue, error) {
	if target == nil || len(target.Id) == 0 {
		return nil, fmt.Errorf("missing project")
	}

	oldIssue, err := t.Load(tx, issue.Id, "")
	if err != nil {
		return nil, err
	}
	oldProject, err := t.ProjectDb().Load(tx, target.Id)
	if err != nil {
		return nil, err
	}

	if issue.Project != nil && oldIssue.Project != nil && oldIssue.Project.Id == target.Id {
		return nil, fmt.Errorf("the issue already belongs to this project, there is nothing to do")
	}

	target = oldProject
	issue = oldIssue

	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_DELETE_ISSUE
	ok, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, issue, permission)
	if !ok {
		return nil, errors.New("db.Base.AccessDenied()")
	}

	project, err := t.ProjectDb().NextNumber(tx, target.Id)
	if err != nil {
		return nil, err
	}
	issue.Pkey = project.IssueKey()
	issue.Project = target
	issue.IdProject.Valid = false

	permission.Name = domain.PERMISSION_CREATE_ISSUE
	ok, err = t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, issue, permission)
	if !ok {
		return nil, errors.New("db.Base.AccessDenied()")
	}

	if _, err = t.Base.Executor(tx).Update(issue); err != nil {
		return nil, err
	}

	if issue, err = t.Load(tx, issue.Id, issue.Pkey); err != nil {
		return nil, err
	} else {
		return issue, nil
	}
}
