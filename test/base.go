package test

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mauleyzaola/issue-tracker/server/application"
)

//reads from config file and returns the app object
//TODO: make this work under an isolated enviroment like root from upstart
func factoryDb() *application.Application {
	home := os.Getenv("HOME")
	rootPath := filepath.Join(home, "go", "src", "github.com", "mauleyzaola", "issue-tracker")

	configFile := filepath.Join(rootPath, "test", "config.json")

	app := application.ParseConfiguration(configFile, rootPath)
	app.Setup.RootChDir = rootPath

	app.Db.Db.Register()

	//execute bootstrappers in test env
	if err := app.BootstrapApplication(); err != nil {
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
