package user

import (
	"github.com/mauleyzaola/issue-tracker/server/api"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/zenazn/goji/web"
)

type Api struct {
	base *api.ApiBase
}

func (t *Api) RoutesAuth(basePath string, router *web.Mux) {
	router.Get(basePath+"/group/grid", t.groupGrid)
	router.Get(basePath+"/group/:id", t.groupLoad)
	router.Get(basePath+"/groups", t.groupList)
	router.Get(basePath+"/group/:id/users", t.groupUsers)
	router.Get(basePath+"/user/:id/groups", t.userGroups)

	router.Get(basePath+"/role/grid", t.roleGrid)
	router.Get(basePath+"/role/:id", t.roleLoad)
	router.Get(basePath+"/roles", t.roleList)

	router.Get(basePath+"/user/grid", t.grid)
	router.Get(basePath+"/user/:id", t.load)
	router.Get(basePath+"/users", t.list)

}

func (t *Api) RoutesSysAdminAuth(basePath string, router *web.Mux) {
	router.Post(basePath+"/group", t.groupSave)
	router.Delete(basePath+"/group/:id", t.groupRemove)
	router.Put(basePath+"/group", t.groupSave)
	router.Post(basePath+"/group/users/add", t.userAdd)
	router.Post(basePath+"/group/users/remove", t.userRemove)

	router.Post(basePath+"/role", t.roleSave)
	router.Put(basePath+"/role", t.roleSave)
	router.Delete(basePath+"/role/:id", t.roleRemove)

	router.Post(basePath+"/user", t.save)
	router.Delete(basePath+"/user/:id", t.remove)
	router.Put(basePath+"/user", t.save)
}

func (t *Api) init(c web.C) {
	if t.base == nil {
		t.base = c.Env[operations.MIDDLEWARE_BASE_API].(*api.ApiBase)
	}
}
