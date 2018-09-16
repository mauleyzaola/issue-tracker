package issue

import (
	"net/http"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
	"github.com/zenazn/goji/web"
)

func (t *Api) status(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	issue := &domain.Issue{Id: t.base.ParamValue("id", c, r), Pkey: t.base.ParamValue("pkey", c, r)}
	status := &domain.Status{Id: t.base.ParamValue("status", c, r)}
	err = t.base.Database.IssueDb.StatusChange(tx, issue, status, nil)
	t.base.Response(tx, err, w, issue.Status)
}

func (t *Api) load(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.IssueDb.Load(tx, t.base.ParamValue("id", c, r), t.base.ParamValue("pkey", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) createMeta(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

	type Meta struct {
		Item         *domain.Issue         `json:"item"`
		Priorities   []domain.Priority     `json:"priorities"`
		Users        []domain.User         `json:"users"`
		Workflows    []domain.Workflow     `json:"workflows"`
		Steps        []domain.WorkflowStep `json:"steps"`
		Comments     []domain.IssueComment `json:"comments"`
		IsSubscribed bool                  `json:"isSubscribed"`
		Parent       *domain.Issue         `json:"parent"`
	}
	item := &Meta{}
	item.Item = &domain.Issue{}
	if value := c.URLParams["pkey"]; len(value) != 0 {
		item.Item, err = t.base.Database.IssueDb.Load(tx, "", value)
		if err != nil {
			t.base.ErrorResponse(tx, err, w)
			return
		}

		steps, err := t.base.Database.StatusDb.WorkflowStepAvailableUser(tx, item.Item.Workflow, item.Item.Status)
		if err != nil {
			t.base.ErrorResponse(tx, err, w)
			return
		}
		item.Steps = steps

		comments, err := t.base.Database.IssueDb.CommentList(tx, item.Item)
		if err != nil {
			t.base.ErrorResponse(tx, err, w)
			return
		}
		item.Comments = comments

		subscribed, err := t.base.Database.IssueDb.IsSubscribed(tx, item.Item, t.base.CurrentSession(c).User)
		if err != nil {
			t.base.ErrorResponse(tx, err, w)
			return
		}
		item.IsSubscribed = subscribed

		//if the issue has a parent, load it here instead of the dal to avoid possible stack overflow
		if item.Item.IdParent.Valid {
			parent, e := t.base.Database.IssueDb.Load(tx, item.Item.IdParent.String, "")
			if e != nil {
				t.base.ErrorResponse(tx, e, w)
			} else {
				item.Parent = parent
			}
		}

	} else {
		item.Item.DateCreated = time.Now()
		item.Item.DueDate = item.Item.DateCreated
		item.Workflows, err = t.base.Database.StatusDb.WorkflowList(tx)
		if err != nil {
			t.base.ErrorResponse(tx, err, w)
			return
		}
	}

	item.Priorities, err = t.base.Database.PriorityDb.List(tx)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item.Users, err = t.base.Database.UserDb.List(tx)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

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
	filter := &database.IssueFilter{}
	if val := t.base.ParamValue("status", c, r); len(val) != 0 {
		filter.Status = &domain.Status{Name: val}
	}
	if val := t.base.ParamValue("priority", c, r); len(val) != 0 {
		filter.Priority = &domain.Priority{Id: val}
	}
	if val := t.base.ParamValue("assignee", c, r); len(val) != 0 {
		filter.Assignee = &domain.User{Id: val}
	}
	if val := t.base.ParamValue("reporter", c, r); len(val) != 0 {
		filter.Reporter = &domain.User{Id: val}
	}
	if val := t.base.ParamValue("project", c, r); len(val) != 0 {
		filter.Project = &domain.Project{Id: val}
	}
	if val := t.base.ParamValue("parent", c, r); len(val) != 0 {
		filter.Parent = &domain.Issue{Id: val}
	}
	if val := t.base.ParamValue("due", c, r); len(val) != 0 {
		filter.Due.Valid = true
		if val == "true" {
			filter.Due.Bool = true
		} else {
			filter.Due.Bool = false
		}
	}
	if val := t.base.ParamValue("resolved", c, r); len(val) != 0 {
		filter.Resolved.Valid = true
		if val == "true" {
			filter.Resolved.Bool = true
		} else {
			filter.Resolved.Bool = false
		}
	}

	if val := t.base.ParamValue("dueDate", c, r); len(val) != 0 {
		t, e := time.Parse(operations.DATE_FORMAT_JSON, val)
		if e == nil {
			filter.DueDate = &t
		}
	}

	err = t.base.Database.IssueDb.Grid(tx, grid, filter)
	t.base.Response(tx, err, w, grid)
}

func (t *Api) children(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	items, err := t.base.Database.IssueDb.Children(tx, &domain.Issue{Id: t.base.ParamValue("id", c, r)})
	t.base.Response(tx, err, w, items)
}

func (t *Api) save(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.Issue{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	if len(item.Id) != 0 {
		err = t.base.Database.IssueDb.Update(tx, item)
	} else {
		err = t.base.Database.IssueDb.Create(tx, item, t.base.ParamValue("parent", c, r))
	}

	t.base.Response(tx, err, w, item)
}

func (t *Api) assignToMe(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.IssueDb.Load(tx, t.base.ParamValue("id", c, r), t.base.ParamValue("pkey", c, r))
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item.Assignee = t.base.CurrentSession(c).User
	err = t.base.Database.IssueDb.Update(tx, item)
	t.base.Response(tx, err, w, item)
}

func (t *Api) reporterIsMe(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.IssueDb.Load(tx, t.base.ParamValue("id", c, r), t.base.ParamValue("pkey", c, r))
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item.Reporter = t.base.CurrentSession(c).User
	err = t.base.Database.IssueDb.Update(tx, item)
	t.base.Response(tx, err, w, item)
}

func (t *Api) remove(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err := t.base.Database.IssueDb.Remove(tx, t.base.ParamValue("id", c, r))
	t.base.Response(tx, err, w, item)
}

func (t *Api) move(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item := &domain.Issue{}
	err = t.base.Decode(item, r)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	item, err = t.base.Database.IssueDb.MoveProject(tx, item, item.Project)
	t.base.Response(tx, err, w, item)
}
