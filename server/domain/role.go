package domain

import (
	"errors"
)

type Role struct {
	Meta *DocumentMetadata `json:"meta" db:"-"`
	Id   string            `json:"id"`
	Name string            `json:"name"`
}

func (u *Role) Validate() (err error) {
	u.Initialize()
	if len(u.Name) == 0 {
		err = errors.New("missing name")
		return
	}

	return
}

func (u *Role) Initialize() {
	u.Meta = &DocumentMetadata{}
	u.Meta.DocumentType = "role"
	u.Meta.FriendName = "Role"
}
