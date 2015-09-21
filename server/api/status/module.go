package status

import (
	"github.com/mauleyzaola/issue-tracker/server/api"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/zenazn/goji/web"
)

type Api struct {
	base *api.ApiBase
}

func (t *Api) RoutesAuth(basePath string, router *web.Mux) {

	router.Get(basePath+"/status/:id", t.load)
	router.Get(basePath+"/statuses/:workflow", t.list)

	router.Get(basePath+"/workflow/grid", t.workflowGrid)
	router.Get(basePath+"/workflow/:id/createmeta", t.workflowCreateMeta)
	router.Get(basePath+"/workflow/:id", t.workflowLoad)
	router.Get(basePath+"/workflows", t.workflowList)

	router.Get(basePath+"/workflowsteps/user/:workflow", t.workflowStepAvailableUser)
	router.Get(basePath+"/workflowsteps/available/:workflow", t.workflowStepAvailableStatus)
	router.Get(basePath+"/workflowsteps/:workflow", t.workflowStepList)

	router.Get(basePath+"/workflowstepmembers/:id", t.workflowStepMembers)
}

func (t *Api) RoutesSysAdminAuth(basePath string, router *web.Mux) {
	router.Post(basePath+"/status", t.save)
	router.Delete(basePath+"/status/:id", t.remove)
	router.Put(basePath+"/status", t.save)

	router.Post(basePath+"/workflow", t.workflowSave)
	router.Delete(basePath+"/workflow/:id", t.workflowRemove)
	router.Put(basePath+"/workflow", t.workflowSave)

	router.Post(basePath+"/workflowstep", t.workflowStepSave)
	router.Delete(basePath+"/workflowstep/:id", t.workflowStepRemove)
	router.Put(basePath+"/workflowstep", t.workflowStepSave)

	router.Get(basePath+"/workflowstep/:id/groups", t.workflowStepMemberGroups)
	router.Get(basePath+"/workflowstep/:id/users", t.workflowStepMemberUsers)

	router.Post(basePath+"/workflowstepmember/add", t.workflowStepMemberAdd)
	router.Post(basePath+"/workflowstepmember/remove", t.workflowStepMemberRemove)
}

func (t *Api) init(c web.C) {
	if t.base == nil {
		t.base = c.Env[operations.MIDDLEWARE_BASE_API].(*api.ApiBase)
	}
}
