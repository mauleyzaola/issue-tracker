package permission

import (
	"testing"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestPermissionSchemeCrud(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a new PermissionScheme")

		scheme := mock.PermissionScheme()
		err := app.Db.PermissionDb.Create(tx, scheme)
		assert.Nil(t, err)
		assert.NotEmpty(t, scheme.Id)

		t.Log("Load the new item")
		scheme2, err := app.Db.PermissionDb.Load(tx, scheme.Id)
		assert.Nil(t, err)
		assert.Equal(t, scheme2.Id, scheme.Id)
		assert.NotNil(t, scheme2.Meta)
		assert.NotEmpty(t, scheme2.Meta.DocumentType)
		assert.NotEmpty(t, scheme2.Meta.FriendName)

		t.Log("Update the item")
		scheme2.Name = scheme.Name + "X"
		err = app.Db.PermissionDb.Update(tx, scheme2)
		assert.Nil(t, err)

		scheme3, err := app.Db.PermissionDb.Load(tx, scheme.Id)
		assert.Nil(t, err)
		assert.NotNil(t, scheme3)
		assert.Equal(t, scheme3.Id, scheme2.Id)
		assert.NotEqual(t, scheme3.Name, scheme.Name)

		t.Log("Remove the item and force error")
		scheme4, err := app.Db.PermissionDb.Remove(tx, scheme.Id)
		assert.Nil(t, err)
		assert.NotNil(t, scheme4)
		assert.NotEmpty(t, scheme4.Id)

		deleted, err := app.Db.PermissionDb.Load(tx, scheme4.Id)
		assert.NotNil(t, err)
		assert.Nil(t, deleted)
	})
}

func TestPermissionSchemeList(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a new list of permission schemes")
		list, err := app.Db.PermissionDb.List(tx)
		assert.Nil(t, err)
		oldLength := len(list)

		counter := 10
		for i := 0; i < counter; i++ {
			scheme := mock.PermissionScheme()
			err = app.Db.PermissionDb.Create(tx, scheme)
			assert.Nil(t, err)
			assert.NotEmpty(t, scheme.Id)
		}

		list, err = app.Db.PermissionDb.List(tx)
		t.Log("Compare the number with the original items")
		assert.Nil(t, err)
		assert.Equal(t, len(list), oldLength+counter)
	})
}

func TestPermissionSchemeProject(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a new permission scheme")
		admin := mock.User()
		admin.IsSystemAdministrator = true
		err := mock.UserCreate(app.Db, tx, admin)
		assert.Nil(t, err)
		assert.NotEmpty(t, admin.Id)
		session, err := mock.SessionCreate(app.Db, tx, admin)
		assert.Nil(t, err)
		assert.NotNil(t, session)
		app.Db.Db.SetCurrentSession(session)

		scheme := mock.PermissionScheme()
		err = mock.PermissionSchemeCreate(app.Db, tx, scheme)
		assert.Nil(t, err)
		assert.NotEmpty(t, scheme.Id)

		oldProjectList, err := app.Db.PermissionDb.Projects(tx, scheme)
		assert.Nil(t, err)
		assert.NotNil(t, oldProjectList)

		counter := 10

		t.Log("Create relation between permission scheme and project")
		for i := 0; i < counter; i++ {
			project := mock.Project(i)
			project.PermissionScheme = scheme
			err = mock.ProjectCreate(app.Db, tx, project)
			assert.Nil(t, err)
		}

		t.Log("Get list of projects related with this permission scheme")
		newProjectList, err := app.Db.PermissionDb.Projects(tx, scheme)
		assert.Nil(t, err)
		assert.NotEqual(t, len(newProjectList), len(oldProjectList))
		assert.Equal(t, len(newProjectList), counter)

		t.Log("Clear all the permissions from the scheme and validate they are not in db")

		t.Log("Remove the permission scheme for each project created")
		for _, pl := range newProjectList {
			p, e := app.Db.ProjectDb.Load(tx, pl.Id)
			assert.Nil(t, e)
			assert.NotNil(t, p)
			assert.Equal(t, p.Id, pl.Id)

			p.PermissionScheme = nil
			p.IdPermissionScheme.String = ""
			p.IdPermissionScheme.Valid = false
			err = app.Db.ProjectDb.Update(tx, p)
			assert.Nil(t, err)
		}

		newProjectList, err = app.Db.PermissionDb.Projects(tx, scheme)
		assert.Nil(t, err)
		assert.Equal(t, len(newProjectList), 0)
	})
}
