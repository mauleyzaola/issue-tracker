package domain

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Project struct {
	Meta               *DocumentMetadata `json:"meta" db:"-"`
	Id                 string            `json:"id"`
	Pkey               string            `json:"pkey"`
	Name               string            `json:"name"`
	IdProjectLead      string            `json:"-"`
	DateCreated        time.Time         `json:"dateCreated"`
	LastModified       *time.Time        `json:"lastModified"`
	Begins             *time.Time        `json:"begins"`
	Ends               *time.Time        `json:"ends"`
	Next               int64             `json:"next"`
	ProjectLead        *User             `json:"projectLead" db:"-"`
	IssueCount         int32             `json:"issueCount"`
	NotResolvedCount   int32             `json:"notResolvedCount"`
	IdPermissionScheme sql.NullString    `json:"-"`
	PermissionScheme   *PermissionScheme `json:"permissionScheme" db:"-"`
}

func (u Project) GetId() string {
	return u.Id
}

func (u Project) GetMeta() *DocumentMetadata {
	return u.Meta
}

func (u *Project) Initialize() {
	u.Meta = &DocumentMetadata{}
	u.Meta.Id = u.Id
	u.Meta.DocumentType = "project"
	u.Meta.FriendName = "Project"
	u.DateCreated = time.Now()

	if u.ProjectLead != nil && len(u.ProjectLead.Id) != 0 {
		u.IdProjectLead = u.ProjectLead.Id
	} else if len(u.IdProjectLead) != 0 {
		if u.ProjectLead == nil {
			u.ProjectLead = &User{}
		}
		u.ProjectLead.Id = u.IdProjectLead
	}

	if u.PermissionScheme != nil && len(u.PermissionScheme.Id) != 0 {
		u.IdPermissionScheme.Valid = true
		u.IdPermissionScheme.String = u.PermissionScheme.Id
	} else if len(u.IdPermissionScheme.String) != 0 {
		if u.PermissionScheme == nil {
			u.PermissionScheme = &PermissionScheme{}
		}
		u.PermissionScheme.Id = u.IdPermissionScheme.String
	}
}

func (u *Project) Validate() (err error) {
	u.Initialize()
	if len(u.Name) == 0 {
		err = errors.New("name is missing")
		return
	}
	if len(u.Pkey) == 0 {
		err = errors.New("key is missing")
		return
	}

	u.Pkey = strings.ToUpper(u.Pkey)
	for i := range u.Pkey {
		ascii := u.Pkey[i]
		if !((ascii >= 48 && ascii <= 57) || (ascii >= 65 && ascii <= 90)) {
			err = errors.New("key is invalid. valid keys must be between A-Z or 0-9")
			break
		}
	}
	return
}

func (u *Project) IssueKey() string {
	return fmt.Sprintf("%s-%d", u.Pkey, u.Next)
}
