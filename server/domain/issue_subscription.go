package domain

import (
	"time"
)

type IssueSubscription struct {
	Id          string    `json:"-"`
	DateCreated time.Time `json:"dateCreated"`
	IdIssue     string    `json:"idIssue"`
	IdUser      string    `json:"idUser"`
	Issue       *Issue    `json:"issue" db:"-"`
	User        *User     `json:"user" db:"-"`
}

func (u *IssueSubscription) Initialize() {
	if u.Issue != nil && len(u.Issue.Id) != 0 {
		u.IdIssue = u.Issue.Id
	} else if len(u.IdIssue) != 0 {
		if u.Issue == nil {
			u.Issue = &Issue{}
		}
		u.Issue.Id = u.IdIssue
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
