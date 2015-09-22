package issue

import (
	"github.com/mauleyzaola/issue-tracker/server/api"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/zenazn/goji/web"
)

type Api struct {
	base *api.ApiBase
}

func (t *Api) RoutesAuth(authPath string, router *web.Mux) {
	router.Get(authPath+"/issue/groupall", t.dashboardGroupAll)
	router.Get(authPath+"/issue/group/bydatatype", t.dashboardGroupDataType)
	router.Get(authPath+"/issue/createmeta", t.createMeta)
	router.Get(authPath+"/issue/createmeta/:pkey", t.createMeta)
	router.Get(authPath+"/issue/browse/:pkey", t.load)
	router.Get(authPath+"/issue/mysubscriptions", t.mySuscriptions)
	router.Get(authPath+"/issue/grid", t.grid)
	router.Get(authPath+"/issue/:id", t.load)
	router.Get(authPath+"/issue", t.load)
	router.Get(authPath+"/issue/:id/children", t.children)

	router.Get(authPath+"/issue/attachment/:id", t.attachmentLoad)
	router.Get(authPath+"/issue/:id/attachments", t.attachmentList)
	router.Post(authPath+"/issue/:id/attachment", t.attachmentAdd)
	router.Delete(authPath+"/issue/:id/attachment", t.attachmentRemove)

	router.Post(authPath+"/issue/status", t.status)

	router.Get(authPath+"/issue/:id/comments", t.comments)
	router.Post(authPath+"/issue/:id/comment", t.commentSave)
	router.Put(authPath+"/issue/:id/comment", t.commentSave)
	router.Delete(authPath+"/issue/:id/comment", t.commentRemove)

	router.Post(authPath+"/issue/:id/subscription/toggle", t.subscriptionToogle)
	router.Get(authPath+"/issue/:id/subscribed", t.subscriptionLoad)
	router.Get(authPath+"/issue/:id/subscribedselected", t.subscriptions)
	router.Post(authPath+"/issue/subscription/toggle/any", t.subscriptionToogleUser)
	router.Post(authPath+"/issue", t.save)
	router.Post(authPath+"/issue/assigntome", t.assignToMe)
	router.Post(authPath+"/issue/reporterisme", t.reporterIsMe)
	router.Delete(authPath+"/issue/:id", t.remove)
	router.Put(authPath+"/issue", t.save)
	router.Post(authPath+"/issue/:id/move", t.move)
}

func (t *Api) init(c web.C) {
	if t.base == nil {
		t.base = c.Env[operations.MIDDLEWARE_BASE_API].(*api.ApiBase)
	}
}
