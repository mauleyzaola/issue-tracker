package project

import (
	"net/http"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/zenazn/goji/web"
)

func (t *Api) roleProjectMembers(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.Project{Id: t.base.ParamValue("id", c, r)}
	items, err := t.base.Database.ProjectDb.RoleProjectMembers(tx, item)
	t.base.Response(tx, err, w, items)
}

func (t *Api) projectRoleMembers(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.ProjectRole{Id: t.base.ParamValue("id", c, r)}
	items, err := t.base.Database.ProjectDb.RoleMembers(tx, item)
	t.base.Response(tx, err, w, items)
}

func (t *Api) projectRoleMemberAdd(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.ProjectRoleMember{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	err = t.base.Database.ProjectDb.RoleMemberAdd(tx, item.ProjectRole, item.User, item.Group)
	t.base.Response(tx, err, w, nil)
}

func (t *Api) projectRoleMemberRemove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.ProjectRoleMember{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	err = t.base.Database.ProjectDb.RoleMemberRemove(tx, item)
	t.base.Response(tx, err, w, nil)
}
