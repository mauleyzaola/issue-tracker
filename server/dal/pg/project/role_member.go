package project

import (
	"database/sql"
	"fmt"

	"github.com/mauleyzaola/issue-tracker/server/domain"
)

func (t *ProjectDb) RoleMemberLoad(tx interface{}, project *domain.Project, role *domain.Role, user *domain.User, group *domain.Group) (*domain.ProjectRoleMember, error) {
	if project == nil || len(project.Id) == 0 {
		return nil, fmt.Errorf("missing project parameter")
	}
	if role == nil || len(role.Id) == 0 {
		return nil, fmt.Errorf("missing role parameter")
	}
	pr, err := t.RoleLoad(tx, project, role)
	if err != nil {
		return nil, err
	}
	var pars []interface{}

	query :=
		`
	select	*
	from		project_role_member
	where	idprojectrole = $1 
	`
	pars = append(pars, pr.Id)

	if user != nil && len(user.Id) != 0 {
		query += "and iduser = $2"
		pars = append(pars, user.Id)
	} else if group != nil && len(group.Id) != 0 {
		query += "and idgroup=$2"
		pars = append(pars, group.Id)
	} else {
		return nil, fmt.Errorf("missing user or group parameters")
	}
	item := &domain.ProjectRoleMember{}
	err = t.Base.Executor(tx).SelectOne(item, query, pars...)
	if err != nil {
		return nil, err
	}
	item.ProjectRole = pr
	item.Initialize()
	return item, nil
}

func projectRoleMemberValidatePars(item *domain.ProjectRoleMember) error {
	item.Initialize()

	if item.ProjectRole == nil || len(item.ProjectRole.Id) == 0 {
		return fmt.Errorf("missing project role parameter")
	} else if item.ProjectRole.Project == nil || len(item.ProjectRole.Project.Id) == 0 {
		return fmt.Errorf("missing project parameter")
	}

	if item.Group != nil && len(item.Group.Id) != 0 {
		return nil
	}

	if item.User != nil && len(item.User.Id) != 0 {
		return nil
	}

	return fmt.Errorf("missing group or user parameters")
}

func (t *ProjectDb) RoleMemberAdd(tx interface{}, projectRole *domain.ProjectRole, user *domain.User, group *domain.Group) error {
	item := &domain.ProjectRoleMember{ProjectRole: projectRole, Group: group, User: user}
	err := projectRoleMemberValidatePars(item)
	if err != nil {
		return err
	}

	_, err = t.RoleMemberLoad(tx, item.ProjectRole.Project, item.ProjectRole.Role, item.User, item.Group)
	if err != sql.ErrNoRows {
		return err
	}

	return t.Base.Executor(tx).Insert(item)
}

func (t *ProjectDb) RoleMemberRemove(tx interface{}, item *domain.ProjectRoleMember) error {
	err := projectRoleMemberValidatePars(item)

	oldItem, err := t.RoleMemberLoad(tx, item.ProjectRole.Project, item.ProjectRole.Role, item.User, item.Group)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}

	_, err = t.Base.Executor(tx).Delete(oldItem)
	return err
}

func (t *ProjectDb) RoleMembers(tx interface{}, item *domain.ProjectRole) ([]domain.ProjectRoleMember, error) {
	var items []domain.ProjectRoleMember
	_, err := t.Base.Executor(tx).Select(&items, "select * from project_role_member where idprojectrole=$1", item.Id)

	for i := range items {
		item := &items[i]
		item.Initialize()
		if len(item.IdGroup.String) != 0 {
			item.Group, err = t.UserDb().GroupLoad(tx, item.Group.Id)
			if err != nil {
				return nil, err
			}
		} else if len(item.IdUser.String) != 0 {
			item.User, err = t.UserDb().Load(tx, item.User.Id)
			if err != nil {
				return nil, err
			}
		}
	}
	return items, nil
}

func (t *ProjectDb) RoleProjectMembers(tx interface{}, item *domain.Project) ([]domain.ProjectRoleMember, error) {
	var err error
	var items []domain.ProjectRoleMember
	projectRoles, err := t.Roles(tx, item)
	if err != nil {
		return nil, err
	}

	for i := range projectRoles {
		pr := &projectRoles[i]
		members, err := t.RoleMembers(tx, pr)
		if err != nil {
			return nil, err
		}
		for j := range members {
			member := &members[j]
			member.ProjectRole = pr
			member.Initialize()
			items = append(items, *member)
		}

	}
	if items == nil {
		items = make([]domain.ProjectRoleMember, 0)
	}
	return items, nil
}
