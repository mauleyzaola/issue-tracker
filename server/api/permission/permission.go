package permission

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
	item, err := t.base.Database.PermissionDb.Load(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) save(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)

	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

	item := &domain.PermissionScheme{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	if len(item.Id) != 0 {
		err = t.base.Database.PermissionDb.Update(tx, item)
	} else {
		err = t.base.Database.PermissionDb.Create(tx, item)
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
	item, err := t.base.Database.PermissionDb.Remove(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) names(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)

	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.PermissionDb.Names(tx)
	for i := range items {
		item := &items[i]
		item.Initialize()
	}
	t.base.Response(tx, err, w, items)
}

func (t *Api) clearMembers(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)

	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	err = t.base.Database.PermissionDb.ClearAll(tx, &domain.PermissionScheme{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, nil)
}

func (t *Api) list(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)

	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.PermissionDb.List(tx)
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
	err = t.base.Database.PermissionDb.Grid(tx, grid)
	t.base.Response(tx, err, w, grid)
}

func (t *Api) items(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)

	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.PermissionDb.Items(tx, &domain.PermissionScheme{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, items)
}

func (t *Api) schemeProjects(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)

	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.PermissionDb.Projects(tx, &domain.PermissionScheme{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, items)
}

func (t *Api) availablesUser(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)

	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	type psIssue struct {
		User  *domain.User
		Issue *domain.Issue
	}
	item := &psIssue{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.PermissionDb.AvailablesUser(tx, item.User, item.Issue)
	t.base.Response(tx, err, w, items)
}

func (t *Api) itemAdd(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)

	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.PermissionSchemeItem{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	err = t.base.Database.PermissionDb.ItemAdd(tx, item)
	item.Initialize()
	t.base.Response(tx, err, w, item)
}

func (t *Api) itemRemove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)

	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.PermissionSchemeItem{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	err = t.base.Database.PermissionDb.ItemRemove(tx, item)
	item.Initialize()
	t.base.Response(tx, err, w, item)
}
