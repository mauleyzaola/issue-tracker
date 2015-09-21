package account

import (
	"net/http"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/zenazn/goji/web"
)

type LoginUser struct {
	Email    string
	Password string
}

func (t *Api) changeMyPassword(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &LoginUser{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	err = t.base.Database.AccountDb.ChangeMyPassword(tx, item.Password)
	t.base.Response(tx, err, w, nil)
}

func (t *Api) myProfileLoad(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

	item, err := t.base.Database.UserDb.Load(tx, t.base.CurrentSession(c).User.Id)
	t.base.Response(tx, err, w, item)
}

func (t *Api) myProfileSave(c web.C, w http.ResponseWriter, r *http.Request) {
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

	err = t.base.Database.UserDb.Update(tx, item)
	t.base.Response(tx, err, w, item)
}

func (t *Api) login(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &LoginUser{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	session, err := t.base.Database.AccountDb.LoginUser(tx, item.Email, item.Password, t.base.IpAddress(r), true)
	t.base.Response(tx, err, w, session)
}

func (t *Api) logout(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

	err = t.base.Database.AccountDb.LogoutUser(tx)
	t.base.Response(tx, err, w, nil)
}

func (t *Api) sessionList(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.SessionDb.List(tx, &domain.User{Id: t.base.CurrentSession(c).User.Id})
	t.base.Response(tx, err, w, items)
}

func (t *Api) sessionRemove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	err = t.base.Database.SessionDb.Remove(tx, t.base.ParamValue("id", c, r), true)
	t.base.Response(tx, err, w, nil)
}

func (t *Api) sessionDelay(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	//TODO: make this value parametrizable
	err = t.base.Database.AccountDb.LoginDelay(tx, time.Hour*24*3)
	t.base.Response(tx, err, w, nil)
}

func (t *Api) changePassword(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

	type loginUser struct {
		Id       string
		Email    string
		Password string
	}
	item := &loginUser{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	err = t.base.Database.AccountDb.ChangePassword(tx, item.Id, item.Password)
	t.base.Response(tx, err, w, nil)
}
