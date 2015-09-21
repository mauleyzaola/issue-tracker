package permission

import (
	"github.com/mauleyzaola/issue-tracker/server/api"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/zenazn/goji/web"
)

type Api struct {
	base *api.ApiBase
}

func (t *Api) RoutesAuth(basePath string, router *web.Mux) {

	router.Get(basePath+"/permissionscheme/grid", t.grid)
	router.Get(basePath+"/permissionscheme/:id", t.load)
	router.Get(basePath+"/permissionscheme/:id/projects", t.schemeProjects)
	router.Get(basePath+"/permissionschemes", t.list)
	router.Post(basePath+"/permissionscheme/:id/clear", t.clearMembers)

	router.Get(basePath+"/permissionnames", t.names)

	router.Get(basePath+"/permissionschemeitems/:id", t.items)
	router.Post(basePath+"/permissionschemeitem/user/available", t.availablesUser)
}

func (t *Api) RoutesSysAdminAuth(basePath string, router *web.Mux) {
	router.Post(basePath+"/permissionscheme", t.save)
	router.Delete(basePath+"/permissionscheme/:id", t.remove)
	router.Put(basePath+"/permissionscheme", t.save)

	router.Post(basePath+"/permissionschemeitem/add", t.itemAdd)
	router.Post(basePath+"/permissionschemeitem/remove", t.itemRemove)
}

func (t *Api) init(c web.C) {
	if t.base == nil {
		t.base = c.Env[operations.MIDDLEWARE_BASE_API].(*api.ApiBase)
	}
}
