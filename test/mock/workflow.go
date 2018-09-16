package mock

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecutils"
)

func Workflow() *domain.Workflow {
	item := &domain.Workflow{Name: tecutils.UUID()}
	return item
}

func WorkflowCreate(op *database.DbOperations, tx interface{}, item *domain.Workflow) error {
	return op.StatusDb.WorkflowCreate(tx, item)
}

func WorkflowStepsCreate(op *database.DbOperations, tx interface{}, workflow *domain.Workflow) error {
	stNew := &domain.Status{Name: "New", Workflow: workflow}
	stInProgress := &domain.Status{Name: "In Progress", Workflow: workflow}
	stResolved := &domain.Status{Name: "Resolved", Workflow: workflow}
	stCancelled := &domain.Status{Name: "Cancelled", Workflow: workflow}

	open := &domain.WorkflowStep{Workflow: workflow, Name: "Open"}
	open.NextStatus = stInProgress

	beginProgress := &domain.WorkflowStep{Workflow: workflow, Name: "Begin Progress"}
	beginProgress.PrevStatus = stNew
	beginProgress.NextStatus = stInProgress

	resolve := &domain.WorkflowStep{Workflow: workflow, Name: "Resolve"}
	resolve.PrevStatus = stInProgress
	resolve.NextStatus = stResolved
	resolve.Resolves = true

	cancels := &domain.WorkflowStep{Workflow: workflow, Name: "Cancel"}
	cancels.PrevStatus = stInProgress
	cancels.NextStatus = stCancelled
	cancels.Cancels = true

	if err := op.StatusDb.Create(tx, stCancelled); err != nil {
		return err
	}

	if err := op.StatusDb.Create(tx, stInProgress); err != nil {
		return err
	}

	if err := op.StatusDb.Create(tx, stNew); err != nil {
		return err
	}

	if err := op.StatusDb.Create(tx, stResolved); err != nil {
		return err
	}

	if err := op.StatusDb.WorkflowStepCreate(tx, open); err != nil {
		return err
	}

	if err := op.StatusDb.WorkflowStepCreate(tx, beginProgress); err != nil {
		return err
	}

	if err := op.StatusDb.WorkflowStepCreate(tx, resolve); err != nil {
		return err
	}

	if err := op.StatusDb.WorkflowStepCreate(tx, cancels); err != nil {
		return err
	}
	return nil
}
