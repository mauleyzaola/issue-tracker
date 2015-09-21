package domain

import (
	"database/sql"
	"errors"
	"time"
)

type WorkflowStep struct {
	Id           string         `json:"id"`
	Name         string         `json:"name"`
	IdWorkflow   string         `json:"-"`
	IdPrevStatus sql.NullString `json:"-"`
	IdNextStatus string         `json:"-"`
	DateCreated  time.Time      `json:"dateCreated"`
	Resolves     bool           `json:"resolves"`
	Cancels      bool           `json:"cancels"`
	PrevStatus   *Status        `json:"prevStatus" db:"-"`
	NextStatus   *Status        `json:"nextStatus" db:"-"`
	Workflow     *Workflow      `json:"workflow" db:"-"`
}

func (u *WorkflowStep) Initialize() {
	if u.NextStatus != nil && len(u.NextStatus.Id) != 0 {
		u.IdNextStatus = u.NextStatus.Id
	} else if len(u.IdNextStatus) != 0 {
		if u.NextStatus == nil {
			u.NextStatus = &Status{}
		}
		u.NextStatus.Id = u.IdNextStatus
	}

	if u.PrevStatus != nil && len(u.PrevStatus.Id) != 0 {
		u.IdPrevStatus.Valid = true
		u.IdPrevStatus.String = u.PrevStatus.Id
	} else if len(u.IdPrevStatus.String) != 0 {
		if u.PrevStatus == nil {
			u.PrevStatus = &Status{}
		}
		u.PrevStatus.Id = u.IdPrevStatus.String
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

func (u *WorkflowStep) Validate() (err error) {
	u.Initialize()

	if len(u.IdPrevStatus.String) == 0 && len(u.IdNextStatus) == 0 {
		err = errors.New("cannot save the step without any prev or next status")
		return
	}

	if u.IdPrevStatus.String == u.IdNextStatus {
		err = errors.New("prev and next status must be different")
		return
	}

	if len(u.IdNextStatus) == 0 {
		err = errors.New("next step is required")
		return
	}

	if u.Cancels && u.Resolves {
		err = errors.New("step cannot resolve and cancel at the same time")
		return
	}

	if len(u.IdWorkflow) == 0 {
		err = errors.New("workflow is missing")
		return
	}

	if len(u.Name) == 0 {
		err = errors.New("name is missing")
		return
	}

	return
}
