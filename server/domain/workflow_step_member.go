package domain

import (
	"database/sql"
	"errors"
)

type WorkflowStepMember struct {
	Id             string         `json:"id"`
	IdWorkflowStep string         `json:"-"`
	IdUser         sql.NullString `json:"-"`
	IdGroup        sql.NullString `json:"-"`
	WorkflowStep   *WorkflowStep  `json:"workflowStep" db:"-"`
	User           *User          `json:"user" db:"-"`
	Group          *Group         `json:"group" db:"-"`
}

func (u *WorkflowStepMember) Initialize() {
	if u.WorkflowStep != nil && len(u.WorkflowStep.Id) != 0 {
		u.IdWorkflowStep = u.WorkflowStep.Id
	} else if len(u.IdWorkflowStep) != 0 {
		if u.WorkflowStep == nil {
			u.WorkflowStep = &WorkflowStep{}
		}
		u.WorkflowStep.Id = u.IdWorkflowStep
	}

	if u.Group != nil && len(u.Group.Id) != 0 {
		u.IdGroup.Valid = true
		u.IdGroup.String = u.Group.Id
	} else if len(u.IdGroup.String) != 0 {
		if u.Group == nil {
			u.Group = &Group{}
		}
		u.Group.Id = u.IdGroup.String
	}

	if u.User != nil && len(u.User.Id) != 0 {
		u.IdUser.Valid = true
		u.IdUser.String = u.User.Id
	} else if len(u.IdUser.String) != 0 {
		if u.User == nil {
			u.User = &User{}
		}
		u.User.Id = u.IdUser.String
	}
}

func (u *WorkflowStepMember) Validate() (err error) {
	u.Initialize()

	if len(u.IdGroup.String) != 0 && len(u.IdUser.String) != 0 {
		err = errors.New("only one member is allowed for each row")
		return
	}

	if len(u.IdGroup.String)+len(u.IdUser.String) == 0 {
		err = errors.New("an user or group is required")
		return
	}

	if len(u.IdWorkflowStep) == 0 {
		err = errors.New("step is missing")
		return
	}

	return
}
