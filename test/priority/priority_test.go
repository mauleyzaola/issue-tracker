package priority

import (
	"strings"
	"testing"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/application"
	"github.com/mauleyzaola/issue-tracker/test"
	"github.com/mauleyzaola/issue-tracker/test/mock"
	"github.com/mauleyzaola/issue-tracker/utils/tecutils"
	"github.com/stretchr/testify/assert"
)

func TestCrud(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a new priority")
		p := mock.Priority()
		err := app.Db.PriorityDb.Create(tx, p)
		assert.Nil(t, err)
		assert.NotEmpty(t, p.Id)
		assert.Equal(t, p.DateCreated.Year(), time.Now().Year())

		t.Log("Load just created priority")
		p2, err := app.Db.PriorityDb.Load(tx, p.Id)
		assert.Nil(t, err)
		assert.NotNil(t, p2)
		assert.Equal(t, p2.Id, p.Id)
		assert.NotNil(t, p2.Meta)
		assert.NotEmpty(t, p2.Meta.DocumentType)
		assert.NotEmpty(t, p2.Meta.FriendName)

		t.Log("Update that priority and compare with original")
		p2.Name = tecutils.UUID()
		err = app.Db.PriorityDb.Update(tx, p2)
		assert.Nil(t, err)
		assert.NotNil(t, p2.LastModified)
		assert.Equal(t, p2.LastModified.Year(), time.Now().Year())

		p3, err := app.Db.PriorityDb.Load(tx, p.Id)
		assert.Nil(t, err)
		assert.NotNil(t, p3)
		assert.Equal(t, p3.Id, p.Id)
		assert.NotEqual(t, p3.Name, p.Name)
		assert.NotNil(t, p3.LastModified)

		t.Log("Delete the priority and force error on second loading")
		p4, err := app.Db.PriorityDb.Remove(tx, p3.Id)
		assert.Nil(t, err)
		assert.NotNil(t, p4)
		assert.Equal(t, p4.Name, p2.Name)

		p5, err := app.Db.PriorityDb.Load(tx, p2.Id)
		assert.NotNil(t, err)
		assert.Nil(t, p5)
	})
}

func TestList(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Given a list of new priorities")
		total := 10
		oldItems, err := app.Db.PriorityDb.List(tx)
		assert.Nil(t, err)
		oldItemCount := len(oldItems)

		t.Log("Insert a number of priorities")
		for i := 0; i < total; i++ {
			err = mock.PriorityCreate(app.Db, tx, mock.Priority())
			assert.Nil(t, err)
		}
		t.Log("Validate the number of items matches")
		newItems, err := app.Db.PriorityDb.List(tx)
		assert.Nil(t, err)
		assert.Equal(t, len(newItems), oldItemCount+total)

		t.Log("Validate get first priority works as well")
		first, err := app.Db.PriorityDb.GetFirst(tx)
		assert.Nil(t, err)
		assert.NotNil(t, first)
		assert.NotEmpty(t, first.Id)
	})
}

func TestPriorityDups(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Create new priority, create another with same name and force error")
		item := mock.Priority()
		err := mock.PriorityCreate(app.Db, tx, item)
		assert.Nil(t, err)

		item2 := mock.Priority()
		item2.Name = strings.ToUpper(item.Name)
		err = mock.PriorityCreate(app.Db, tx, item2)
		assert.NotNil(t, err)
		if err != nil {
			t.Log(err)
		}
	})
}

func TestPriorityDups2(t *testing.T) {
	test.Runner(func(app *application.Application, tx interface{}) {
		t.Log("Create new priority, update another with same name and force error")
		item := mock.Priority()
		err := mock.PriorityCreate(app.Db, tx, item)
		assert.Nil(t, err)

		item2 := mock.Priority()
		err = mock.PriorityCreate(app.Db, tx, item2)
		assert.Nil(t, err)

		item3, err := app.Db.PriorityDb.Load(tx, item2.Id)
		assert.Nil(t, err)
		item3.Name = strings.ToUpper(item.Name)
		err = app.Db.PriorityDb.Update(tx, item3)
		assert.NotNil(t, err)
		if err != nil {
			t.Log(err)
		}
	})
}
