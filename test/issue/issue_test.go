package issue

import (
	"testing"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestIssueCrud(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		session, err := mock.SessionSetContext(app.Db, tx, true)
		assert.Nil(t, err)
		assert.NotNil(t, session)

		issue := mock.Issue()
		err = mock.IssueCreate(app.Db, tx, issue)
		if !assert.Nil(t, err) {
			t.Log(err)
		}
	})
}
