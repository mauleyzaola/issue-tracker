package domain

import (
	"errors"
	"time"
)

type Session struct {
	Id          string    `json:"id"`
	IdUser      string    `json:"-"`
	User        *User     `json:"user" db:"-"`
	DateCreated time.Time `json:"dateCreated"`
	Expires     time.Time `json:"expires"`
	IpAddress   string    `json:"ipAddress"`
}

func (u *Session) Initialize() {
	u.DateCreated = time.Now()
	u.Expires = u.DateCreated.Add(u.Expiration())

	if u.User != nil && len(u.User.Id) != 0 {
		u.IdUser = u.User.Id
	} else if len(u.IdUser) != 0 {
		if u.User == nil {
			u.User = &User{}
		}
		u.User.Id = u.IdUser
	}
}

func (u *Session) Expiration() time.Duration {
	return time.Hour * 36
}

func (u *Session) Validate() error {
	u.Initialize()
	if len(u.IpAddress) == 0 {
		return errors.New("missing ip address")
	}
	return nil
}
