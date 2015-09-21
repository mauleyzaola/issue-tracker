package domain

import (
	"errors"
)

type PermissionScheme struct {
	Meta *DocumentMetadata `json:"meta" db:"-"`
	Id   string            `json:"id"`
	Name string            `json:"name"`
}

func (u *PermissionScheme) Validate() (err error) {
	if len(u.Name) == 0 {
		err = errors.New("name is missing")
		return
	}

	return
}

func (u *PermissionScheme) Initialize() {
	u.Meta = &DocumentMetadata{}
	u.Meta.DocumentType = "permissionScheme"
	u.Meta.FriendName = "Permission Scheme"
}
