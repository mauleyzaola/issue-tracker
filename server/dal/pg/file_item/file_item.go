package fileitem

import (
	"fmt"
	"strings"
	"time"

	"github.com/mauleyzaola/gorp"
	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

const (
	FILE_ITEM_FIELDS = "id,iduser,bytes,extension,name,datecreated,mimetype"
)

type FileItemDb struct {
	base *pg.Db
}

func New(db database.Db) *FileItemDb {
	base := db.(*pg.Db)
	return &FileItemDb{base: base}
}

func (t *FileItemDb) Create(tx interface{}, item *domain.FileItem) error {
	item.User = t.base.CurrentSession().User
	item.DateCreated = time.Now()
	err := item.Validate()
	if err != nil {
		return err
	}

	err = t.base.Executor(tx).Insert(item)
	item.FileData = nil
	return err
}

func (t *FileItemDb) Data(tx interface{}, id string) (*domain.FileItem, error) {
	item := &domain.FileItem{}
	err := t.base.Executor(tx).SelectOne(item, "select * from file_item where id=$1", id)
	return item, err
}

func (t *FileItemDb) Load(tx interface{}, id string) (*domain.FileItem, error) {
	item := &domain.FileItem{}
	query := fmt.Sprintf("select %s from file_item where id=$1", FILE_ITEM_FIELDS)
	err := t.base.Executor(tx).SelectOne(item, query, id)
	if err != nil {
		return nil, err
	}
	item.Initialize()
	return item, nil
}

func (t *FileItemDb) Remove(tx interface{}, id string) (*domain.FileItem, error) {
	item, err := t.Load(tx, id)
	if err != nil {
		return nil, err
	}
	_, err = t.base.Executor(tx).Delete(item)
	return item, err
}

func (t *FileItemDb) DirectoryGrid(tx interface{}, grid *tecgrid.NgGrid) error {
	type Directory struct {
		DateCreated string  `json:"dateCreated"`
		Bytes       float64 `json:"bytes"`
		ItemCount   int64   `json:"itemCount"`
	}
	query := "select	datecreated, sum(bytes) as bytes, count(*) as itemcount " +
		"from " +
		"(select	bytes, substr(cast(date_trunc('month', datecreated) as text),1,7) as datecreated " +
		"from 	file_item) as t " +
		"where 1 = 1 "
	if len(grid.GetQuery()) != 0 {
		query += "and lower(datecreated) like '%" + grid.GetQuery() + "%' "
	}
	query += "group by t.datecreated "
	fields := strings.Split("datecreated,bytes,itemcount", ",")

	var rows []Directory
	grid.MainQuery = query
	tran := tx.(gorp.SqlExecutor)
	return grid.ExecuteSqlParameters(tran, &rows, fields, nil)
}

func (t *FileItemDb) FileGrid(tx interface{}, grid *tecgrid.NgGrid, filter *database.FileFilter) error {
	var pars []interface{}
	query := "select " + FILE_ITEM_FIELDS + " " +
		"from " +
		"(select " + FILE_ITEM_FIELDS + ", " +
		"substr(cast(date_trunc('month', datecreated) as text),1,7) as yearmonth " +
		"from file_item " +
		"where 1 = 1 "

	if len(grid.GetQuery()) != 0 {
		query += "and lower(name) like '%" + grid.GetQuery() + "%' "
	}

	query += ") as t " +
		"where 1 = 1 "

	if filter != nil {
		if len(filter.YearMonth) != 0 {
			pars = append(pars, filter.YearMonth)
			query += fmt.Sprintf("and yearmonth=$%v", len(pars))
		}
	}

	var rows []domain.FileItem
	grid.MainQuery = query
	return grid.ExecuteSqlParameters(t.base.GetTransaction(tx), &rows, strings.Split(FILE_ITEM_FIELDS, ","), pars)

}
