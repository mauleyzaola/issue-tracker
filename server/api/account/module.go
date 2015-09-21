package account

import (
	"github.com/mauleyzaola/issue-tracker/server/api"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/zenazn/goji/web"
)

type Api struct {
	base *api.ApiBase
}

func (t *Api) RoutesAuth(basePath string, router *web.Mux) {
	router.Post(basePath+"/account/changemypassword", t.changeMyPassword)
	router.Get(basePath+"/account/myprofile", t.myProfileLoad)
	router.Put(basePath+"/account/myprofile", t.myProfileSave)
	router.Post(basePath+"/account/logout", t.logout)
	router.Get(basePath+"/account/token/list", t.sessionList)
	router.Delete(basePath+"/account/session/:id", t.sessionRemove)
	router.Post(basePath+"/account/delaysession", t.sessionDelay)
}

func (t *Api) RoutesNoAuth(basePath string, router *web.Mux) {
	router.Post(basePath+"/account/login", t.login)
	router.Post(basePath+"/account/logout", t.logout)
}

func (t *Api) RoutesSysAdminAuth(basePath string, router *web.Mux) {
	router.Post(basePath+"/account/changepassword", t.changePassword)
}

func (t *Api) init(c web.C) {
	if t.base == nil {
		t.base = c.Env[operations.MIDDLEWARE_BASE_API].(*api.ApiBase)
	}
}
