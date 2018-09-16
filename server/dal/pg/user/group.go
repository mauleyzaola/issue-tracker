package user

import (
	"database/sql"
	"strings"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

func (t *UserDb) GroupCreate(tx interface{}, item *domain.Group) (err error) {
	err = item.Validate()
	if err != nil {
		return err
	}
	item.DateCreated = time.Now()
	return t.Base.Executor(tx).Insert(item)
}

func (t *UserDb) GroupLoad(tx interface{}, id string) (*domain.Group, error) {
	item := &domain.Group{}
	err := t.Base.Executor(tx).SelectOne(item, "select * from groups where id=$1", id)
	if err != nil {
		return nil, err
	}
	item.Initialize()
	return item, nil
}
func (t *UserDb) GroupRemove(tx interface{}, id string) (item *domain.Group, err error) {
	item, err = t.GroupLoad(tx, id)
	if err != nil {
		return
	}
	query := "delete from user_group where idgroup=$1"
	_, err = t.Base.Executor(tx).Exec(query, id)
	if err != nil {
		return
	}

	query = "delete from permission_scheme_item where idgroup=$1"
	_, err = t.Base.Executor(tx).Exec(query, id)
	if err != nil {
		return
	}

	query = "delete from project_role_member where idgroup=$1"
	_, err = t.Base.Executor(tx).Exec(query, id)
	if err != nil {
		return
	}

	query = "delete from workflow_step_member where idgroup=$1"
	_, err = t.Base.Executor(tx).Exec(query, id)
	if err != nil {
		return
	}

	_, err = t.Base.Executor(tx).Delete(item)
	return
}
func (t *UserDb) GroupUpdate(tx interface{}, item *domain.Group) (err error) {
	oldItem, err := t.GroupLoad(tx, item.Id)
	if err != nil {
		return
	}
	err = item.Validate()
	if err != nil {
		return
	}
	item.DateCreated = oldItem.DateCreated
	item.LastModified = &time.Time{}
	*item.LastModified = time.Now()
	_, err = t.Base.Executor(tx).Update(item)
	return
}
func (t *UserDb) GroupList(tx interface{}) (items []domain.Group, err error) {
	_, err = t.Base.Executor(tx).Select(&items, "select * from groups order by name")
	return
}
func (t *UserDb) GroupGrid(tx interface{}, grid *tecgrid.NgGrid) (err error) {
	query := "select * from groups "
	if len(grid.Query) != 0 {
		query += "where lower(name) like '%" + grid.GetQuery() + "%' "
	}
	grid.MainQuery = query
	fields := strings.Split("id,name,datecreated,lastmodified", ",")
	var rows []domain.Group
	return grid.ExecuteSqlParameters(t.Base.GetTransaction(tx), &rows, fields, nil)
}

func (t *UserDb) UserGroupListGroups(tx interface{}, u *domain.User) (selected []domain.Group, unselected []domain.Group, err error) {
	query := "select g.* " +
		"from groups g " +
		"where exists(select null from user_group where idgroup=g.id and iduser=$1) " +
		"order by g.name"
	_, err = t.Base.Executor(tx).Select(&selected, query, u.Id)
	if err != nil {
		return
	}

	query = "select g.* " +
		"from groups g " +
		"where not exists(select null from user_group where idgroup=g.id and iduser=$1) " +
		"order by g.name"
	_, err = t.Base.Executor(tx).Select(&unselected, query, u.Id)
	if err != nil {
		return
	}

	for i := range selected {
		s := &selected[i]
		s.Initialize()
	}

	for i := range unselected {
		s := &unselected[i]
		s.Initialize()
	}

	return
}

func (t *UserDb) UserGroupListUsers(tx interface{}, g *domain.Group) (selected []domain.User, unselected []domain.User, err error) {
	query := "select u.* " +
		"from users u " +
		"where exists(select null from user_group where iduser=u.id and idgroup=$1) " +
		"order by name, lastname "
	_, err = t.Base.Executor(tx).Select(&selected, query, g.Id)
	if err != nil {
		return
	}

	query = "select u.* " +
		"from users u " +
		"where not exists(select null from user_group where iduser=u.id and idgroup=$1) " +
		"order by name, lastname "
	_, err = t.Base.Executor(tx).Select(&unselected, query, g.Id)

	if err != nil {
		return
	}

	for i := range selected {
		r := &selected[i]
		r.Initialize()
	}
	for i := range unselected {
		r := &unselected[i]
		r.Initialize()
	}
	return
}

func (t *UserDb) userGroupFind(tx interface{}, u *domain.UserGroup) (err error) {
	u.Validate()
	query := "select * from user_group where idgroup=$1 and iduser=$2"
	err = t.Base.Executor(tx).SelectOne(u, query, u.Group.Id, u.User.Id)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func (t *UserDb) UserGroupAdd(tx interface{}, u *domain.UserGroup) error {
	err := t.userGroupFind(tx, u)
	if err != nil {
		return err
	}
	if len(u.Id) != 0 {
		return nil
	}
	return t.Base.Executor(tx).Insert(u)
}

func (t *UserDb) UserGroupRemove(tx interface{}, u *domain.UserGroup) error {
	err := t.userGroupFind(tx, u)
	if err != nil {
		return err
	}
	if len(u.Id) == 0 {
		return nil
	}
	_, err = t.Base.Executor(tx).Delete(u)
	return err
}

func (t *UserDb) UserGroupIsMember(tx interface{}, group *domain.Group, user *domain.User) (ok bool, err error) {
	rowCount, err := t.Base.Executor(tx).SelectInt("select count(*) from user_group where idgroup=$1 and iduser=$2", group.Id, user.Id)
	if err != nil {
		return
	}
	ok = rowCount != 0
	return
}
