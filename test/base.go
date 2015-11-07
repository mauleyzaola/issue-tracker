package test

import (
	"log"
	"path"
	"path/filepath"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/tecutils"
)

//reads from config file and returns the app object
//TODO: make this work under an isolated enviroment like root from upstart
func factoryDb() *application.Application {
	pkg, err := tecutils.GetPackageFullPath("github.com/mauleyzaola/issue-tracker/test")
	rootChDir, err := filepath.Abs(pkg)
	if err != nil {
		log.Fatal(err)
	}
	app := application.ParseConfiguration(path.Join(rootChDir, "config.json"))
	app.Setup.RootChDir = rootChDir

	//execute bootstrappers in test env
	if err = app.BootstrapApplication(); err != nil {
		log.Fatal(err)
	}

	return app
}

var app *application.Application

func initialize() {
	if app == nil {
		app = factoryDb()
	}
}

func Runner(fn func(a *application.Application, tx interface{})) {
	initialize()
	tx, err := app.Db.Db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	if fn != nil {
		fn(app, tx)
	}
	defer app.Db.Db.Rollback(tx)
}
