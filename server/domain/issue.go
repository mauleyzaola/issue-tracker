package domain

import (
	"database/sql"
	"errors"
	"time"
)

type Issue struct {
	Meta            *DocumentMetadata `json:"meta" db:"-"`
	Id              string            `json:"id"`
	IdParent        sql.NullString    `json:"-"`
	Pkey            string            `json:"pkey"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	DateCreated     time.Time         `json:"dateCreated"`
	LastModified    *time.Time        `json:"lastModified"`
	IdStatus        string            `json:"-"`
	IdWorkflow      string            `json:"-"`
	IdPriority      string            `json:"-"`
	IdProject       sql.NullString    `json:"-"`
	IdAssignee      sql.NullString    `json:"-"`
	IdReporter      string            `json:"-"`
	DueDate         time.Time         `json:"dueDate"`
	ResolvedDate    *time.Time        `json:"resolvedDate"`
	CancelledDate   *time.Time        `json:"cancelledDate"`
	Status          *Status           `json:"status" db:"-"`
	Priority        *Priority         `json:"priority" db:"-"`
	Project         *Project          `json:"project" db:"-"`
	Assignee        *User             `json:"assignee" db:"-"`
	Reporter        *User             `json:"reporter" db:"-"`
	Workflow        *Workflow         `json:"workflow" db:"-"`
	DocumentRelated *DocumentMetadata `json:"documentRelated" db:"-"`
}

func (u *Issue) GetMeta() *DocumentMetadata {
	return u.Meta
}

func (u *Issue) GetId() string {
	return u.Id
}

func (u *Issue) AcceptUpdates() bool {
	return u.CancelledDate == nil && u.ResolvedDate == nil
}

func (u *Issue) Validate() (err error) {
	u.Initialize()
	if len(u.Name) == 0 {
		err = errors.New("name is missing")
		return
	}

	if len(u.IdWorkflow) == 0 {
		err = errors.New("workflow is missing")
		return
	}

	if len(u.IdPriority) == 0 {
		err = errors.New("priority is missing")
		return
	}

	if u.DueDate.Year() < 2000 {
		err = errors.New("due date is missing")
		return
	}

	return
}

func (u *Issue) Initialize() {
	u.Meta = &DocumentMetadata{}
	u.Meta.DocumentType = "issue"
	u.Meta.FriendName = "Issue"
	u.Meta.Id = u.Id

	if u.Assignee != nil && len(u.Assignee.Id) != 0 {
		u.IdAssignee.Valid = true
		u.IdAssignee.String = u.Assignee.Id
	} else if len(u.IdAssignee.String) != 0 {
		if u.Assignee == nil {
			u.Assignee = &User{}
		}
		u.Assignee.Id = u.IdAssignee.String
	}

	if u.Project != nil && len(u.Project.Id) != 0 {
		u.IdProject.Valid = true
		u.IdProject.String = u.Project.Id
	} else if len(u.IdProject.String) != 0 {
		if u.Project == nil {
			u.Project = &Project{}
		}
		u.Project.Id = u.IdProject.String
	}

	if u.Priority != nil && len(u.Priority.Id) != 0 {
		u.IdPriority = u.Priority.Id
	} else if len(u.IdPriority) != 0 {
		if u.Priority == nil {
			u.Priority = &Priority{}
		}
		u.Priority.Id = u.IdPriority
	}

	if u.Reporter != nil && len(u.Reporter.Id) != 0 {
		u.IdReporter = u.Reporter.Id
	} else if len(u.IdReporter) != 0 {
		if u.Reporter == nil {
			u.Reporter = &User{}
		}
		u.Reporter.Id = u.IdReporter
	}

	if u.Status != nil && len(u.Status.Id) != 0 {
		u.IdStatus = u.Status.Id
	} else if len(u.IdStatus) != 0 {
		if u.Status == nil {
			u.Status = &Status{}
		}
		u.Status.Id = u.IdStatus
	}

	if u.Workflow != nil && len(u.Workflow.Id) != 0 {
		u.IdWorkflow = u.Workflow.Id
	} else if len(u.IdWorkflow) != 0 {
		if u.Workflow == nil {
			u.Workflow = &Workflow{}
		}
		u.Workflow.Id = u.IdWorkflow
	}
}
