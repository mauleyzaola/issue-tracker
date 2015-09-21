package bootstrap

import (
	"database/sql"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
)

type BootstrapDb struct {
	Base     *pg.Db
	statusDb *database.Status
	userDb   *database.User
}

func New(db database.Db) *BootstrapDb {
	base := db.(*pg.Db)
	return &BootstrapDb{Base: base}
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

func (db *BootstrapDb) BootstrapAll(tx interface{}, environment string, chdir string) error {
	err := db.UpgradeDbScripts(environment, chdir)
	if err != nil {
		return err
	}

	//create default user if it doesn't exist
	_, _, err = db.BootstrapAdminUser(tx)
	if err != nil {
		return err
	}

	//add missing permission names
	err = db.CreatePermissionNames(tx)
	if err != nil {
		return err
	}

	//point to the root directory of the application again
	return db.ResetApplicationPath(chdir)
}

func (db *BootstrapDb) BootstrapAdminUser(tx interface{}) (bool, *domain.User, error) {
	var users []domain.User
	_, err := db.Base.Executor(tx).Select(&users, "select * from users where issystemadministrator=$1 and isactive=$1", true)
	if err != nil && err != sql.ErrNoRows {
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
	err = db.UserDb().Create(tx, admin)

	return false, admin, nil
}

func (db *BootstrapDb) ResetApplicationPath(chdir string) error {
	return os.Chdir(chdir)
}

func (t *BootstrapDb) UpgradeDbScripts(environment string, chdir string) error {
	err := t.ResetApplicationPath(chdir)
	if err != nil {
		return err
	}

	err = t.ResetApplicationPath("../dbmigrations/pg")
	if err != nil {
		return err
	}

	var params []string
	if len(environment) != 0 {
		params = append(params, "-env="+environment)
	}
	params = append(params, "up")

	cmd := exec.Command("goose", params...)
	result, err := cmd.CombinedOutput()
	if len(result) != 0 {
		log.Println(string(result))
	}

	return err
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
