package bootstrap

import (
	"testing"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/stretchr/testify/assert"
)

func TestBootstrapSequence(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		_, user, err := app.Db.BootstrapDb.BootstrapAdminUser(tx)
		if !assert.Nil(t, err) {
			t.Log(err)
			return
		}

		assert.NotNil(t, user)
		t.Log(user.Email)

		err = app.Db.BootstrapDb.CreatePermissionNames(tx)
		if !assert.Nil(t, err) {
			t.Log(err)
			return
		}
		names, err := app.Db.PermissionDb.Names(tx)
		if !assert.Nil(t, err) {
			t.Log(err)
			return
		}
		assert.NotNil(t, names)
		assert.Equal(t, true, len(names) != 0)
	})
}
