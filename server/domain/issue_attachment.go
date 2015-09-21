package domain

import (
	"errors"
	"time"
)

type IssueAttachment struct {
	Id          string            `json:"id"`
	Meta        *DocumentMetadata `json:"meta" db:"-"`
	DateCreated time.Time         `json:"dateCreated"`
	IdIssue     string            `json:"-"`
	IdFileItem  string            `json:"-"`
	IdUser      string            `json:"-"`
	FileItem    *FileItem         `json:"fileItem" db:"-"`
	User        *User             `json:"user" db:"-"`
	Issue       *Issue            `json:"issue" db:"-"`
}

func (u IssueAttachment) GetMeta() *DocumentMetadata {
	return u.Meta
}

func (u IssueAttachment) GetId() string {
	return u.Id
}

func (u *IssueAttachment) Initialize() {
	u.Meta = &DocumentMetadata{}
	u.Meta.DocumentType = "issueAttachment"
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

	if u.FileItem != nil && len(u.FileItem.Id) != 0 {
		u.IdFileItem = u.FileItem.Id
	} else if len(u.IdFileItem) != 0 {
		if u.FileItem == nil {
			u.FileItem = &FileItem{}
		}
		u.FileItem.Id = u.IdFileItem
	}

	if u.Issue != nil && len(u.Issue.Id) != 0 {
		u.IdIssue = u.Issue.Id
	} else if len(u.IdIssue) != 0 {
		if u.Issue == nil {
			u.Issue = &Issue{}
		}
		u.Issue.Id = u.IdIssue
	}
}

func (u *IssueAttachment) Validate() (err error) {
	u.Initialize()

	if len(u.IdFileItem) == 0 {
		err = errors.New("missing file item")
		return
	}

	if len(u.IdUser) == 0 {
		err = errors.New("missing user")
		return
	}

	if len(u.IdIssue) == 0 {
		err = errors.New("missing issue")
		return
	}

	return
}
