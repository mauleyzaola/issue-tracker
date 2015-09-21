package domain

import (
	"errors"
	"time"
)

type Workflow struct {
	Meta         *DocumentMetadata `json:"meta" db:"-"`
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	DateCreated  time.Time         `json:"dateCreated"`
	LastModified *time.Time        `json:"lastModified"`
}

func (u *Workflow) Validate() (err error) {
	if len(u.Name) == 0 {
		err = errors.New("name is missing")
	}

	return
}

func (u *Workflow) Initialize() {
	u.Meta = &DocumentMetadata{}
	u.Meta.DocumentType = "workflow"
	u.Meta.FriendName = "Workflow"
}
