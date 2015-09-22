package mock

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
)

func SessionCreate(op *database.DbOperations, tx interface{}, user *domain.User) (*domain.Session, error) {
	return op.AccountDb.LoginUser(tx, user.Email, user.Password, "127.0.0.1", false)
}

func SessionWithUser(tx interface{}, op *database.DbOperations, isAdmin bool) (*domain.Session, error) {
	user := User()
	user.IsSystemAdministrator = isAdmin
	err := UserCreate(op, tx, user)
	if err != nil {
		return nil, err
	}
	return op.AccountDb.LoginUser(tx, user.Email, user.Password, "127.0.0.1", false)
}

func SessionSetContext(op *database.DbOperations, tx interface{}, isAdmin bool) (*domain.Session, error) {
	session, err := SessionWithUser(tx, op, isAdmin)
	if err != nil {
		return nil, err
	}
	op.Db.SetCurrentSession(session)
	return session, nil
}
