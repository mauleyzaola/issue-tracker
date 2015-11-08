package application

import (
	"database/sql"
	"fmt"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

func (a *Application) initDb() error {
	if a.Setup.RelationalDb != nil {
		var (
			dialect          gorp.Dialect
			connectionString string
		)

		connectionString, err := a.Setup.RelationalDb.ConnectionString()
		if err != nil {
			return err
		}

		db, err := sql.Open(a.Setup.RelationalDb.Driver, connectionString)
		if err != nil {
			return err
		}
		switch a.Setup.RelationalDb.Driver {
		case "postgres":
			dialect = gorp.PostgresDialect{}
		default:
			return fmt.Errorf("unsopported db dialect")
		}
		a.Setup.Relational = &gorp.DbMap{Db: db, Dialect: dialect}
	} else {
		return fmt.Errorf("cannot open any relational database")
	}

	return nil
}
