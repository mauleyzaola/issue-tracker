package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
)

//configuration files are parsed and application is initialized
func main() {
	jsonFile := flag.String("config", "config.json", "Config file. Make a copy of config.json.sample")
	flag.Parse()

	log.Printf("Configuration file:%s\n", *jsonFile)

	app := application.ParseConfiguration(*jsonFile)

	//get the base root path of the application
	rootChDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	app.Setup.RootChDir = rootChDir

	err = app.BootstrapApplication()
	if err != nil {
		log.Fatal(err)
	}

	//register all api routes
	app.Router()

	graceful.PostHook(func() {
		log.Println("Application is closing now... releasing resources")
		//		app.Setup.Db.Db.Close()
		//		app.Db.Db.Close()
		log.Println(app.Db == nil)
	})
	goji.Serve()
}
