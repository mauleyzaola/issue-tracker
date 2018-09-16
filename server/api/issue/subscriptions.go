package issue

import (
	"net/http"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
	"github.com/zenazn/goji/web"
)

func (t *Api) subscriptionToogle(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	selected, err := t.base.Database.IssueDb.SubscriptionToggle(tx, &domain.Issue{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, &struct {
		Selected bool `json:"selected"`
	}{
		selected,
	})
}

func (t *Api) subscriptionLoad(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	ok, err := t.base.Database.IssueDb.IsSubscribed(tx, &domain.Issue{Id: t.base.ParamValue("id", c, r)}, t.base.CurrentSession(c).User)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	t.base.Response(tx, err, w, &struct {
		Subscribed bool `json:"subscribed"`
	}{
		ok,
	})
}

func (t *Api) subscriptions(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	selected, unselected, err := t.base.Database.IssueDb.SubscribersIssue(tx, &domain.Issue{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, &struct {
		Selected   []domain.User `json:"selected"`
		Unselected []domain.User `json:"unselected"`
	}{
		selected,
		unselected,
	})
}

func (t *Api) subscriptionToogleUser(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.IssueSubscription{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	selected, err := t.base.Database.IssueDb.SubscriptionToggleUser(tx, item.Issue, item.User)
	t.base.Response(tx, err, w, &struct {
		Selected bool `json:"selected"`
	}{
		selected,
	})
}

func (t *Api) mySuscriptions(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	grid := tecgrid.ParseQueryString(r.URL.Query())
	err = t.base.Database.IssueDb.MySubscriptions(tx, grid)
	t.base.Response(tx, err, w, grid)
}
