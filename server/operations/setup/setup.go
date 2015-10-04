package setup

import (
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
)

type Application struct {
	RootChDir   string
	Environment string

	//should be a pointer to *gorp.DbMap
	Postgres    interface{}
	BaseUrl     string
	BaseApiName string

	//should be a pointer to *mgo.Database
	Mongo      interface{}
	PostgresDb *ConfigurationDatabase
	Db         *database.DbOperations
}

type ConfigurationDatabase struct {
	Host         string
	DatabaseName string
	Username     string
	Password     string
}
