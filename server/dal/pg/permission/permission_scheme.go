package permission

import (
	"strings"

	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

type PermissionDb struct {
	Base      *pg.Db
	projectDb *database.Project
	userDb    *database.User
}

func New(db database.Db) *PermissionDb {
	base := db.(*pg.Db)
	return &PermissionDb{Base: base}
}

func (t *PermissionDb) ProjectDb() database.Project {
	return *t.projectDb
}

func (t *PermissionDb) SetProjectDb(item *database.Project) {
	t.projectDb = item
}

func (t *PermissionDb) UserDb() database.User {
	return *t.userDb
}

func (t *PermissionDb) SetUserDb(item *database.User) {
	t.userDb = item
}

func (t *PermissionDb) Create(tx interface{}, item *domain.PermissionScheme) error {
	err := item.Validate()
	if err != nil {
		return err
	}
	err = t.Base.Executor(tx).Insert(item)
	return err
}

func (t *PermissionDb) Update(tx interface{}, item *domain.PermissionScheme) error {
	_, err := t.Load(tx, item.Id)
	if err != nil {
		return err
	}
	err = item.Validate()
	if err != nil {
		return err
	}
	_, err = t.Base.Executor(tx).Update(item)
	return err
}

func (t *PermissionDb) Load(tx interface{}, id string) (*domain.PermissionScheme, error) {
	item := &domain.PermissionScheme{}
	err := t.Base.Executor(tx).SelectOne(item, "select * from permission_scheme where id=$1", id)
	if err != nil {
		return nil, err
	}
	item.Initialize()
	return item, nil
}

func (t *PermissionDb) List(tx interface{}) ([]domain.PermissionScheme, error) {
	var items []domain.PermissionScheme
	_, err := t.Base.Executor(tx).Select(&items, "select * from permission_scheme order by name")
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (t *PermissionDb) Remove(tx interface{}, id string) (*domain.PermissionScheme, error) {
	oldItem, err := t.Load(tx, id)
	if err != nil {
		return nil, err
	}
	_, err = t.Base.Executor(tx).Delete(oldItem)
	if err != nil {
		return nil, err
	}
	return oldItem, nil
}

func (t *PermissionDb) Grid(tx interface{}, grid *tecgrid.NgGrid) error {
	query := "select * from permission_scheme "
	if len(grid.Query) != 0 {
		query += "where lower(name) like '%" + grid.GetQuery() + "%' "
	}
	grid.MainQuery = query
	var rows []domain.PermissionScheme
	return grid.ExecuteSqlParameters(t.Base.GetTransaction(tx), &rows, strings.Split("id,name", ","), nil)
}

func (t *PermissionDb) Projects(tx interface{}, item *domain.PermissionScheme) ([]database.ProjectQuery, error) {
	var items []database.ProjectQuery
	_, err := t.Base.Executor(tx).Select(&items, "select * from view_projects where idpermissionscheme=$1 order by name", item.Id)
	if err != nil {
		return nil, err
	}
	for i := range items {
		item := &items[i]
		item.Initialize()
	}
	return items, nil
}

func (t *PermissionDb) ClearAll(tx interface{}, item *domain.PermissionScheme) error {
	_, err := t.Base.Executor(tx).Exec("delete from permission_scheme_item where idpermissionscheme=$1", item.Id)
	return err
}
