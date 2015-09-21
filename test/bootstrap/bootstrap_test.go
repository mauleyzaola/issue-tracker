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
		assert.Nil(t, err)
		assert.NotNil(t, user)
		t.Log(user.Email)

		err = app.Db.BootstrapDb.CreatePermissionNames(tx)
		assert.Nil(t, err)
		names, err := app.Db.PermissionDb.Names(tx)
		assert.Nil(t, err)
		assert.NotNil(t, names)
		assert.Equal(t, true, len(names) != 0)
	})
}
