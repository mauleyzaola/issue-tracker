package project

import (
	"github.com/mauleyzaola/issue-tracker/server/api"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/zenazn/goji/web"
)

type Api struct {
	base *api.ApiBase
}

func (t *Api) RoutesAuth(basePath string, router *web.Mux) {

	router.Get(basePath+"/project/createmeta", t.createMeta)
	router.Get(basePath+"/project/grid", t.grid)
	router.Get(basePath+"/project/:id", t.load)
	router.Delete(basePath+"/project/:id", t.remove)
	router.Post(basePath+"/project", t.save)
	router.Put(basePath+"/project", t.save)
	router.Get(basePath+"/project/:id/members", t.roleProjectMembers)

	router.Get(basePath+"/projectrole/:id/members", t.projectRoleMembers)
	router.Post(basePath+"/projectrole/members/add", t.projectRoleMemberAdd)
	router.Post(basePath+"/projectrole/members/remove", t.projectRoleMemberRemove)
}

func (t *Api) init(c web.C) {
	if t.base == nil {
		t.base = c.Env[operations.MIDDLEWARE_BASE_API].(*api.ApiBase)
	}
}
