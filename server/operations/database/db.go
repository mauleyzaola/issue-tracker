package database

//base del db
type Db interface {

	//Creates a new database transaction and returns its object
	Begin() (interface{}, error)

	//Commmits a database transaction
	Commit(tran interface{}) error

	//Rollbacks a database transaction
	Rollback(tran interface{}) error

	//Starts the SQL Tracing
	SqlTraceOn()

	//Stops the SQL Tracing
	SqlTraceOff()

	//If the engine is compatible with any relational database,
	//it will register all the internal structs to tables
	Register()

	//Contains the object which holds information about the current session connected
	//to the system
	CurrentSession

	//Closes the db connections and releases resources
	Close()
}
