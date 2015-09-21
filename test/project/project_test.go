package project

import (
	"strings"
	"testing"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestProjectCrud(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("create a new project")
		session, err := mock.SessionSetContext(tx, app.Db, true)
		assert.Nil(t, err)
		assert.NotNil(t, session)

		item := mock.Project(555)
		err = app.Db.ProjectDb.Create(tx, item)
		assert.Nil(t, err)
		assert.NotNil(t, item)
		assert.NotEmpty(t, item.Id)
		assert.Equal(t, item.ProjectLead.Id, session.User.Id)
		assert.Equal(t, item.DateCreated.Year(), time.Now().Year())

		t.Log("load project and update")
		item2, err := app.Db.ProjectDb.Load(tx, item.Id)
		assert.Nil(t, err)
		assert.NotNil(t, item2)
		assert.Equal(t, item2.DateCreated.Year(), time.Now().Year())
		assert.Nil(t, item.LastModified)

		item2.Begins = &time.Time{}
		*item2.Begins, err = time.Parse("2006-01-02", "2015-09-15")
		assert.Nil(t, err)
		err = app.Db.ProjectDb.Update(tx, item2)
		assert.Nil(t, err)

		t.Log("reload updated project")
		item3, err := app.Db.ProjectDb.Load(tx, item2.Id)
		assert.NotNil(t, item3)
		assert.NotNil(t, item3.LastModified)
		assert.Equal(t, item3.Begins.Year(), 2015)

		t.Log("remove project")
		item4, err := app.Db.ProjectDb.Remove(tx, item3.Id)
		assert.Nil(t, err)
		assert.NotNil(t, item4)

		t.Log("reload removed project and force error")
		item5, err := app.Db.ProjectDb.Load(tx, item3.Id)
		assert.NotNil(t, err)
		assert.Nil(t, item5)
	})
}

func TestProjectDups(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("create a new project")
		session, err := mock.SessionSetContext(tx, app.Db, true)
		assert.Nil(t, err)
		assert.NotNil(t, session)

		item := mock.Project(1)
		err = mock.ProjectCreate(tx, app.Db, item)
		assert.Nil(t, err)
		assert.NotNil(t, item)
		assert.NotEmpty(t, item.Id)

		item2 := mock.Project(2)
		item2.Name = strings.ToUpper(item.Name)
		err = mock.ProjectCreate(tx, app.Db, item2)
		if assert.NotNil(t, err) {
			t.Log(err)
		}
	})
}

func TestProjectDups2(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("create a new project")
		session, err := mock.SessionSetContext(tx, app.Db, true)
		assert.Nil(t, err)
		assert.NotNil(t, session)

		item := mock.Project(1)
		err = mock.ProjectCreate(tx, app.Db, item)
		assert.Nil(t, err)
		assert.NotNil(t, item)
		assert.NotEmpty(t, item.Id)

		t.Log("create a second project")
		item2 := mock.Project(2)
		err = mock.ProjectCreate(tx, app.Db, item2)
		assert.Nil(t, err)
		assert.NotNil(t, item2)

		t.Log("update the second project with the name of first one and produce error")
		item2.Name = item.Name
		err = app.Db.ProjectDb.Update(tx, item2)
		if assert.NotNil(t, err) {
			t.Log(err)
		}
	})
}
