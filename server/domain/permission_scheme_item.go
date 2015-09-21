package domain

import (
	"database/sql"
	"errors"
)

type PermissionSchemeItem struct {
	Id                 string            `json:"id"`
	IdPermissionScheme string            `json:"-"`
	IdPermissionName   string            `json:"-"`
	IdRole             sql.NullString    `json:"-"`
	IdGroup            sql.NullString    `json:"-"`
	IdUser             sql.NullString    `json:"-"`
	User               *User             `json:"user" db:"-"`
	Group              *Group            `json:"group" db:"-"`
	Role               *Role             `json:"role" db:"-"`
	PermissionScheme   *PermissionScheme `json:"permissionScheme" db:"-"`
	PermissionName     *PermissionName   `json:"permissionName" db:"-"`
}

func (u *PermissionSchemeItem) Validate() (err error) {
	u.Initialize()
	if len(u.IdPermissionName) == 0 {
		err = errors.New("missing permission name")
		return
	}

	if len(u.IdPermissionScheme) == 0 {
		err = errors.New("missing permission scheme")
		return
	}
	return
}

func (u *PermissionSchemeItem) Initialize() {
	if u.Group != nil && len(u.Group.Id) != 0 {
		u.IdGroup.Valid = true
		u.IdGroup.String = u.Group.Id
	} else if len(u.IdGroup.String) != 0 {
		if u.Group == nil {
			u.Group = &Group{}
		}
		u.Group.Id = u.IdGroup.String
	}

	if u.User != nil && len(u.User.Id) != 0 {
		u.IdUser.Valid = true
		u.IdUser.String = u.User.Id
	} else if len(u.IdUser.String) != 0 {
		if u.User == nil {
			u.User = &User{}
		}
		u.User.Id = u.IdUser.String
	}

	if u.Role != nil && len(u.Role.Id) != 0 {
		u.IdRole.Valid = true
		u.IdRole.String = u.Role.Id
	} else if len(u.IdRole.String) != 0 {
		if u.Role == nil {
			u.Role = &Role{}
		}
		u.Role.Id = u.IdRole.String
	}

	if u.PermissionScheme != nil && len(u.PermissionScheme.Id) != 0 {
		u.IdPermissionScheme = u.PermissionScheme.Id
	} else if len(u.IdPermissionScheme) != 0 {
		if u.PermissionScheme == nil {
			u.PermissionScheme = &PermissionScheme{}
		}
		u.PermissionScheme.Id = u.IdPermissionScheme
	}

	if u.PermissionName != nil && len(u.PermissionName.Id) != 0 {
		u.IdPermissionName = u.PermissionName.Id
	} else if len(u.IdPermissionName) != 0 {
		if u.PermissionName == nil {
			u.PermissionName = &PermissionName{}
		}
		u.PermissionName.Id = u.IdPermissionName
	}
}
