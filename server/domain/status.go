package domain

import (
	"errors"
	"time"
)

type Status struct {
	Id           string     `json:"id"`
	IdWorkflow   string     `json:"-"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	DateCreated  time.Time  `json:"dateCreated"`
	LastModified *time.Time `json:"lastModified"`
	Workflow     *Workflow  `json:"workflow" db:"-"`
}

func (u *Status) Initialize() {
	if u.Workflow != nil && len(u.Workflow.Id) != 0 {
		u.IdWorkflow = u.Workflow.Id
	} else if len(u.IdWorkflow) != 0 {
		if u.Workflow == nil {
			u.Workflow = &Workflow{}
		}
		u.Workflow.Id = u.IdWorkflow
	}
}

func (u *Status) Validate() (err error) {
	u.Initialize()
	if len(u.Name) == 0 {
		err = errors.New("name is missing")
		return
	}

	if u.Workflow == nil || len(u.Workflow.Id) == 0 {
		err = errors.New("missing workflow")
		return
	}

	return
}
