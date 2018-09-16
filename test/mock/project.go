package mock

import (
	"fmt"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecutils"
)

func Project(sequence int) *domain.Project {
	item := &domain.Project{}
	item.Name = tecutils.UUID()
	item.Pkey = fmt.Sprintf("XTEST%v", sequence)

	return item
}

func ProjectCreate(op *database.DbOperations, tx interface{}, item *domain.Project) error {
	return op.ProjectDb.Create(tx, item)
}
