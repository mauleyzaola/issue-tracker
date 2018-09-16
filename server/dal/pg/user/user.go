package user

import (
	"errors"
	"strings"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
	"github.com/mauleyzaola/issue-tracker/utils/tecutils"
)

type UserDb struct {
	Base      *pg.Db
	accountDb *database.Account
}

func New(db database.Db) *UserDb {
	base := db.(*pg.Db)
	return &UserDb{Base: base}
}

func (t *UserDb) AccountDb() database.Account {
	return *t.accountDb
}

func (t *UserDb) SetAccountDb(item *database.Account) {
	t.accountDb = item
}

func (t *UserDb) CountSystemAdministrators(tx interface{}) (int64, error) {
	return t.Base.Executor(tx).SelectInt("select count(*) from users where issystemadministrator = $1 and isactive = $2", true, true)
}

func (t *UserDb) Create(tx interface{}, item *domain.User) error {
	err := item.Validate()
	if err != nil {
		return err
	}

	item.Password = tecutils.Encrypt(item.Password)
	item.DateCreated = time.Now()

	err = t.Base.Executor(tx).Insert(item)
	if err != nil {
		return err
	}

	meta := &domain.UserMeta{IdUser: item.Id}
	meta.Initialize()

	err = t.Base.Executor(tx).Insert(meta)
	if err != nil {
		return err
	}

	oldItem, err := t.Load(tx, item.Id)
	if err != nil {
		return err
	}
	*item = *oldItem
	item.Metadata = meta

	return err
}

func (t *UserDb) Load(tx interface{}, id string) (*domain.User, error) {
	item := &domain.User{}
	err := t.Base.Executor(tx).SelectOne(item, "select * from users where id=$1", id)
	if err != nil {
		return nil, err
	}

	meta := &domain.UserMeta{}
	err = t.Base.Executor(tx).SelectOne(meta, "select * from user_meta where iduser=$1", id)
	if err != nil {
		return nil, err
	}
	meta.Initialize()
	item.Metadata = meta

	item.Initialize()

	return item, err
}

func (t *UserDb) Remove(tx interface{}, id string) (*domain.User, error) {
	if t.Base.CurrentSession().User.Id == id {
		return nil, errors.New("cannot remove yourself")
	}

	item, err := t.Load(tx, id)
	if err != nil {
		return nil, err
	}

	_, err = t.Base.Executor(tx).Exec("delete from user_group where iduser=$1", id)
	if err != nil {
		return nil, err
	}

	_, err = t.Base.Executor(tx).Delete(item.Metadata)
	if err != nil {
		return nil, err
	}

	_, err = t.Base.Executor(tx).Delete(item)
	if err != nil {
		return nil, err
	}

	if count, _ := t.CountSystemAdministrators(tx); count == 0 {
		return nil, errors.New("there must be at least one sysadmin user who can login")
	}

	return item, err
}

func (t *UserDb) Update(tx interface{}, item *domain.User) error {
	err := item.Validate()
	if err != nil {
		return err
	}

	oldItem, err := t.Load(tx, item.Id)
	if err != nil {
		return err
	}
	item.DateCreated = oldItem.DateCreated
	item.LastModified = &time.Time{}
	*item.LastModified = time.Now()
	item.LastLogin = oldItem.LastLogin
	item.LoginCount = oldItem.LoginCount
	item.Password = oldItem.Password
	item.TokenEmail = oldItem.TokenEmail
	item.TokenExpires = oldItem.TokenExpires

	_, err = t.Base.Executor(tx).Update(item)
	if err != nil {
		return err
	}

	count, err := t.CountSystemAdministrators(tx)
	if count == 0 && err == nil {
		return errors.New("there must be at least one sysadmin user who can login")
	} else if err != nil {
		return err
	}

	item.Metadata.Initialize()
	item.Metadata.IdUser = oldItem.Metadata.IdUser
	item.Metadata.Id = oldItem.Metadata.Id
	_, err = t.Base.Executor(tx).Update(item.Metadata)
	if err != nil {
		return err
	}

	item.Initialize()
	return nil
}

func (t *UserDb) List(tx interface{}) ([]domain.User, error) {
	var items []domain.User
	_, err := t.Base.Executor(tx).Select(&items, "select * from users order by name, lastname")
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (t *UserDb) Grid(tx interface{}, grid *tecgrid.NgGrid) error {
	query := "select * from users "
	if len(grid.Query) != 0 {
		query += "where lower(name) like '%" + grid.GetQuery() + "%' or lower(lastname) like '%" + grid.GetQuery() + "%' or lower(email) like '%" + grid.GetQuery() + "%'"
	}
	grid.MainQuery = query

	fields := strings.Split("id,name,datecreated,lastmodified,lastname,email,logincount,lastlogin,isactive,issystemadministrator", ",")
	var rows []domain.User
	return grid.ExecuteSqlParameters(t.Base.GetTransaction(tx), &rows, fields, nil)
}
