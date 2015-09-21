package database

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
)

type Session interface {
	UserDb() User
	SetUserDb(item *User)

	//Returns a list of all the active sessions for a given user
	List(tx interface{}, user *domain.User) ([]domain.Session, error)

	//Loads the session from database
	Load(tx interface{}, token string) (*domain.Session, error)

	//Removes the session from database
	Remove(tx interface{}, token string, validateCurrentSession bool) error

	//Updates the session from database
	Update(tx interface{}, session *domain.Session) error

	//Creates a new session for a login user
	Create(tx interface{}, session *domain.Session) error

	//Middleware function to retrieve the current session based on the token provided
	ValidateToken(tx interface{}, token string) (*domain.Session, error)
}
