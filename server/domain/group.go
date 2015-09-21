package domain

import (
	"errors"
	"time"
)

type Group struct {
	Meta         *DocumentMetadata `json:"meta" db:"-"`
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	DateCreated  time.Time         `json:"dateCreated"`
	LastModified *time.Time        `json:"lastModified"`
}

func (u *Group) Validate() (err error) {
	u.Initialize()
	if len(u.Name) == 0 {
		err = errors.New("name is missing")
		return
	}

	return
}

func (t *Group) Initialize() {
	t.Meta = &DocumentMetadata{}
	t.Meta.DocumentType = "group"
	t.Meta.FriendName = "Group"
}
