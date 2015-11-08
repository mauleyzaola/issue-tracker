package setup

import (
	"fmt"

	"github.com/go-gorp/gorp"
	mgo "gopkg.in/mgo.v2"
)

type Application struct {
	RootChDir    string
	Environment  string
	Relational   *gorp.DbMap
	BaseUrl      string
	BaseApiName  string
	Mongo        *mgo.Database
	RelationalDb *ConfigurationDatabase
	MongoDb      *ConfigurationDatabase
	IndexerUrl   string
	Indexer      interface{}
}

type ConfigurationDatabase struct {
	Host         string
	DatabaseName string
	Username     string
	Password     string
	Driver       string
}

func (t *ConfigurationDatabase) ConnectionString() (string, error) {
	switch t.Driver {
	case "mssql":
		return fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;", t.Host, t.Username, t.Password, t.DatabaseName), nil
	case "postgres":
		return fmt.Sprintf("user=%s host=%s dbname=%s password=%s sslmode=disable", t.Username, t.Host, t.DatabaseName, t.Password), nil
	default:
		return "", fmt.Errorf("not supported database driver")
	}
}
