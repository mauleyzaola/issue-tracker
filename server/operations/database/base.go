package database

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
)

//Defines operations with connected users
type CurrentSession interface {
	CurrentSession() *domain.Session
	SetCurrentSession(user *domain.Session)
}

//Container of all db operations
type DbOperations struct {
	//Main database implementation
	Db Db

	//Login, authentication operations
	AccountDb Account

	//ERP Bootstrap operations
	BootstrapDb Bootstrap

	//FileItem operations
	FileItemDb FileItem

	//Issue, IssueAttachment, IssueComment and IssueSubscription operations
	IssueDb Issue

	//PermissionName, PermissionScheme, PermissionSchemeItem operations
	PermissionDb Permission

	//Priority operations
	PriorityDb Priority

	//Project, ProjectRole and ProjectRoleMember operations
	ProjectDb Project

	//Session manager
	SessionDb Session

	//Status, Workflow and WorkflowStep operations
	StatusDb Status

	//User, UserGroup, Group and Role operations
	UserDb User
}
