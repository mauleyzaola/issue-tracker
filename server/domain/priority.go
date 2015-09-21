package domain

import (
	"errors"
	"time"
)

type Priority struct {
	Meta         *DocumentMetadata `json:"meta" db:"-"`
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	DateCreated  time.Time         `json:"dateCreated"`
	LastModified *time.Time        `json:"lastModified"`
}

func (u *Priority) Validate() (err error) {
	u.Initialize()
	if len(u.Name) == 0 {
		err = errors.New("name is missing")
	}

	return
}

func (u *Priority) Initialize() {
	u.Meta = &DocumentMetadata{}
	u.Meta.DocumentType = "priority"
	u.Meta.FriendName = "Priority"
}
