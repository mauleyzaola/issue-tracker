package user

import (
	"strings"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

func (t *UserDb) RoleCreate(tx interface{}, item *domain.Role) error {
	err := item.Validate()
	if err != nil {
		return err
	}
	return t.Base.Executor(tx).Insert(item)
}

func (t *UserDb) RoleUpdate(tx interface{}, item *domain.Role) error {
	_, err := t.RoleLoad(tx, item.Id)
	if err != nil {
		return err
	}
	_, err = t.Base.Executor(tx).Update(item)
	return err
}

func (t *UserDb) RoleLoad(tx interface{}, id string) (*domain.Role, error) {
	item := &domain.Role{}
	err := t.Base.Executor(tx).SelectOne(item, "select * from roles where id=$1", id)
	if err != nil {
		return nil, err
	}
	item.Initialize()
	return item, nil
}

func (t *UserDb) RoleList(tx interface{}) ([]domain.Role, error) {
	var items []domain.Role
	_, err := t.Base.Executor(tx).Select(&items, "select * from roles order by name")
	return items, err
}

func (t *UserDb) RoleRemove(tx interface{}, id string) (*domain.Role, error) {
	item, err := t.RoleLoad(tx, id)
	if err != nil {
		return nil, err
	}

	_, err = t.Base.Executor(tx).Exec("delete from permission_scheme_item where idrole=$1", id)
	if err != nil {
		return nil, err
	}

	_, err = t.Base.Executor(tx).Exec("delete from project_role_member where idprojectrole in (select id from project_role where idrole=$1)", id)
	if err != nil {
		return nil, err
	}

	_, err = t.Base.Executor(tx).Exec("delete from project_role where idrole=$1", id)
	if err != nil {
		return nil, err
	}

	_, err = t.Base.Executor(tx).Delete(item)
	return item, err

}

func (t *UserDb) RoleGrid(tx interface{}, grid *tecgrid.NgGrid) error {
	query := "select * from roles "
	if len(grid.Query) != 0 {
		query += "where lower(name) like '%" + grid.GetQuery() + "%' "
	}
	grid.MainQuery = query
	fields := strings.Split("id,name", ",")
	var rows []domain.Role
	return grid.ExecuteSqlParameters(t.Base.GetTransaction(tx), &rows, fields, nil)
}
