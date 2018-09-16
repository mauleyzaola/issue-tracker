package user

import (
	"net/http"

	"github.com/mauleyzaola/issue-tracker/server/domain"
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
	item, err := t.base.Database.UserDb.Load(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) remove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.UserDb.Remove(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) list(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.UserDb.List(tx)
	for i := range items {
		item := &items[i]
		item.Initialize()
	}
	t.base.Response(tx, err, w, items)
}

func (t *Api) grid(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	grid := tecgrid.ParseQueryString(r.URL.Query())
	err = t.base.Database.UserDb.Grid(tx, grid)
	t.base.Response(tx, err, w, grid)
}

func (t *Api) save(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.User{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	if len(item.Id) != 0 {
		err = t.base.Database.UserDb.Update(tx, item)
	} else {
		err = t.base.Database.UserDb.Create(tx, item)
	}
	t.base.Response(tx, err, w, item)
}
