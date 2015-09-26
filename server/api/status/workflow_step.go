package status

import (
	"fmt"
	"net/http"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/zenazn/goji/web"
)

type workflowStepGroups struct {
	Selected   []domain.Group `json:"selected"`
	Unselected []domain.Group `json:"unselected"`
}

type workflowStepUsers struct {
	Selected   []domain.User `json:"selected"`
	Unselected []domain.User `json:"unselected"`
}

func (t *Api) workflowStepAvailableUser(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	workflow := &domain.Workflow{Id: t.base.ParamValue("workflow", c, r)}
	prevStatus := &domain.Status{Id: t.base.ParamValue("status", c, r)}
	items, err := t.base.Database.StatusDb.WorkflowStepAvailableUser(tx, workflow, prevStatus)
	t.base.Response(tx, err, w, items)
}

func (t *Api) workflowStepAvailableStatus(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	workflow := &domain.Workflow{Id: t.base.ParamValue("workflow", c, r)}
	prevStatus := &domain.Status{Id: t.base.ParamValue("status", c, r)}
	items, err := t.base.Database.StatusDb.WorkflowStepAvailableStatus(tx, workflow, prevStatus)
	t.base.Response(tx, err, w, items)
}

func (t *Api) workflowStepList(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.StatusDb.WorkflowSteps(tx, &domain.Workflow{Id: t.base.ParamValue("workflow", c, r)})
	t.base.Response(tx, err, w, items)
}

func (t *Api) workflowStepMembers(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.StatusDb.WorkflowStepMembers(tx, &domain.WorkflowStep{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, items)
}

func (t *Api) workflowStepSave(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.WorkflowStep{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	if len(item.Id) != 0 {
		err = t.base.Database.StatusDb.WorkflowStepUpdate(tx, item)
	} else {
		err = t.base.Database.StatusDb.WorkflowStepCreate(tx, item)
	}
	t.base.Response(tx, err, w, item)
}

func (t *Api) workflowStepRemove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.StatusDb.WorkflowStepRemove(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) workflowStepMemberGroups(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	selected, unselected, err := t.base.Database.StatusDb.WorkflowStepMemberGroups(tx, &domain.WorkflowStep{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, &workflowStepGroups{Selected: selected, Unselected: unselected})
}

func (t *Api) workflowStepMemberUsers(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	selected, unselected, err := t.base.Database.StatusDb.WorkflowStepMemberUsers(tx, &domain.WorkflowStep{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, &workflowStepUsers{Selected: selected, Unselected: unselected})
}

func (t *Api) workflowStepMemberAdd(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.WorkflowStepMember{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

	err = t.base.Database.StatusDb.WorkflowStepMemberAdd(tx, item)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	if len(item.IdGroup.String) != 0 {
		selected, unselected, err := t.base.Database.StatusDb.WorkflowStepMemberGroups(tx, item.WorkflowStep)
		t.base.Response(tx, err, w, &workflowStepGroups{Selected: selected, Unselected: unselected})
	} else if len(item.IdUser.String) != 0 {
		selected, unselected, err := t.base.Database.StatusDb.WorkflowStepMemberUsers(tx, item.WorkflowStep)
		t.base.Response(tx, err, w, &workflowStepUsers{Selected: selected, Unselected: unselected})
	} else {
		t.base.ErrorResponse(tx, fmt.Errorf("invalid workflow step member"), w)
	}
}

func (t *Api) workflowStepMemberRemove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.WorkflowStepMember{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

	err = t.base.Database.StatusDb.WorkflowStepMemberRemove(tx, item)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	if len(item.IdGroup.String) != 0 {
		selected, unselected, err := t.base.Database.StatusDb.WorkflowStepMemberGroups(tx, item.WorkflowStep)
		t.base.Response(tx, err, w, &workflowStepGroups{Selected: selected, Unselected: unselected})
	} else if len(item.IdUser.String) != 0 {
		selected, unselected, err := t.base.Database.StatusDb.WorkflowStepMemberUsers(tx, item.WorkflowStep)
		t.base.Response(tx, err, w, &workflowStepUsers{Selected: selected, Unselected: unselected})
	} else {
		t.base.ErrorResponse(tx, fmt.Errorf("invalid workflow step member"), w)
	}
}
