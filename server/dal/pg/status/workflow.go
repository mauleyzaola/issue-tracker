package status

import (
	"strings"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

func (t *StatusDb) WorkflowLoad(tx interface{}, id string) (item *domain.Workflow, err error) {
	item = &domain.Workflow{}
	err = t.Base.Executor(tx).SelectOne(item, "select * from workflow where id=$1", id)
	if err != nil {
		return nil, err
	}
	item.Initialize()
	return item, nil
}

func (t *StatusDb) WorkflowCreate(tx interface{}, item *domain.Workflow) error {
	err := item.Validate()
	if err != nil {
		return err
	}
	item.DateCreated = time.Now()
	return t.Base.Executor(tx).Insert(item)
}

func (t *StatusDb) WorkflowUpdate(tx interface{}, item *domain.Workflow) error {
	err := item.Validate()
	if err != nil {
		return err
	}
	oldItem, err := t.WorkflowLoad(tx, item.Id)
	if err != nil {
		return err
	}
	item.DateCreated = oldItem.DateCreated
	item.LastModified = &time.Time{}
	*item.LastModified = time.Now()
	_, err = t.Base.Executor(tx).Update(item)
	return err
}

func (t *StatusDb) WorkflowRemove(tx interface{}, id string) (*domain.Workflow, error) {
	item, err := t.WorkflowLoad(tx, id)
	if err != nil {
		return nil, err
	}

	_, err = t.Base.Executor(tx).Exec("delete from workflow_step where idworkflow=$1", id)
	if err != nil {
		return nil, err
	}

	_, err = t.Base.Executor(tx).Exec("delete from status where idworkflow=$1", id)
	if err != nil {
		return nil, err
	}

	_, err = t.Base.Executor(tx).Delete(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (t *StatusDb) WorkflowList(tx interface{}) ([]domain.Workflow, error) {
	var items []domain.Workflow
	_, err := t.Base.Executor(tx).Select(&items, "select * from workflow order by name")
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (t *StatusDb) WorkflowGrid(tx interface{}, grid *tecgrid.NgGrid) error {
	query := "select * from workflow "
	if len(grid.Query) != 0 {
		query += "where lower(name) like '%" + grid.GetQuery() + "%'"
	}
	grid.MainQuery = query
	fields := strings.Split("id,name,datecreated,lastmodified", ",")
	var rows []domain.Workflow
	return grid.ExecuteSqlParameters(t.Base.GetTransaction(tx), &rows, fields, nil)
}

func (t *StatusDb) WorkflowCreateMeta(tx interface{}, item *domain.Workflow) (*database.WorkflowMeta, error) {
	var err error
	res := &database.WorkflowMeta{}
	if len(item.Id) != 0 {
		res.Item, err = t.WorkflowLoad(tx, item.Id)
		if err != nil {
			return nil, err
		}
		res.Steps, err = t.WorkflowSteps(tx, res.Item)
		if err != nil {
			return nil, err
		}
		res.Statuses, err = t.List(tx, res.Item)
		if err != nil {
			return nil, err
		}
	}

	if res.Item == nil {
		res.Item = &domain.Workflow{}
	}
	if res.Statuses == nil {
		res.Statuses = make([]domain.Status, 0)
	}
	if res.Steps == nil {
		res.Steps = make([]domain.WorkflowStep, 0)
	}
	return res, nil
}
