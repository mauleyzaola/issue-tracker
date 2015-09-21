package domain

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id                    string            `json:"id"`
	LastName              string            `json:"lastName"`
	Name                  string            `json:"name"`
	Email                 string            `json:"email"`
	Password              string            `json:"-"`
	TokenEmail            string            `json:"tokenEmail"`
	TokenExpires          time.Time         `json:"tokenExpires"`
	LoginCount            int               `json:"loginCount"`
	LastLogin             *time.Time        `json:"lastLogin"`
	IsActive              bool              `json:"isActive"`
	IsSystemAdministrator bool              `json:"isSystemAdministrator"`
	DateCreated           time.Time         `json:"dateCreated"`
	LastModified          *time.Time        `json:"lastModified"`
	Meta                  *DocumentMetadata `json:"meta" db:"-"`
	Metadata              *UserMeta         `json:"metadata" db:"-"`
}

func (u User) GetMeta() *DocumentMetadata {
	return u.Meta
}

func (u *User) GetId() string {
	return u.Id
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.Name, u.LastName)
}

func (u *User) Validate() (err error) {
	u.Initialize()
	if len(u.Email) == 0 {
		err = errors.New("email is missing")
		return
	}

	return
}

func (u *User) Initialize() {
	u.Meta = &DocumentMetadata{}
	u.Meta.Id = u.Id
	u.Meta.DocumentType = "user"
	u.Meta.FriendName = "User"
}
