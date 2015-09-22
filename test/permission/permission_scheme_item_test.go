package permission

import (
	"testing"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestPermissionSchemeItemCrud(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a set of permissions")

		scheme := mock.PermissionScheme()
		err := app.Db.PermissionDb.Create(tx, scheme)
		assert.Nil(t, err)
		assert.NotEmpty(t, scheme.Id)

		t.Log("Add items")
		names, err := app.Db.PermissionDb.Names(tx)
		assert.Nil(t, err)
		assert.Equal(t, true, len(names) > 0)
		firstName := &names[0]

		user := mock.User()
		err = mock.UserCreate(app.Db, tx, user)
		assert.Nil(t, err)
		assert.NotEmpty(t, user.Id)

		item := mock.PermissionSchemeItem()
		item.User = user
		item.PermissionName = firstName
		item.PermissionScheme = scheme
		err = item.Validate()
		assert.Nil(t, err)

		mock.PermissionSchemeItemCreate(app.Db, tx, item)
		assert.Nil(t, err)
		assert.NotEmpty(t, item.Id)

		t.Log("Load added items")
		addedItems, err := app.Db.PermissionDb.Items(tx, scheme)
		assert.Nil(t, err)
		assert.Equal(t, len(addedItems), 1)

		t.Log("Remove items")
		err = app.Db.PermissionDb.ItemRemove(tx, item)
		assert.Nil(t, err)

		t.Log("Validate the items are not there anymore")
		addedItems, err = app.Db.PermissionDb.Items(tx, scheme)
		assert.Nil(t, err)
		assert.Equal(t, len(addedItems), 0)
	})
}
