package project

import (
	"testing"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestProjectRole(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("create a new project")
		session, err := mock.SessionSetContext(tx, app.Db, true)
		assert.Nil(t, err)
		assert.NotNil(t, session)

		item := mock.Project(1)
		err = app.Db.ProjectDb.Create(tx, item)
		assert.Nil(t, err)
		assert.NotNil(t, item)
		assert.NotEmpty(t, item.Id)
		assert.Equal(t, item.ProjectLead.Id, session.User.Id)
		assert.Equal(t, item.DateCreated.Year(), time.Now().Year())

		t.Log("create new role, user and group")
		role := mock.Role()
		user := mock.User()
		group := mock.Group()
		err = mock.RoleCreate(tx, app.Db, role)
		assert.Nil(t, err)
		err = mock.UserCreate(tx, app.Db, user)
		assert.Nil(t, err)
		err = mock.GroupCreate(tx, app.Db, group)
		assert.Nil(t, err)

		t.Log("meta should load all objects")
		meta, err := app.Db.ProjectDb.CreateMeta(tx, item.Id)
		assert.Nil(t, err)
		if assert.NotNil(t, meta) {
			assert.NotEmpty(t, meta.Users)
		}
		assert.Equal(t, true, len(meta.ProjectRoles) > 0)
		assert.NotNil(t, meta.Item)
		assert.NotEmpty(t, meta.Item.Id)
		var pr *domain.ProjectRole
		for i := range meta.ProjectRoles {
			item := &meta.ProjectRoles[i]
			item.Initialize()
			if item.Role.Id == role.Id {
				pr = item
				break
			}
		}

		t.Log("assign user to projectRoleMember")
		member := &domain.ProjectRoleMember{}
		member.ProjectRole = pr
		member.ProjectRole.Project = item
		member.User = user
		member.Group = group
		err = member.Validate()
		if assert.NotNil(t, err) {
			t.Log(err)
		}

		member.Group = nil
		member.IdGroup.String = ""
		err = member.Validate()
		assert.Nil(t, err)

		err = app.Db.ProjectDb.RoleMemberAdd(tx, member.ProjectRole, member.User, member.Group)
		assert.Nil(t, err)

		t.Log("assign group to projectRoleMember")
		member = &domain.ProjectRoleMember{}
		member.ProjectRole = pr
		member.ProjectRole.Project = item
		member.Group = group
		err = member.Validate()
		if assert.Nil(t, err) {
			t.Log(err)
		}

		err = app.Db.ProjectDb.RoleMemberAdd(tx, member.ProjectRole, member.User, member.Group)
		assert.Nil(t, err)

		t.Log("load all members for all roles related with project")
		members, err := app.Db.ProjectDb.RoleProjectMembers(tx, item)
		assert.Nil(t, err)
		if assert.NotEmpty(t, members) {
			t.Logf("there are %v members in this project", len(members))
			for i := range members {
				m := &members[i]
				t.Logf("member pr:%s user:%s group:%s", m.IdProjectRole, m.IdUser.String, m.IdGroup.String)
			}
		}

		t.Log("remove user and group from projectRolemMembers")
		err = app.Db.ProjectDb.RoleMemberRemove(tx, member)
		assert.Nil(t, err)

		member = &domain.ProjectRoleMember{}
		member.ProjectRole = pr
		member.User = user
		err = app.Db.ProjectDb.RoleMemberRemove(tx, member)
		assert.Nil(t, err)
	})
}
