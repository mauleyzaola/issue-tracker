package file_item

import (
	"github.com/mauleyzaola/issue-tracker/server/api"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/zenazn/goji/web"
)

type Api struct {
	base *api.ApiBase
}

func (t *Api) RoutesAuth(basePath string, router *web.Mux) {
	router.Post(basePath+"/file", t.upload)
	router.Get(basePath+"/file/:id", t.download)
	router.Get(basePath+"/files/directory/grid", t.directoryGrid)
	router.Get(basePath+"/files/file/grid", t.fileGrid)
}

func (t *Api) init(c web.C) {
	if t.base == nil {
		t.base = c.Env[operations.MIDDLEWARE_BASE_API].(*api.ApiBase)
	}
}
