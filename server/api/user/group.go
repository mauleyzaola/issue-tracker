package user

import (
	"net/http"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
	"github.com/zenazn/goji/web"
)

func (t *Api) groupLoad(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.UserDb.GroupLoad(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) groupRemove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.UserDb.GroupRemove(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) groupList(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.UserDb.GroupList(tx)
	for i := range items {
		item := &items[i]
		item.Initialize()
	}
	t.base.Response(tx, err, w, items)
}

func (t *Api) groupGrid(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	grid := tecgrid.ParseQueryString(r.URL.Query())
	err = t.base.Database.UserDb.GroupGrid(tx, grid)
	t.base.Response(tx, err, w, grid)
}

func (t *Api) groupSave(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.Group{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	if len(item.Id) != 0 {
		err = t.base.Database.UserDb.GroupUpdate(tx, item)
	} else {
		err = t.base.Database.UserDb.GroupCreate(tx, item)
	}
	t.base.Response(tx, err, w, item)
}

func (t *Api) groupUsers(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	selected, unselected, err := t.base.Database.UserDb.UserGroupListUsers(tx, &domain.Group{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, &struct {
		Selected   []domain.User `json:"selected"`
		Unselected []domain.User `json:"unselected"`
	}{
		selected,
		unselected,
	})
}

func (t *Api) userGroups(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	selected, unselected, err := t.base.Database.UserDb.UserGroupListGroups(tx, &domain.User{Id: t.base.ParamValue("id", c, r)})

	t.base.Response(tx, err, w, &struct {
		Selected   []domain.Group `json:"selected"`
		Unselected []domain.Group `json:"unselected"`
	}{
		selected,
		unselected,
	})
}

func (t *Api) userAdd(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.UserGroup{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	err = t.base.Database.UserDb.UserGroupAdd(tx, item)
	t.base.Response(tx, err, w, item)
}

func (t *Api) userRemove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.UserGroup{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	err = t.base.Database.UserDb.UserGroupRemove(tx, item)
	t.base.Response(tx, err, w, item)
}
