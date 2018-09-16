package database

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecweb/setup"
)

type Bootstrap interface {
	PriorityDb() Priority
	SetPriorityDb(item *Priority)

	StatusDb() Status
	SetStatusDb(item *Status)

	UserDb() User
	SetUserDb(item *User)

	//Bootstraps the database and configuration for the first time the app runs
	//database and db objects must exist before this method is executed
	BootstrapAll(tx interface{}, conn *setup.ConfigurationDatabase, environment string, chdir string) error

	//Creates the first admin user in database, if it already exists returns true on first return value
	BootstrapAdminUser(tx interface{}) (exists bool, admin *domain.User, err error)

	//Make the chdir point to the root of the application
	ResetApplicationPath(chdir string) error

	//Executes any migration tool to upgrade db objects and data
	UpgradeDbScripts(conn *setup.ConfigurationDatabase, environment string, chdir string) error

	//Bootstraps all the permission names
	CreatePermissionNames(tx interface{}) error

	//Add basic workflows as part of the startup
	BootstrapWorkflows(tx interface{}) error

	//Add some values to catalog
	BootstrapPriorities(tx interface{}) error
}
