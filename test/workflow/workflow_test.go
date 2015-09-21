package workflow

import (
	"testing"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestCrud(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a new workflow, perform some crud operations")
		item := mock.Workflow()

		t.Log("Create new workflow")
		err := mock.WorkflowCreate(tx, app.Db, item)
		assert.Nil(t, err)
		assert.NotEmpty(t, item.Id)

		t.Log("Create a second workflow with same name to force error")
		item2 := &domain.Workflow{Name: item.Name}
		err = mock.WorkflowCreate(tx, app.Db, item2)

		assert.NotNil(t, err)
		t.Log(err)
		app.Db.Db.Rollback(tx)
		tx, err = app.Db.Db.Begin()
		assert.Nil(t, err)
		assert.NotNil(t, tx)

		t.Log("Update original item, reload and compare changes")
		err = mock.WorkflowCreate(tx, app.Db, item)
		assert.Nil(t, err)
		assert.NotEmpty(t, item.Id)

		item3, err := app.Db.StatusDb.WorkflowLoad(tx, item.Id)
		assert.Nil(t, err)
		assert.NotNil(t, item3)

		item3.Name += "x"
		err = app.Db.StatusDb.WorkflowUpdate(tx, item3)
		assert.Nil(t, err)
		assert.NotNil(t, item3.LastModified)

		item4, err := app.Db.StatusDb.WorkflowLoad(tx, item3.Id)
		assert.Nil(t, err)
		assert.NotNil(t, item4)
		assert.NotEqual(t, item4.Name, item.Name)

		t.Log("Remove item, reload and force error")
		item5, err := app.Db.StatusDb.WorkflowRemove(tx, item4.Id)
		assert.Nil(t, err)
		assert.NotNil(t, item5)
		assert.Equal(t, item5.Name, item4.Name)

		removed, err := app.Db.StatusDb.WorkflowLoad(tx, item5.Id)
		assert.NotNil(t, err)
		assert.Nil(t, removed)
	})
}

func TestWorkflowStep(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a set of workflow/workflow step, perform basic crud operations")
		workflow := mock.Workflow()
		err := mock.WorkflowCreate(tx, app.Db, workflow)
		assert.Nil(t, err)
		assert.NotEmpty(t, workflow.Id)

		stNew := &domain.Status{Name: "New", Workflow: workflow}
		stInProgress := &domain.Status{Name: "In Progress", Workflow: workflow}
		stResolved := &domain.Status{Name: "Resolved", Workflow: workflow}
		stCancelled := &domain.Status{Name: "Cancelled", Workflow: workflow}

		open := &domain.WorkflowStep{Workflow: workflow, Name: "Open"}
		open.NextStatus = stInProgress

		beginProgress := &domain.WorkflowStep{Workflow: workflow, Name: "Begin Progress"}
		beginProgress.PrevStatus = stNew
		beginProgress.NextStatus = stInProgress

		resolve := &domain.WorkflowStep{Workflow: workflow, Name: "Resolve"}
		resolve.PrevStatus = stInProgress
		resolve.NextStatus = stResolved
		resolve.Resolves = true

		cancels := &domain.WorkflowStep{Workflow: workflow, Name: "Cancel"}
		cancels.PrevStatus = stInProgress
		cancels.NextStatus = stCancelled
		cancels.Cancels = true

		err = app.Db.StatusDb.Create(tx, stCancelled)
		assert.Nil(t, err)
		err = app.Db.StatusDb.Create(tx, stInProgress)
		assert.Nil(t, err)
		err = app.Db.StatusDb.Create(tx, stNew)
		assert.Nil(t, err)
		err = app.Db.StatusDb.Create(tx, stResolved)
		assert.Nil(t, err)

		err = app.Db.StatusDb.WorkflowStepCreate(tx, open)
		assert.Nil(t, err)
		err = app.Db.StatusDb.WorkflowStepCreate(tx, beginProgress)
		assert.Nil(t, err)
		err = app.Db.StatusDb.WorkflowStepCreate(tx, resolve)
		assert.Nil(t, err)
		err = app.Db.StatusDb.WorkflowStepCreate(tx, cancels)
		assert.Nil(t, err)

		t.Log("Break some workflow rules and see it fail")
		resolve.Cancels = true
		err = app.Db.StatusDb.WorkflowStepUpdate(tx, resolve)
		assert.NotNil(t, err)
		if err != nil {
			t.Log(err)
		}
		resolve.Cancels = false

		cancels.Resolves = true
		err = app.Db.StatusDb.WorkflowStepUpdate(tx, cancels)
		assert.NotNil(t, err)
		if err != nil {
			t.Log(err)
		}
	})
}
