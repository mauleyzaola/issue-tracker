package domain

import (
	"database/sql"
	"fmt"
)

type ProjectRoleMember struct {
	Id            string         `json:"id"`
	IdProjectRole string         `json:"-"`
	IdGroup       sql.NullString `json:"-"`
	IdUser        sql.NullString `json:"-"`
	Group         *Group         `json:"group" db:"-"`
	User          *User          `json:"user" db:"-"`
	ProjectRole   *ProjectRole   `json:"projectRole" db:"-"`
}

func (u *ProjectRoleMember) Validate() (err error) {
	u.Initialize()

	if len(u.IdGroup.String) != 0 && len(u.IdUser.String) != 0 {
		return fmt.Errorf("cannot set group and user properties for the same member")
	}
	return
}

func (u *ProjectRoleMember) Initialize() {
	if u.Group != nil && len(u.Group.Id) != 0 {
		u.IdGroup.Valid = true
		u.IdGroup.String = u.Group.Id
	} else if len(u.IdGroup.String) != 0 {
		if u.Group == nil {
			u.Group = &Group{}
		}
		u.Group.Id = u.IdGroup.String
	}
	u.IdGroup.Valid = len(u.IdGroup.String) != 0

	if u.User != nil && len(u.User.Id) != 0 {
		u.IdUser.Valid = true
		u.IdUser.String = u.User.Id
	} else if len(u.IdUser.String) != 0 {
		if u.User == nil {
			u.User = &User{}
		}
		u.User.Id = u.IdUser.String
	}
	u.IdUser.Valid = len(u.IdUser.String) != 0

	if u.ProjectRole != nil && len(u.ProjectRole.Id) != 0 {
		u.IdProjectRole = u.ProjectRole.Id
	} else if len(u.IdProjectRole) != 0 {
		if u.ProjectRole == nil {
			u.ProjectRole = &ProjectRole{}
		}
		u.ProjectRole.Id = u.IdProjectRole
	}
}
