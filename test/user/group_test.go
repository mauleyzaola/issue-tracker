package user

import (
	"testing"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
	"github.com/mauleyzaola/tecutils"
	"github.com/stretchr/testify/assert"
)

func TestGroupCrud(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a group, test from crud operations")

		t.Log("Create new group")
		group := mock.Group()
		err := mock.GroupCreate(app.Db, tx, group)
		assert.Nil(t, err)
		assert.NotEmpty(t, group.Id)
		assert.Equal(t, group.DateCreated.Year(), time.Now().Year())

		t.Log("Modify properties and load again")
		group2, err := app.Db.UserDb.GroupLoad(tx, group.Id)
		assert.Nil(t, err)
		assert.NotNil(t, group2)
		assert.Equal(t, group2.Id, group.Id)
		assert.Equal(t, group2.Name, group.Name)
		group2.Name = tecutils.UUID()
		err = app.Db.UserDb.GroupUpdate(tx, group2)
		assert.Nil(t, err)
		assert.NotNil(t, group2.LastModified)
		assert.Equal(t, group2.LastModified.Year(), time.Now().Year())

		group3, err := app.Db.UserDb.GroupLoad(tx, group2.Id)
		assert.Nil(t, err)
		assert.NotNil(t, group3)
		assert.Equal(t, group3.Id, group.Id)
		assert.NotEqual(t, group3.Name, group.Name)

		t.Log("Delete group and trigger error")
		group4, err := app.Db.UserDb.GroupRemove(tx, group.Id)
		assert.Nil(t, err)
		assert.NotNil(t, group4)
		assert.Equal(t, group4.Name, group2.Name)
		nullGroup, err := app.Db.UserDb.GroupLoad(tx, group.Id)
		assert.NotNil(t, err)
		assert.Nil(t, nullGroup)
	})
}

func TestGroupMembers(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a group, modify its members")

		t.Log("Create a new group")
		group := mock.Group()
		err := mock.GroupCreate(app.Db, tx, group)
		assert.Nil(t, err)
		assert.NotEmpty(t, group.Id)
		assert.NotNil(t, group.Meta)
		assert.NotEmpty(t, group.Meta.DocumentType)
		assert.NotEmpty(t, group.Meta.FriendName)

		user1 := mock.User()
		user1.Name = "User1"
		user2 := mock.User()
		user2.Name = "User2"
		err = mock.UserCreate(app.Db, tx, user1)
		assert.Nil(t, err)
		assert.NotNil(t, user1.Id)
		err = mock.UserCreate(app.Db, tx, user2)
		assert.Nil(t, err)
		assert.NotNil(t, user2.Id)

		t.Log("Take the members group list")
		selected, _, err := app.Db.UserDb.UserGroupListUsers(tx, group)
		assert.Nil(t, err)

		t.Log("Selected members should be zero")
		assert.Equal(t, len(selected), 0)

		t.Log("Add two users to the group")
		ug := &domain.UserGroup{User: user1, Group: group}
		err = app.Db.UserDb.UserGroupAdd(tx, ug)
		assert.Nil(t, err)
		selected, _, err = app.Db.UserDb.UserGroupListUsers(tx, group)
		assert.Equal(t, len(selected), 1)

		ug = &domain.UserGroup{User: user2, Group: group}
		err = app.Db.UserDb.UserGroupAdd(tx, ug)
		assert.Nil(t, err)
		selected, _, err = app.Db.UserDb.UserGroupListUsers(tx, group)
		assert.Equal(t, len(selected), 2)
	})
}
