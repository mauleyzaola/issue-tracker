package priority

import (
	"github.com/mauleyzaola/issue-tracker/server/api"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/zenazn/goji/web"
)

type Api struct {
	base *api.ApiBase
}

func (t *Api) RoutesAuth(basePath string, router *web.Mux) {
	router.Get(basePath+"/priority/grid", t.grid)
	router.Get(basePath+"/priority/:id", t.load)
	router.Get(basePath+"/priorities", t.list)
	router.Post(basePath+"/priority", t.save)
	router.Delete(basePath+"/priority/:id", t.remove)
	router.Put(basePath+"/priority", t.save)
}

func (t *Api) init(c web.C) {
	if t.base == nil {
		t.base = c.Env[operations.MIDDLEWARE_BASE_API].(*api.ApiBase)
	}
}
