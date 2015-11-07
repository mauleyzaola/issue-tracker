package application

func (a *Application) BootstrapApplication() error {

	//bootstrap database, users and default configuration
	tx, err := a.Db.Db.Begin()
	if err != nil {
		return err
	}

	err = a.Db.BootstrapDb.BootstrapAll(tx, a.Setup.Environment, a.Setup.RootChDir)
	if err != nil {
		a.Setup.Db.Db.Rollback(tx)
		return err
	}

	return a.Db.Db.Commit(tx)
}
