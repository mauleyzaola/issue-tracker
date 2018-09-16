package setup

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
	"github.com/mauleyzaola/gorp"
	mgo "gopkg.in/mgo.v2"
)

func (t *Setup) InitDb() error {

	//configure relational database
	if t.RelationalDb != nil {
		var (
			dialect          gorp.Dialect
			connectionString string
		)

		connectionString, err := t.RelationalDb.ConnectionString()
		if err != nil {
			return err
		}

		switch t.RelationalDb.Driver {
		case "mssql":
			dialect = gorp.SqlServerDialect{}
		case "postgres":
			dialect = gorp.PostgresDialect{}
		default:
			return fmt.Errorf("unsopported relational driver")
		}

		db, err := sql.Open(t.RelationalDb.Driver, connectionString)
		if err != nil {
			return err
		}
		t.Relational = &gorp.DbMap{Db: db, Dialect: dialect}
	}

	//configure NoSql database
	if t.NoSqlDb != nil {
		connectionString := fmt.Sprintf("%s", t.NoSqlDb.Host)
		session, err := mgo.Dial(connectionString)
		if err != nil {
			return err
		} else {
			mgoDatabase := session.DB(t.NoSqlDb.DatabaseName)
			t.NoSql = mgoDatabase
		}

	}

	return nil
}
