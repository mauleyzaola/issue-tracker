package issue

import (
	"net/http"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/zenazn/goji/web"
)

func (t *Api) attachmentLoad(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.IssueDb.AttachmentLoad(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) attachmentAdd(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.IssueAttachment{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err = t.base.Database.IssueDb.AttachmentAdd(tx, item.Issue, item.FileItem)
	t.base.Response(tx, err, w, item)
}

func (t *Api) attachmentList(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.IssueDb.AttachmentList(tx, &domain.Issue{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, items)
}

func (t *Api) attachmentRemove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	err = t.base.Database.IssueDb.AttachmentRemove(tx, &domain.IssueAttachment{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, nil)
}
