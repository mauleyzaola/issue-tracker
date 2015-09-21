package issue

import (
	"net/http"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/zenazn/goji/web"
)

func (t *Api) comments(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.IssueDb.CommentList(tx, &domain.Issue{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, items)
}

func (t *Api) commentSave(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.IssueComment{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

	if len(item.Id) != 0 {
		err = t.base.Database.IssueDb.CommentUpdate(tx, item)
	} else {
		err = t.base.Database.IssueDb.CommentAdd(tx, item)
	}
	t.base.Response(tx, err, w, item)
}

func (t *Api) commentRemove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.IssueDb.CommentRemove(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}
