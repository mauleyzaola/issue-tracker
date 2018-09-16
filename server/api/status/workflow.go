package status

import (
	"net/http"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
	"github.com/zenazn/goji/web"
)

func (t *Api) workflowGrid(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	grid := tecgrid.ParseQueryString(r.URL.Query())
	err = t.base.Database.StatusDb.WorkflowGrid(tx, grid)
	t.base.Response(tx, err, w, grid)
}

func (t *Api) workflowLoad(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.StatusDb.WorkflowLoad(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) workflowCreateMeta(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.StatusDb.WorkflowCreateMeta(tx, &domain.Workflow{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, item)
}

func (t *Api) workflowList(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.StatusDb.WorkflowList(tx)
	t.base.Response(tx, err, w, items)
}

func (t *Api) workflowSave(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.Workflow{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

	if len(item.Id) != 0 {
		err = t.base.Database.StatusDb.WorkflowUpdate(tx, item)
	} else {
		err = t.base.Database.StatusDb.WorkflowCreate(tx, item)
	}
	t.base.Response(tx, err, w, item)
}

func (t *Api) workflowRemove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.StatusDb.WorkflowRemove(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}
