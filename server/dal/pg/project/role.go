package project

import (
	"database/sql"

	"github.com/mauleyzaola/issue-tracker/server/domain"
)

func (t *ProjectDb) RoleAdd(tx interface{}, item *domain.ProjectRole) error {
	item.Initialize()
	newItem, err := t.RoleLoad(tx, item.Project, item.Role)
	if err == sql.ErrNoRows {
		err = t.Base.Executor(tx).Insert(item)
	} else if err == nil {
		item = newItem
	}
	return err
}

func (t *ProjectDb) RoleRemove(tx interface{}, item *domain.ProjectRole) error {
	item.Initialize()
	newItem, err := t.RoleLoad(tx, item.Project, item.Role)
	if err != nil {
		return err
	} else {
		item = newItem
		_, err = t.Base.Executor(tx).Delete(item)
		return err
	}
}

func (t *ProjectDb) RoleLoad(tx interface{}, project *domain.Project, role *domain.Role) (*domain.ProjectRole, error) {
	var err error
	item := &domain.ProjectRole{}
	item.Initialize()
	newItem := &domain.ProjectRole{}
	err = t.Base.Executor(tx).SelectOne(newItem, "select * from project_role where idproject=$1 and idrole=$2", project.Id, role.Id)
	if err != nil {
		return nil, err
	}
	return newItem, nil
}

func (t *ProjectDb) Roles(tx interface{}, item *domain.Project) ([]domain.ProjectRole, error) {
	var items []domain.ProjectRole
	var err error
	_, err = t.Base.Executor(tx).Select(&items, "select * from project_role where idproject = $1", item.Id)

	if err != nil {
		return nil, err
	}

	for i := range items {
		pr := &items[i]
		pr.Initialize()
		r, e := t.UserDb().RoleLoad(tx, pr.Role.Id)
		if e != nil {
			err = e
			return nil, err
		}
		pr.Role = r
		pr.Project = item
	}
	return items, nil
}

func (t *ProjectDb) RoleCreateAll(tx interface{}, item *domain.Project) error {
	roles, err := t.UserDb().RoleList(tx)
	if err != nil {
		return err
	}

	var projectRoles []domain.ProjectRole
	_, err = t.Base.Executor(tx).Select(&projectRoles, "select * from project_role where idproject = $1", item.Id)
	if err != nil {
		return err
	}

	mapProjectRoles := make(map[string]*domain.ProjectRole)
	for i := range projectRoles {
		item := &projectRoles[i]
		item.Initialize()
		mapProjectRoles[item.IdRole] = item
	}

	for i := range roles {
		role := &roles[i]
		if _, ok := mapProjectRoles[role.Id]; !ok {
			newPr := &domain.ProjectRole{Role: role, Project: item}
			err = t.RoleAdd(tx, newPr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
