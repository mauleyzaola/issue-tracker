package priority

import (
	"fmt"
	"strings"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

type PriorityDb struct {
	Base *pg.Db
}

func New(db database.Db) *PriorityDb {
	base := db.(*pg.Db)
	return &PriorityDb{Base: base}
}

func (t *PriorityDb) Create(tx interface{}, item *domain.Priority) error {
	err := item.Validate()
	if err != nil {
		return err
	}
	item.DateCreated = time.Now()
	err = t.ValidateDups(tx, item)
	if err != nil {
		return err
	}

	return t.Base.Executor(tx).Insert(item)
}

func (t *PriorityDb) Load(tx interface{}, id string) (*domain.Priority, error) {
	item := &domain.Priority{}
	err := t.Base.Executor(tx).SelectOne(item, "select * from priority where id=$1", id)
	if err != nil {
		return nil, err
	}
	item.Initialize()
	return item, nil
}

func (t *PriorityDb) Remove(tx interface{}, id string) (*domain.Priority, error) {
	oldItem, err := t.Load(tx, id)
	if err != nil {
		return nil, err
	}
	_, err = t.Base.Executor(tx).Delete(oldItem)
	if err != nil {
		return nil, err
	}
	return oldItem, nil
}

func (t *PriorityDb) Update(tx interface{}, item *domain.Priority) error {
	oldItem, err := t.Load(tx, item.Id)
	if err != nil {
		return err
	}
	item.DateCreated = oldItem.DateCreated
	item.LastModified = &time.Time{}
	*item.LastModified = time.Now()
	err = t.ValidateDups(tx, item)
	if err != nil {
		return err
	}

	_, err = t.Base.Executor(tx).Update(item)
	return err
}

func (t *PriorityDb) List(tx interface{}) ([]domain.Priority, error) {
	var items []domain.Priority
	_, err := t.Base.Executor(tx).Select(&items, "select * from priority order by name")
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (t *PriorityDb) GetFirst(tx interface{}) (*domain.Priority, error) {
	items, err := t.List(tx)
	if err != nil || len(items) == 0 {
		return nil, err
	}
	return &items[0], nil
}

func (t *PriorityDb) Grid(tx interface{}, grid *tecgrid.NgGrid) error {
	query := "select * from priority "
	if len(grid.Query) != 0 {
		query += "where lower(name) like '%" + grid.GetQuery() + "%' "
	}
	grid.MainQuery = query
	fields := strings.Split("id,name,datecreated,lastmodified", ",")
	var rows []domain.Priority
	return grid.ExecuteSqlParameters(t.Base.GetTransaction(tx), &rows, fields, nil)
}

func (t *PriorityDb) ValidateDups(tx interface{}, item *domain.Priority) error {
	var (
		err   error
		count int64
	)
	if len(item.Id) != 0 {
		count, err = t.Base.Executor(tx).SelectInt("select count(*) from priority where lower(name)=$1 and id<>$2", strings.ToLower(item.Name), item.Id)
	} else {
		count, err = t.Base.Executor(tx).SelectInt("select count(*) from priority where lower(name)=$1", strings.ToLower(item.Name))
	}
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}
	return fmt.Errorf("there is already another priority with the name %s", item.Name)
}
