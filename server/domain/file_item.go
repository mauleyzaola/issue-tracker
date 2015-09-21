package domain

import (
	"errors"
	"time"
)

type FileItem struct {
	Id          string            `json:"id"`
	IdUser      string            `json:"-"`
	Bytes       int64             `json:"bytes"`
	Extension   string            `json:"extension"`
	Name        string            `json:"name"`
	DateCreated time.Time         `json:"dateCreated"`
	User        *User             `json:"user" db:"-"`
	Meta        *DocumentMetadata `json:"meta" db:"-"`
	MimeType    string            `json:"mimeType"`
	FileData    []byte            `json:"-"`
}

func (u FileItem) GetMeta() *DocumentMetadata {
	return u.Meta
}

func (u *FileItem) Initialize() {
	u.Meta = &DocumentMetadata{}
	u.Meta.DocumentType = "fileItem"
	u.Meta.FriendName = "File"
	u.Meta.Id = u.Id

	if u.User != nil && len(u.User.Id) != 0 {
		u.IdUser = u.User.Id
	} else if len(u.IdUser) != 0 {
		if u.User == nil {
			u.User = &User{}
		}
		u.User.Id = u.IdUser
	}
}

func (u *FileItem) Validate() (err error) {
	u.Initialize()

	if len(u.Extension) == 0 {
		err = errors.New("Missing file extension")
		return
	}
	if len(u.Name) == 0 {
		err = errors.New("Missing file name")
		return
	}
	if len(u.MimeType) == 0 {
		err = errors.New("Missing mime")
		return
	}

	if u.FileData == nil || len(u.FileData) == 0 {
		err = errors.New("File data is empty")
		return
	}

	return
}
