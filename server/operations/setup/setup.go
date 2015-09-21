package setup

import (
	"github.com/go-gorp/gorp"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	mgo "gopkg.in/mgo.v2"
)

type Application struct {
	RootChDir   string
	Environment string
	Postgres    *gorp.DbMap
	BaseUrl     string
	BaseApiName string
	Mongo       *mgo.Database
	PostgresDb  *ConfigurationDatabase
	Db          *database.DbOperations
}

type ConfigurationDatabase struct {
	Host         string
	DatabaseName string
	Username     string
	Password     string
}
