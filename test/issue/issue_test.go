package issue

import (
	"testing"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
	"github.com/mauleyzaola/tecgrid"
	"github.com/stretchr/testify/assert"
)

func TestIssueCrud(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		session, err := mock.SessionSetContext(app.Db, tx, true)
		assert.Nil(t, err)
		assert.NotNil(t, session)

		item := mock.Issue()
		err = mock.IssueCreate(app.Db, tx, item)
		if !assert.Nil(t, err) {
			t.Log(err)
		}

		item2, err := app.Db.IssueDb.Remove(tx, item.Id)
		assert.Nil(t, err)
		assert.NotNil(t, item2)
	})
}

func TestIssueChangeStatus(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		session, err := mock.SessionSetContext(app.Db, tx, true)
		assert.Nil(t, err)
		assert.NotNil(t, session)

		item := mock.Issue()
		err = mock.IssueCreate(app.Db, tx, item)
		if !assert.Nil(t, err) {
			t.Log(err)
		}

		item, err = app.Db.IssueDb.Load(tx, item.Id, "")
		assert.Nil(t, err)
		assert.NotNil(t, item)
		assert.NotEmpty(t, item.Id)

		steps, err := app.Db.StatusDb.WorkflowStepAvailableUser(tx, item.Workflow, item.Status)
		assert.Nil(t, err)
		assert.Equal(t, true, len(steps) > 0)

		nextStatus := steps[0]
		err = app.Db.IssueDb.StatusChange(tx, item, nextStatus.NextStatus, nil)
		assert.Nil(t, err)

		item2, err := app.Db.IssueDb.Load(tx, item.Id, "")
		assert.Nil(t, err)
		assert.NotNil(t, item2)
		assert.NotEqual(t, item.Status.Id, item2.Status.Id)
		t.Log(item.Status.Id, item2.Status.Id)
		t.Log(item.Status.Name, item2.Status.Name)
	})
}

func TestIssueMoveProject(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("given a two projects")
		session, err := mock.SessionSetContext(app.Db, tx, true)
		assert.Nil(t, err)
		assert.NotNil(t, session)

		p1 := mock.Project(1)
		p2 := mock.Project(1)
		p2.Pkey = "XXX"
		err = mock.ProjectCreate(app.Db, tx, p1)
		assert.Nil(t, err)
		assert.NotEmpty(t, p1.Id)

		err = mock.ProjectCreate(app.Db, tx, p2)
		assert.Nil(t, err)
		assert.NotEmpty(t, p2.Id)

		t.Log("create an issue on the first project")
		issue := mock.Issue()
		issue.Project = p1
		err = mock.IssueCreate(app.Db, tx, issue)
		assert.Nil(t, err)
		assert.NotEmpty(t, issue.Id)
		assert.NotEmpty(t, issue.Pkey)

		oldkey := issue.Pkey

		t.Log("move the issue to the second project")
		_, err = app.Db.IssueDb.MoveProject(tx, issue, p2)
		assert.Nil(t, err)

		t.Log("validate the issue's key matches its new project")
		issue, err = app.Db.IssueDb.Load(tx, issue.Id, issue.Pkey)
		assert.Nil(t, err)
		assert.NotNil(t, issue)
		assert.NotEqual(t, oldkey, issue.Pkey)

		t.Log("list all the issues for the second project and validate issue is there")
		grid := &tecgrid.NgGrid{}
		grid.PageNumber = 1
		grid.PageSize = 10
		filter := &database.IssueFilter{}
		filter.Project = p2
		err = app.Db.IssueDb.Grid(tx, grid, filter)
		assert.Nil(t, err)

		rows, ok := grid.Rows.(*[]database.IssueGrid)
		assert.Equal(t, true, ok)
		assert.NotNil(t, rows)
		assert.Equal(t, true, len(*rows) == 1)
	})
}
