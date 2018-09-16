package application

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/mauleyzaola/gorp"
	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/dal/pg/account"
	"github.com/mauleyzaola/issue-tracker/server/dal/pg/bootstrap"
	"github.com/mauleyzaola/issue-tracker/server/dal/pg/file_item"
	"github.com/mauleyzaola/issue-tracker/server/dal/pg/issue"
	"github.com/mauleyzaola/issue-tracker/server/dal/pg/permission"
	"github.com/mauleyzaola/issue-tracker/server/dal/pg/priority"
	"github.com/mauleyzaola/issue-tracker/server/dal/pg/project"
	"github.com/mauleyzaola/issue-tracker/server/dal/pg/session"
	"github.com/mauleyzaola/issue-tracker/server/dal/pg/status"
	"github.com/mauleyzaola/issue-tracker/server/dal/pg/user"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecweb/setup"
)

type Application struct {
	Setup *setup.Setup
	Db    *database.DbOperations
}

func ParseConfiguration(fileName, rootDir string) (app *Application) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal("Cannot read from configuration file ", err.Error())
		panic(err)
	}

	app = &Application{}
	app.Setup = &setup.Setup{}
	if err = json.Unmarshal(data, app.Setup); err != nil {
		log.Fatal("cannot parse config.json: ", err.Error())
	}

	//if there is no environment on the configuration file, exit with error
	if len(app.Setup.Environment) == 0 {
		log.Fatal(fmt.Errorf("there is no environment set in the configuration file"))
	}

	app.Setup.RootChDir = rootDir

	err = app.Setup.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	var db database.Db
	var accountDb database.Account
	var bootstrapDb database.Bootstrap
	var fileItemDb database.FileItem
	var issueDb database.Issue
	var permissionDb database.Permission
	var priorityDb database.Priority
	var projectDb database.Project
	var sessionDb database.Session
	var statusDb database.Status
	var userDb database.User

	dbmap, ok := app.Setup.Relational.(*gorp.DbMap)
	if !ok {
		log.Fatal("cannot cast to *dbmap")
	}
	db = pg.New(dbmap)

	accountDb = account.New(db)
	bootstrapDb = bootstrap.New(db)
	fileItemDb = fileitem.New(db)
	issueDb = issue.New(db)
	permissionDb = permission.New(db)
	priorityDb = priority.New(db)
	projectDb = project.New(db)
	sessionDb = session.New(db)
	statusDb = status.New(db)
	userDb = user.New(db)

	app.Db = &database.DbOperations{}
	ops := app.Db
	ops.Db = db
	app.Db = ops

	ops.AccountDb = accountDb
	ops.BootstrapDb = bootstrapDb
	ops.FileItemDb = fileItemDb
	ops.IssueDb = issueDb
	ops.PermissionDb = permissionDb
	ops.PriorityDb = priorityDb
	ops.ProjectDb = projectDb
	ops.SessionDb = sessionDb
	ops.StatusDb = statusDb
	ops.UserDb = userDb

	//attach dependencies for each implementation
	ops.AccountDb.SetSessionDb(&sessionDb)
	ops.AccountDb.SetUserDb(&userDb)

	ops.BootstrapDb.SetPriorityDb(&priorityDb)
	ops.BootstrapDb.SetStatusDb(&statusDb)
	ops.BootstrapDb.SetUserDb(&userDb)

	ops.IssueDb.SetFileItemDb(&fileItemDb)
	ops.IssueDb.SetPermissionDb(&permissionDb)
	ops.IssueDb.SetPriorityDb(&priorityDb)
	ops.IssueDb.SetProjectDb(&projectDb)
	ops.IssueDb.SetStatusDb(&statusDb)
	ops.IssueDb.SetUserDb(&userDb)

	ops.PermissionDb.SetProjectDb(&projectDb)
	ops.PermissionDb.SetUserDb(&userDb)

	ops.ProjectDb.SetUserDb(&userDb)

	ops.SessionDb.SetUserDb(&userDb)

	ops.StatusDb.SetUserDb(&userDb)

	ops.UserDb.SetAccountDb(&accountDb)

	//register all database objects
	log.Println("registering all the tables in gorp")
	db.Register()

	return
}
