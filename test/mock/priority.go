package mock

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/tecutils"
)

func Priority() *domain.Priority {
	return &domain.Priority{Name: tecutils.UUID()}
}

func PriorityCreate(tx interface{}, op *database.DbOperations, item *domain.Priority) error {
	return op.PriorityDb.Create(tx, item)
}
