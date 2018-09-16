package project

import (
	"net/http"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
	"github.com/zenazn/goji/web"
)

func (t *Api) load(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.ProjectDb.Load(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) save(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)

	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

	item := &domain.Project{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	if len(item.Id) != 0 {
		err = t.base.Database.ProjectDb.Update(tx, item)
	} else {
		err = t.base.Database.ProjectDb.Create(tx, item)
	}

	t.base.Response(tx, err, w, item)
}

func (t *Api) remove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.ProjectDb.Remove(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) createMeta(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)

	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.ProjectDb.CreateMeta(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) grid(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	grid := tecgrid.ParseQueryString(r.URL.Query())
	filter := &database.ProjectFilter{}
	filter.ProjectLead = t.base.ParamValue("projectLead", c, r)
	if val := t.base.ParamValue("resolved", c, r); len(val) != 0 {
		filter.Resolved.Valid = true
		filter.Resolved.Bool = val == "true"
	}
	err = t.base.Database.ProjectDb.Grid(tx, grid, filter)
	t.base.Response(tx, err, w, grid)
}
