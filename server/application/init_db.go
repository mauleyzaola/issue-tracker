package application

import (
	"database/sql"
	"fmt"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

func (a *Application) initDb() error {

	if a.Setup.PostgresDb != nil {
		connectionString := fmt.Sprintf("user=%s host=%s dbname=%s password=%s sslmode=disable", a.Setup.PostgresDb.Username, a.Setup.PostgresDb.Host, a.Setup.PostgresDb.DatabaseName, a.Setup.PostgresDb.Password)
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			return err
		} else {
			a.Setup.Postgres = &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
		}

	}

	return nil
}
