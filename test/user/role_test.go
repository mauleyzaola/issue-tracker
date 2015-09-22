package user

import (
	"testing"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
	"github.com/mauleyzaola/tecutils"
	"github.com/stretchr/testify/assert"
)

func TestRoleCrud(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a role")

		t.Log("Create it on database")
		role := mock.Role()
		err := mock.RoleCreate(app.Db, tx, role)
		assert.Nil(t, err)
		assert.NotEmpty(t, role.Id)
		assert.NotNil(t, role.Meta)
		assert.NotEmpty(t, role.Meta.DocumentType)
		assert.NotEmpty(t, role.Meta.FriendName)

		t.Log("Load the new created role")
		role2, err := app.Db.UserDb.RoleLoad(tx, role.Id)
		assert.Nil(t, err)
		assert.Equal(t, role2.Name, role.Name)
		assert.NotNil(t, role2.Meta)
		assert.NotEmpty(t, role2.Meta.DocumentType)
		assert.NotEmpty(t, role2.Meta.FriendName)

		t.Log("Modify its properties and update")
		role2.Name = tecutils.UUID()
		err = app.Db.UserDb.RoleUpdate(tx, role2)
		assert.Nil(t, err)

		t.Log("Load again and compare with original")
		role3, err := app.Db.UserDb.RoleLoad(tx, role2.Id)
		assert.Nil(t, err)
		assert.Equal(t, role3.Name, role2.Name)
		assert.NotEqual(t, role3.Name, role.Name)

		t.Log("Delete the group")
		role4, err := app.Db.UserDb.RoleRemove(tx, role3.Id)
		assert.Nil(t, err)
		assert.NotNil(t, role4)

		t.Log("Load again and force error")
		role5, err := app.Db.UserDb.RoleLoad(tx, role4.Id)
		assert.NotNil(t, err)
		assert.Nil(t, role5)
	})
}

func TestRoleList(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a list of roles")

		t.Log("Get the current list of roles")
		oldItems, err := app.Db.UserDb.RoleList(tx)
		oldCount := len(oldItems)

		t.Log("Create the list and get the number")
		for i := 0; i < 10; i++ {
			newRole := mock.Role()
			err = mock.RoleCreate(app.Db, tx, newRole)
			assert.Nil(t, err)
			assert.NotEmpty(t, newRole.Id)
		}
		t.Log("Validate the number should be different than original list")
		newItems, err := app.Db.UserDb.RoleList(tx)
		newItemCount := len(newItems)
		assert.Nil(t, err)
		assert.Equal(t, oldCount+10, newItemCount)
	})
}
