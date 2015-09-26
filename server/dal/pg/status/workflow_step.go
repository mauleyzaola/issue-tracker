package status

import (
	"errors"
	"fmt"

	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
)

type StatusDb struct {
	Base   *pg.Db
	userDb *database.User
}

func New(db database.Db) *StatusDb {
	base := db.(*pg.Db)
	return &StatusDb{Base: base}
}

func (t *StatusDb) UserDb() database.User {
	return *t.userDb
}
func (t *StatusDb) SetUserDb(item *database.User) {
	t.userDb = item
}

func (t *StatusDb) WorkflowSteps(tx interface{}, workflow *domain.Workflow) ([]domain.WorkflowStep, error) {
	var statusses []domain.Status
	_, err := t.Base.Executor(tx).Select(&statusses, "select * from status where idworkflow=$1 order by name", workflow.Id)
	if err != nil {
		return nil, err
	}

	var steps []domain.WorkflowStep
	var st map[string]domain.Status
	st = make(map[string]domain.Status)
	for _, item := range statusses {
		st[item.Id] = item
	}

	_, err = t.Base.Executor(tx).Select(&steps, "select * from workflow_step where idworkflow=$1 order by datecreated", workflow.Id)

	for i, _ := range steps {
		item := &steps[i]
		item.NextStatus = &domain.Status{}
		*item.NextStatus = st[item.IdNextStatus]
		if item.IdPrevStatus.Valid {
			item.PrevStatus = &domain.Status{}
			*item.PrevStatus = st[item.IdPrevStatus.String]
		}
	}
	return steps, nil
}

func (t *StatusDb) WorkflowStepDups(tx interface{}, step *domain.WorkflowStep) error {
	step.Initialize()

	query := "select count(*) " +
		"from workflow_step " +
		"where idworkflow = '%s' "

	query = fmt.Sprintf(query, step.IdWorkflow)
	if len(step.Id) != 0 {
		query += fmt.Sprintf("and id != '%s' ", step.Id)
	}
	if len(step.IdNextStatus) != 0 {
		query += fmt.Sprintf("and idnextstatus = '%s' ", step.NextStatus.Id)
	}
	if len(step.IdPrevStatus.String) != 0 {
		query += fmt.Sprintf("and idprevstatus = '%s' ", step.PrevStatus.Id)
	}

	rowCount, err := t.Base.Executor(tx).SelectInt(query)
	if err != nil {
		return err
	}

	if rowCount != 0 {
		return errors.New("there is already a step for the same workflow with the same values")
	}

	if len(step.IdPrevStatus.String) == 0 {
		query = "select count(*)  " +
			"from workflow_step " +
			"where idworkflow = $1 " +
			"and idprevstatus is null "

		if len(step.Id) != 0 {
			query += "and id <> $2 "
			rowCount, err = t.Base.Executor(tx).SelectInt(query, step.IdWorkflow, step.Id)
		} else {
			rowCount, err = t.Base.Executor(tx).SelectInt(query, step.IdWorkflow)
		}

		if err != nil {
			return err
		}
		if rowCount != 0 {
			return errors.New("there can only be one step without prev status for a workflow")
		}
	}

	return nil
}

func (t *StatusDb) WorkflowStepCreate(tx interface{}, item *domain.WorkflowStep) error {
	err := item.Validate()
	if err != nil {
		return err
	}

	err = t.WorkflowStepDups(tx, item)
	if err != nil {
		return err
	}

	if len(item.IdPrevStatus.String) != 0 {
		err = t.Base.Executor(tx).SelectOne(item.PrevStatus, "select * from status where id=$1", item.PrevStatus.Id)
		if err != nil {
			return err
		}
	}

	if len(item.IdNextStatus) != 0 {
		err = t.Base.Executor(tx).SelectOne(item.NextStatus, "select * from status where id=$1", item.NextStatus.Id)
		if err != nil {
			return err
		}
	}

	return t.Base.Executor(tx).Insert(item)
}

func (t *StatusDb) WorkflowStepUpdate(tx interface{}, item *domain.WorkflowStep) error {
	err := item.Validate()
	if err != nil {
		return err
	}

	err = t.WorkflowStepDups(tx, item)
	if err != nil {
		return err
	}

	oldItem := &domain.WorkflowStep{}
	err = t.Base.Executor(tx).SelectOne(oldItem, "select * from workflow_step where id=$1", item.Id)
	if err != nil {
		return err
	}

	item.IdWorkflow = oldItem.IdWorkflow
	item.DateCreated = oldItem.DateCreated
	_, err = t.Base.Executor(tx).Update(item)
	return err

}

func (t *StatusDb) WorkflowStepLoad(tx interface{}, id string) (*domain.WorkflowStep, error) {
	item := &domain.WorkflowStep{}
	err := t.Base.Executor(tx).SelectOne(item, "select * from workflow_step where id=$1", id)
	if err != nil {
		return nil, err
	}
	item.Initialize()
	item.Workflow, err = t.WorkflowLoad(tx, item.Workflow.Id)
	if err != nil {
		return nil, err
	}
	if item.PrevStatus != nil {
		item.PrevStatus, err = t.Load(tx, item.PrevStatus.Id)
		if err != nil {
			return nil, err
		}
	}
	if item.NextStatus != nil {
		item.NextStatus, err = t.Load(tx, item.NextStatus.Id)
		if err != nil {
			return nil, err
		}
	}
	return item, nil
}

func (t *StatusDb) WorkflowStepRemove(tx interface{}, id string) (*domain.WorkflowStep, error) {
	return nil, nil
}

func (t *StatusDb) WorkflowStepAvailableStatus(tx interface{}, workflow *domain.Workflow, prevStatus *domain.Status) ([]domain.WorkflowStep, error) {
	query := "select * from workflow_step where idworkflow =$1 "
	var steps []domain.WorkflowStep
	var err error
	items := make([]domain.WorkflowStep, 0)

	if prevStatus == nil || len(prevStatus.Id) == 0 {
		query += "and idprevstatus is null "
		_, err = t.Base.Executor(tx).Select(&steps, query, workflow.Id)
		if err != nil {
			return nil, err
		}
	} else {
		query += "and idprevstatus = $2 "
		_, err = t.Base.Executor(tx).Select(&steps, query, workflow.Id, prevStatus.Id)
		if err != nil {
			return nil, err
		}
	}

	for i, _ := range steps {
		item := &steps[i]
		status, err := t.Load(tx, item.IdNextStatus)
		if err != nil {
			return nil, err
		}
		item.NextStatus = status
		items = append(items, *item)
	}

	return steps, nil
}

func (t *StatusDb) WorkflowStepAvailableUser(tx interface{}, workflow *domain.Workflow, prevStatus *domain.Status) ([]domain.WorkflowStep, error) {
	var items []domain.WorkflowStep
	steps, err := t.WorkflowStepAvailableStatus(tx, workflow, prevStatus)
	if err != nil {
		return nil, err
	}

	if t.Base.CurrentSession().User.IsSystemAdministrator {
		return steps, nil
	}

	for i := range steps {
		step := &steps[i]
		members, e := t.WorkflowStepMembers(tx, step)
		if e != nil {
			return nil, e
		}
		found := false
		for j := range members {
			member := &members[j]
			if member.User != nil && member.User.Id == t.Base.CurrentSession().User.Id {
				found = true
				break
			} else if member.Group != nil {
				ok, e := t.UserDb().UserGroupIsMember(tx, member.Group, t.Base.CurrentSession().User)
				if e != nil {
					return nil, err
				}
				if ok {
					found = true
					break
				}
			}
		}
		if found {
			items = append(items, *step)
		}
	}

	if items == nil {
		items = make([]domain.WorkflowStep, 0)
	}
	return items, nil
}

func (t *StatusDb) WorkflowStepMemberLoad(tx interface{}, item *domain.WorkflowStepMember) (*domain.WorkflowStepMember, error) {
	oldItem := &domain.WorkflowStepMember{}
	err := item.Validate()
	if err != nil {
		return nil, err
	}
	if len(item.IdGroup.String) != 0 {
		err = t.Base.Executor(tx).SelectOne(oldItem, "select * from workflow_step_member where idworkflowstep=$1 and idgroup=$2", item.WorkflowStep.Id, item.Group.Id)
		if err != nil {
			return nil, err
		}
	} else if len(item.IdUser.String) != 0 {
		err = t.Base.Executor(tx).SelectOne(oldItem, "select * from workflow_step_member where idworkflowstep=$1 and iduser=$2", item.WorkflowStep.Id, item.User.Id)
		if err != nil {
			return nil, err
		}
	}
	return oldItem, nil
}

func (t *StatusDb) WorkflowStepMemberGroups(tx interface{}, item *domain.WorkflowStep) (selected []domain.Group, unselected []domain.Group, err error) {
	query := `
		select	u.*
		from		groups u
		where	not exists(	select	null
					from	workflow_step_member
					where	idworkflowstep = $1
					and	idgroup = u.id)
		order by u.name
	`

	_, err = t.Base.Executor(tx).Select(&unselected, query, item.Id)
	if err != nil {
		return
	}

	query =
		`
		select	u.*
		from	workflow_step_member m
		join		groups u on u.id = m.idgroup
		where	m.idworkflowstep = $1
		order by u.name
	`

	_, err = t.Base.Executor(tx).Select(&selected, query, item.Id)
	if err != nil {
		return
	}

	for i := range unselected {
		u := &unselected[i]
		u.Initialize()
	}
	for i := range selected {
		u := &selected[i]
		u.Initialize()
	}

	return
}

func (t *StatusDb) WorkflowStepMemberUsers(tx interface{}, item *domain.WorkflowStep) (selected []domain.User, unselected []domain.User, err error) {
	query := `
		select	u.*
		from	users u
		where	not exists(	select	null
					from	workflow_step_member
					where	idworkflowstep = $1
					and	iduser = u.id)
		order by u.name, u.lastname
	`

	_, err = t.Base.Executor(tx).Select(&unselected, query, item.Id)
	if err != nil {
		return
	}

	query =
		`
		select	u.*
		from	workflow_step_member m
		join	users u on u.id = m.iduser
		where	m.idworkflowstep = $1
		order by u.name, u.lastname
	`

	_, err = t.Base.Executor(tx).Select(&selected, query, item.Id)
	if err != nil {
		return
	}

	for i := range unselected {
		u := &unselected[i]
		u.Initialize()
	}
	for i := range selected {
		u := &selected[i]
		u.Initialize()
	}

	return
}

func (t *StatusDb) WorkflowStepMemberAdd(tx interface{}, item *domain.WorkflowStepMember) error {
	err := t.WorkflowStepMemberRemove(tx, item)
	if err != nil {
		return err
	}
	err = t.Base.Executor(tx).Insert(item)
	return err
}

func (t *StatusDb) WorkflowStepMemberRemove(tx interface{}, item *domain.WorkflowStepMember) error {
	err := item.Validate()
	if err != nil {
		return err
	}
	if len(item.IdGroup.String) != 0 {
		_, err = t.Base.Executor(tx).Exec("delete from workflow_step_member where idworkflowstep=$1 and idgroup=$2", item.WorkflowStep.Id, item.Group.Id)

	} else if len(item.IdUser.String) != 0 {
		_, err = t.Base.Executor(tx).Exec("delete from workflow_step_member where idworkflowstep=$1 and iduser=$2", item.WorkflowStep.Id, item.User.Id)
	}
	return err
}

func (t *StatusDb) WorkflowStepMembers(tx interface{}, item *domain.WorkflowStep) ([]domain.WorkflowStepMember, error) {
	var items []domain.WorkflowStepMember
	_, err := t.Base.Executor(tx).Select(&items, "select * from workflow_step_member where idworkflowstep=$1", item.Id)
	if err != nil {
		return nil, err
	}
	for i := range items {
		item := &items[i]
		item.Initialize()
		if item.Group != nil {
			item.Group, err = t.UserDb().GroupLoad(tx, item.Group.Id)
			if err != nil {
				return nil, err
			}
		} else if item.User != nil {
			item.User, err = t.UserDb().Load(tx, item.User.Id)
			if err != nil {
				return nil, err
			}
		}
	}
	return items, nil
}
