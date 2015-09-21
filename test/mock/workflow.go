package mock

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/tecutils"
)

func Workflow() *domain.Workflow {
	item := &domain.Workflow{Name: tecutils.UUID()}
	return item
}

func WorkflowCreate(tx interface{}, op *database.DbOperations, item *domain.Workflow) error {
	return op.StatusDb.WorkflowCreate(tx, item)
}
