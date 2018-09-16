package file_item

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
	"github.com/mauleyzaola/issue-tracker/utils/tecutils"
	"github.com/zenazn/goji/web"
)

func (t *Api) upload(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	item := &domain.FileItem{}

	// the FormFile function takes in the POST input id/name = file
	file, header, err := r.FormFile("file")
	defer file.Close()

	if err != nil {
		t.base.ErrorResponse(nil, err, w)
		return
	}

	item.MimeType = header.Header.Get(operations.CONTENT_TYPE_HEADER)
	item.Name = header.Filename
	item.Id = tecutils.UUID()

	tmpFilePath := path.Join(os.TempDir(), item.Id)
	out, err := os.Create(tmpFilePath)
	if err != nil {
		t.base.ErrorResponse(nil, err, w)
		return
	}

	//copy the file content to hd
	_, err = io.Copy(out, file)
	if err != nil {
		t.base.ErrorResponse(nil, err, w)
		return
	}

	//open the file to know its real size on disk
	f, err := os.Open(tmpFilePath)
	stat, _ := f.Stat()
	item.Bytes = int64(stat.Size())
	item.FileData, err = ioutil.ReadAll(f)
	if err != nil {
		t.base.ErrorResponse(nil, err, w)
		return
	}
	item.Extension = path.Ext(item.Name)

	//remove temporary file
	out.Close()
	defer os.Remove(tmpFilePath)

	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

	err = t.base.Database.FileItemDb.Create(tx, item)
	t.base.Response(tx, err, w, item)
}

func (t *Api) download(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	id := t.base.ParamValue("id", c, r)
	item, err := t.base.Database.FileItemDb.Data(tx, id)
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}

	w.Header().Set(operations.CONTENT_TYPE_HEADER, item.MimeType)
	w.Header().Set(operations.CONTENT_DISPOSITION, fmt.Sprintf("inline; filename='%s'", item.Name))
	w.Write(item.FileData)
}

func (t *Api) directoryGrid(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	grid := tecgrid.ParseQueryString(r.URL.Query())
	err = t.base.Database.FileItemDb.DirectoryGrid(tx, grid)
	t.base.Response(tx, err, w, grid)
}

func (t *Api) fileGrid(c web.C, w http.ResponseWriter, r *http.Request) {
	t.init(c)
	tx, err := t.base.Database.Db.Begin()
	if err != nil {
		t.base.ErrorResponse(tx, err, w)
		return
	}
	grid := tecgrid.ParseQueryString(r.URL.Query())
	filter := &database.FileFilter{YearMonth: t.base.ParamValue("yearMonth", c, r)}
	err = t.base.Database.FileItemDb.FileGrid(tx, grid, filter)
	t.base.Response(tx, err, w, grid)
}
