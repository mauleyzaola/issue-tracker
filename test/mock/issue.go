package mock

import (
	"time"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecutils"
)

func Issue() *domain.Issue {
	item := &domain.Issue{}
	item.Assignee = User()
	item.DueDate = time.Now()
	item.Name = tecutils.UUID()
	item.Priority = Priority()
	item.Reporter = User()
	item.Workflow = Workflow()

	return item
}

func IssueCreate(op *database.DbOperations, tx interface{}, item *domain.Issue) error {
	var err error
	if item.Workflow == nil {
		item.Workflow = Workflow()
	}
	if len(item.Workflow.Id) == 0 {
		err = WorkflowCreate(op, tx, item.Workflow)
		if err != nil {
			return err
		}
	}
	if err = WorkflowStepsCreate(op, tx, item.Workflow); err != nil {
		return err
	}
	return op.IssueDb.Create(tx, item, "")
}
