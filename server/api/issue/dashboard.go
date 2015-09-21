package issue

import (
	"fmt"
	"net/http"

	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/zenazn/goji/web"
)

func (t *Api) dashboardGroupAll(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	var items interface{}
	dataType := t.base.ParamValue("dataType", c, r)
	switch dataType {
	case "assignee":
		items, err = t.base.Database.IssueDb.DashbordIssueGroupByAssignee(tx, nil)
	case "reporter":
		items, err = t.base.Database.IssueDb.DashbordIssueGroupByReporter(tx, nil)
	case "priority":
		items, err = t.base.Database.IssueDb.DashboardIssueGroupByPriority(tx, nil)
	case "status":
		items, err = t.base.Database.IssueDb.DashboardIssueGroupByStatus(tx, nil)
	case "dueDate":
		items, err = t.base.Database.IssueDb.DashbordIssueGroupByDueDate(tx, nil)
	case "project":
		items, err = t.base.Database.IssueDb.DashboardIssueGroupByProject(tx, nil)
	default:
		err = fmt.Errorf("Invalid dataType parameter")
		t.base.ErrorResponse(tx, err, w)
		return
	}
	t.base.Response(tx, err, w, items)
}

func (t *Api) dashboardGroupDataType(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	var items interface{}

	dataType := t.base.ParamValue("dataType", c, r)
	assignee := t.base.ParamValue("assignee", c, r)
	reporter := t.base.ParamValue("reporter", c, r)
	status := t.base.ParamValue("status", c, r)
	project := t.base.ParamValue("project", c, r)
	filter := &database.IssueGroupFilter{Assignee: assignee, Reporter: reporter, Project: project, Status: status}
	switch dataType {
	case "assignee":
		items, err = t.base.Database.IssueDb.DashbordIssueGroupByAssignee(tx, filter)
	case "reporter":
		items, err = t.base.Database.IssueDb.DashbordIssueGroupByReporter(tx, filter)
	case "priority":
		items, err = t.base.Database.IssueDb.DashboardIssueGroupByPriority(tx, filter)
	case "status":
		items, err = t.base.Database.IssueDb.DashboardIssueGroupByStatus(tx, filter)
	case "dueDate":
		items, err = t.base.Database.IssueDb.DashbordIssueGroupByDueDate(tx, filter)
	case "project":
		items, err = t.base.Database.IssueDb.DashboardIssueGroupByProject(tx, filter)
	default:
		err = fmt.Errorf("Invalid dataType parameter")
		t.base.ErrorResponse(tx, err, w)
		return
	}
	t.base.Response(tx, err, w, items)
}
