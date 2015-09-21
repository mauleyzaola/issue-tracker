package application

import (
	"net/http"

	"github.com/mauleyzaola/issue-tracker/server/api/account"
	"github.com/mauleyzaola/issue-tracker/server/api/file_item"
	"github.com/mauleyzaola/issue-tracker/server/api/issue"
	"github.com/mauleyzaola/issue-tracker/server/api/permission"
	"github.com/mauleyzaola/issue-tracker/server/api/priority"
	"github.com/mauleyzaola/issue-tracker/server/api/project"
	"github.com/mauleyzaola/issue-tracker/server/api/status"
	"github.com/mauleyzaola/issue-tracker/server/api/user"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

func (a *Application) Router() {

	//attach goji's default middleware
	goji.Use(middleware.RealIP)

	//attach gerp middleware
	goji.Use(a.MiddlewareAttach)
	goji.Use(a.MiddlewareJsonResponseHeaders)

	//todo: make this path parametrizable through application config
	baseApiPath := "/api"

	//apply auth middleware
	authPath := baseApiPath + "/auth"
	auth := web.New()
	goji.Handle(authPath+"/*", auth)
	auth.Use(a.MiddlewareAuth)
	goji.Get(authPath, http.RedirectHandler(authPath+"/", http.StatusMovedPermanently))

	//apply sysadmin middleware
	adminPath := baseApiPath + "/admin"
	admin := web.New()
	goji.Handle(adminPath+"/*", admin)
	admin.Use(a.MiddlewareAuth)
	admin.Use(a.MiddlewareSysAdmin)
	goji.Get(adminPath, http.RedirectHandler(adminPath+"/", http.StatusMovedPermanently))

	//apply no auth middleware
	noAuthPath := baseApiPath + "/noauth"
	noAuth := web.New()
	goji.Handle(noAuthPath+"/*", noAuth)
	goji.Get(noAuthPath, http.RedirectHandler(noAuthPath+"/", http.StatusMovedPermanently))

	//assign routes

	accountRt := &account.Api{}
	accountRt.RoutesAuth(authPath, auth)
	accountRt.RoutesNoAuth(noAuthPath, noAuth)
	accountRt.RoutesSysAdminAuth(adminPath, admin)

	fileItemRt := &file_item.Api{}
	fileItemRt.RoutesAuth(authPath, auth)

	issueRt := &issue.Api{}
	issueRt.RoutesAuth(authPath, auth)

	permissionRt := &permission.Api{}
	permissionRt.RoutesAuth(authPath, auth)
	permissionRt.RoutesSysAdminAuth(adminPath, admin)

	priorityRt := &priority.Api{}
	priorityRt.RoutesAuth(authPath, auth)

	projectRt := &project.Api{}
	projectRt.RoutesAuth(authPath, auth)

	statusRt := &status.Api{}
	statusRt.RoutesAuth(authPath, auth)
	statusRt.RoutesSysAdminAuth(adminPath, admin)

	userRt := &user.Api{}
	userRt.RoutesAuth(authPath, auth)
	userRt.RoutesSysAdminAuth(adminPath, admin)
}
