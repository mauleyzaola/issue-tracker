package database

import (
	"time"

	"github.com/mauleyzaola/issue-tracker/server/domain"
)

type Account interface {
	SessionDb() Session
	SetSessionDb(item *Session)

	UserDb() User
	SetUserDb(item *User)

	//Tries to load an user given its id or email
	FullLoadUser(tx interface{}, user *domain.User) error

	//Creates a new session object and increases the login count
	//If validatePassword parameter is false, then no authentication method is used
	LoginUser(tx interface{}, email, password string, ipaddress string, validatePassword bool) (*domain.Session, error)

	//Removes the active session for the connected user
	LogoutUser(tx interface{}) error

	//Changes the password of any user
	ChangePassword(tx interface{}, id string, password string) error

	//Changes the password for the connected user
	ChangeMyPassword(tx interface{}, password string) error

	//Delays the session expire date based on the duration provided as parameter
	LoginDelay(tx interface{}, delay time.Duration) error

	//Creates a token for password recovery
	PasswordRecoverCreateToken(tx interface{}, user *domain.User) error

	//Verifies the token for password recovery is good.
	VerifyTokenEmail(tx interface{}, user *domain.User, ipaddress string) error
}
