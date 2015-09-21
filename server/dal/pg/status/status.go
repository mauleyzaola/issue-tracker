package status

import (
	"time"

	"github.com/mauleyzaola/issue-tracker/server/domain"
)

func (t *StatusDb) Load(tx interface{}, id string) (*domain.Status, error) {
	item := &domain.Status{}
	err := t.Base.Executor(tx).SelectOne(item, "select * from status where id=$1", id)
	if err != nil {
		return nil, err
	}
	item.Initialize()

	item.Workflow, err = t.WorkflowLoad(tx, item.Workflow.Id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (t *StatusDb) Create(tx interface{}, item *domain.Status) error {
	err := item.Validate()
	if err != nil {
		return err
	}
	item.DateCreated = time.Now()
	return t.Base.Executor(tx).Insert(item)
}

func (t *StatusDb) Update(tx interface{}, item *domain.Status) error {
	err := item.Validate()
	if err != nil {
		return err
	}
	oldItem, err := t.Load(tx, item.Id)
	if err != nil {
		return err
	}
	item.DateCreated = oldItem.DateCreated
	item.Workflow = oldItem.Workflow
	_, err = t.Base.Executor(tx).Update(item)
	return err
}

func (t *StatusDb) Remove(tx interface{}, id string) (*domain.Status, error) {
	item, err := t.Load(tx, id)
	if err != nil {
		return nil, err
	}
	_, err = t.Base.Executor(tx).Delete(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (t *StatusDb) List(tx interface{}, workflow *domain.Workflow) ([]domain.Status, error) {
	var items []domain.Status
	_, err := t.Base.Executor(tx).Select(&items, "select * from status where idworkflow = $1 order by name", workflow.Id)
	if err != nil {
		return nil, err
	}
	return items, nil
}
