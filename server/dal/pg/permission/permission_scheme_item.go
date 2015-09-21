package permission

import (
	"database/sql"
	"errors"

	"github.com/mauleyzaola/issue-tracker/server/domain"
)

func (t *PermissionDb) ItemAdd(tx interface{}, item *domain.PermissionSchemeItem) error {
	err := item.Validate()
	if err != nil {
		return err
	}

	oldItem, err := t.ItemLoad(tx, item)
	if err == sql.ErrNoRows {
		err = t.Base.Executor(tx).Insert(item)
	} else if err == nil {
		item = oldItem
	}

	return nil
}

func (t *PermissionDb) ItemLoad(tx interface{}, item *domain.PermissionSchemeItem) (*domain.PermissionSchemeItem, error) {
	var query string
	err := item.Validate()
	if err != nil {
		return nil, err
	}

	oldItem := &domain.PermissionSchemeItem{}
	if item.IdGroup.Valid {
		query = "select * from permission_scheme_item where idpermissionscheme=$1 and idpermissionname=$2 and idgroup=$3"
		err = t.Base.Executor(tx).SelectOne(oldItem, query, item.PermissionScheme.Id, item.PermissionName.Id, item.Group.Id)
	} else if item.IdRole.Valid {
		query = "select * from permission_scheme_item where idpermissionscheme=$1 and idpermissionname=$2 and idrole=$3"
		err = t.Base.Executor(tx).SelectOne(oldItem, query, item.PermissionScheme.Id, item.PermissionName.Id, item.Role.Id)
	} else if item.IdUser.Valid {
		query = "select * from permission_scheme_item where idpermissionscheme=$1 and idpermissionname=$2 and iduser=$3"
		err = t.Base.Executor(tx).SelectOne(oldItem, query, item.PermissionScheme.Id, item.PermissionName.Id, item.User.Id)
	} else {
		return nil, errors.New("No se ha encontrado ningun grupo, rol o usuario en los parametros")
	}
	if err != nil {
		return nil, err
	}
	return oldItem, nil
}

func (t *PermissionDb) ItemRemove(tx interface{}, item *domain.PermissionSchemeItem) error {
	oldItem, err := t.ItemLoad(tx, item)
	if err != nil {
		return err
	}
	_, err = t.Base.Executor(tx).Delete(oldItem)
	return err
}

func (t *PermissionDb) Items(tx interface{}, item *domain.PermissionScheme) ([]domain.PermissionSchemeItem, error) {
	var items []domain.PermissionSchemeItem
	_, err := t.Base.Executor(tx).Select(&items, "select * from permission_scheme_item where idpermissionscheme=$1", item.Id)
	if err != nil {
		return nil, err
	}

	permissionNames, err := t.Names(tx)
	names := make(map[string]*domain.PermissionName)
	for i := range permissionNames {
		name := &permissionNames[i]
		names[name.Id] = name
	}

	if err != nil {
		return nil, err
	}

	for i := range items {
		item := &items[i]
		item.Initialize()

		if name, ok := names[item.PermissionName.Id]; ok {
			item.PermissionName = name
		}

		if item.Group != nil {
			item.Group, err = t.UserDb().GroupLoad(tx, item.Group.Id)
			if err != nil {
				return nil, err
			}
		} else if item.Role != nil {
			item.Role, err = t.UserDb().RoleLoad(tx, item.Role.Id)
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
