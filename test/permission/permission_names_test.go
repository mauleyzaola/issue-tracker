package permission

import (
	"testing"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestPermissionNames(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Retrieve the original list")
		names, err := app.Db.PermissionDb.Names(tx)
		assert.Nil(t, err)

		assert.Equal(t, true, len(names) > 0)
		for i := range names {
			name := &names[i]
			assert.NotNil(t, name.Meta)
			assert.NotEmpty(t, name.Meta.FriendName)
		}
	})
}

func TestPermissionAllowedUser(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a new user, group and a new set of permissions")

		names, err := app.Db.PermissionDb.Names(tx)
		assert.Nil(t, err)
		assert.NotEqual(t, len(names), 0)

		user := mock.User()
		err = mock.UserCreate(tx, app.Db, user)
		assert.Nil(t, err)
		assert.NotEmpty(t, user.Id)

		group := mock.Group()
		err = mock.GroupCreate(tx, app.Db, group)
		assert.Nil(t, err)
		assert.NotEmpty(t, group.Id)

		scheme := mock.PermissionScheme()
		err = mock.PermissionSchemeCreate(tx, app.Db, scheme)
		assert.Nil(t, err)
		assert.NotEmpty(t, scheme.Id)

		t.Log("Create user, group and assign some permissions")
		firstName := &names[0]
		itemWithUser := mock.PermissionSchemeItem()
		itemWithUser.User = user
		itemWithUser.PermissionScheme = scheme
		itemWithUser.PermissionName = firstName
		err = mock.PermissionSchemeItemCreate(tx, app.Db, itemWithUser)
		assert.Nil(t, err)
		assert.NotEmpty(t, itemWithUser.Id)

		itemWithGroup := mock.PermissionSchemeItem()
		itemWithGroup.Group = group
		itemWithGroup.PermissionName = firstName
		itemWithGroup.PermissionScheme = scheme
		err = mock.PermissionSchemeItemCreate(tx, app.Db, itemWithGroup)
		assert.Nil(t, err)
		assert.NotEmpty(t, itemWithGroup.Id)

		t.Log("Verify the permissions have been assigned to the user")
		ok, err := app.Db.PermissionDb.AllowedUser(tx, user, nil, firstName)
		assert.Nil(t, err)
		assert.Equal(t, ok, true)

		t.Log("Validate the user can or cannot access some permissions")
		allowedNames, err := app.Db.PermissionDb.AvailablesUser(tx, user, nil)
		assert.Nil(t, err)
		assert.NotNil(t, allowedNames)
		found := false
		for _, e := range allowedNames {
			if e.Id == firstName.Id {
				found = true
				break
			}
		}
		assert.Equal(t, found, true)
	})
}
