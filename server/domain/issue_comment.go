package domain

import (
	"errors"
	"time"
)

type IssueComment struct {
	Id           string            `json:"id"`
	Meta         *DocumentMetadata `json:"meta" db:"-"`
	DateCreated  time.Time         `json:"dateCreated"`
	LastModified *time.Time        `json:"lastModified"`
	IdIssue      string            `json:"-"`
	IdUser       string            `json:"-"`
	Body         string            `json:"body"`
	User         *User             `json:"user" db:"-"`
	Issue        *Issue            `json:"issue" db:"-"`
}

func (u IssueComment) GetMeta() *DocumentMetadata {
	return u.Meta
}

func (u IssueComment) GetId() string {
	return u.Id
}

func (u *IssueComment) Validate() (err error) {
	u.Initialize()

	if len(u.Body) == 0 {
		err = errors.New("missing body")
		return
	}
	if len(u.IdIssue) == 0 {
		err = errors.New("missing issue")
		return
	}

	return
}

func (u *IssueComment) Initialize() {
	if u.Meta == nil {
		u.Meta = &DocumentMetadata{}
	}
	u.Meta.DocumentType = "issueComment"
	u.Meta.FriendName = "Comment"
	u.Meta.Id = u.Id

	if u.Issue != nil && len(u.Issue.Id) != 0 {
		u.IdIssue = u.Issue.Id
	} else if len(u.IdIssue) != 0 {
		if u.Issue == nil {
			u.Issue = &Issue{}
		}
		u.IdIssue = u.Issue.Id
	}

	if u.User != nil && len(u.User.Id) != 0 {
		u.IdUser = u.User.Id
	} else if len(u.IdUser) != 0 {
		if u.User == nil {
			u.User = &User{}
		}
		u.User.Id = u.IdUser
	}
}
