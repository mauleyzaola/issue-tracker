package issue

import (
	"testing"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
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
