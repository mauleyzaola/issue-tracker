package setup

import (
	"fmt"
)

type Setup struct {
	RootChDir   string
	Environment string

	//should be a pointer to *gorp.DbMap
	Relational interface{}

	//should be a pointer to *mgo.Database
	NoSql interface{}

	BaseUrl      string
	BaseApiName  string
	RelationalDb *ConfigurationDatabase
	NoSqlDb      *ConfigurationDatabase

	FileStorage string

	Indexer *Indexer
}

type Indexer struct {
	Url string
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
