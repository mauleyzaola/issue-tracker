package domain

import (
	"errors"
)

type ProjectRole struct {
	Meta      *DocumentMetadata `json:"meta" db:"-"`
	Id        string            `json:"id"`
	IdProject string            `json:"-"`
	IdRole    string            `json:"-"`
	Project   *Project          `json:"project" db:"-"`
	Role      *Role             `json:"role" db:"-"`
}

func (u *ProjectRole) Validate() (err error) {
	u.Initialize()

	if len(u.IdRole) == 0 {
		err = errors.New("missing role")
		return
	}

	return
}

func (u *ProjectRole) Initialize() {
	u.Meta = &DocumentMetadata{}
	u.Meta.DocumentType = "projectRole"
	u.Meta.FriendName = "Project Role"

	if u.Project != nil && len(u.Project.Id) != 0 {
		u.IdProject = u.Project.Id
	} else if len(u.IdProject) != 0 {
		if u.Project == nil {
			u.Project = &Project{}
		}
		u.Project.Id = u.IdProject
	}

	if u.Role != nil && len(u.Role.Id) != 0 {
		u.IdRole = u.Role.Id
	} else if len(u.IdRole) != 0 {
		if u.Role == nil {
			u.Role = &Role{}
		}
		u.Role.Id = u.IdRole
	}
}
