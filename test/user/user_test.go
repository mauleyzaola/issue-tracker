package user

import (
	"testing"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
	"github.com/mauleyzaola/issue-tracker/utils/tecutils"
	"github.com/stretchr/testify/assert"
)

func TestCountSystemAdmins(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a number of user administrators")

		t.Log("Get the list of current admins")
		count, err := app.Db.UserDb.CountSystemAdministrators(tx)
		assert.Nil(t, err)

		t.Log("Add another administrator")
		newAdmin := mock.User()
		newAdmin.IsActive = true
		newAdmin.IsSystemAdministrator = true
		err = app.Db.UserDb.Create(tx, newAdmin)
		assert.Nil(t, err)
		assert.NotEmpty(t, newAdmin.Id)

		t.Log("Compare the list, should be another one on the list")
		count2, err := app.Db.UserDb.CountSystemAdministrators(tx)
		assert.Nil(t, err)
		assert.Equal(t, count2-1, count)
	})
}

func TestUserCrud(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given an user, make several crud tests")

		t.Log("Create a new user")
		user := mock.User()
		user.IsSystemAdministrator = true
		user.Password = "admin"
		err := mock.UserCreate(app.Db, tx, user)
		assert.Nil(t, err)
		assert.NotEmpty(t, user.Id)
		assert.NotNil(t, user.Meta)
		assert.NotEmpty(t, user.Meta.DocumentType)
		assert.NotEmpty(t, user.Meta.FriendName)
		assert.Equal(t, user.DateCreated.Year(), time.Now().Year())

		t.Log("Retrieve the user's data")
		user2, err := app.Db.UserDb.Load(tx, user.Id)
		assert.Nil(t, err)
		assert.NotNil(t, user2)
		assert.Equal(t, user2.Email, user.Email)

		t.Log("Update the user's data")
		user2.Name = tecutils.UUID()
		err = app.Db.UserDb.Update(tx, user2)
		assert.Nil(t, err)
		assert.NotNil(t, user2.LastModified)
		assert.Equal(t, user2.LastModified.Year(), time.Now().Year())

		user3, err := app.Db.UserDb.Load(tx, user.Id)
		assert.Nil(t, err)
		assert.NotNil(t, user3)
		assert.Equal(t, user3.Id, user.Id)
		assert.NotEqual(t, user3.Name, user.Name)

		t.Log("Create a session manually")
		session, err := app.Db.AccountDb.LoginUser(tx, user.Email, "admin", "127.0.0.1", true)
		assert.Nil(t, err)
		assert.NotNil(t, session)
		app.Db.Db.SetCurrentSession(session)

		t.Log("Remove the user from database with the same session should raise error")
		user4, err := app.Db.UserDb.Remove(tx, user3.Id)
		assert.NotNil(t, err)
		assert.Nil(t, user4)

		t.Log("Create another user should enable removing the new one")
		another1 := mock.User()
		err = mock.UserCreate(app.Db, tx, another1)
		assert.Nil(t, err)
		assert.NotEmpty(t, another1.Id)
		deleted, err := app.Db.UserDb.Remove(tx, another1.Id)
		assert.Nil(t, err)
		assert.NotNil(t, deleted)
		app.Db.Db.SetCurrentSession(nil)
	})
}

func TestUserGrid(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a list of users")
		t.Log("Insert some users")
		items, err := app.Db.UserDb.List(tx)
		assert.Nil(t, err)
		count := len(items)
		var name, email string
		for i := 0; i < 15; i++ {
			user := mock.User()
			err = mock.UserCreate(app.Db, tx, user)
			assert.Nil(t, err)
			if i == 0 {
				name = user.Name
				email = user.Email
			}
		}
		t.Log("Get a list of users and compare the number of items")
		items, err = app.Db.UserDb.List(tx)
		assert.Equal(t, count+15, len(items))

		t.Log("Get a grid of users")
		grid := &tecgrid.NgGrid{}
		grid.PageNumber = 3
		grid.PageSize = 5
		err = app.Db.UserDb.Grid(tx, grid)
		assert.Nil(t, err)
		assert.NotNil(t, grid.Rows)

		t.Log("filter grid results")
		grid.Query = name
		grid.PageNumber = 1
		err = app.Db.UserDb.Grid(tx, grid)
		assert.Nil(t, err)
		assert.NotNil(t, grid.Rows)
		rows, ok := grid.Rows.(*[]domain.User)
		assert.Equal(t, true, ok)
		assert.NotNil(t, rows)
		assert.EqualValues(t, 1, len(*rows))

		grid.Query = email
		err = app.Db.UserDb.Grid(tx, grid)
		assert.Nil(t, err)
		assert.NotNil(t, grid.Rows)
		rows, ok = grid.Rows.(*[]domain.User)
		assert.Equal(t, true, ok)
		assert.NotNil(t, rows)
		assert.EqualValues(t, 1, len(*rows))
	})
}
