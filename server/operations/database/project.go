package database

import (
	"database/sql"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

type Project interface {
	PermissionDb() Permission
	SetPermissionDb(item *Permission)

	UserDb() User
	SetUserDb(item *User)

	//Creates a project in database
	Create(tx interface{}, item *domain.Project) error

	//Loads a project from database including inner objects
	Load(tx interface{}, id string) (*domain.Project, error)

	//Loads just the basic project information from database
	LoadSimple(tx interface{}, id string) (*domain.Project, error)

	//Removes a project from database
	Remove(tx interface{}, id string) (*domain.Project, error)

	//Updates a project in datebase
	Update(tx interface{}, item *domain.Project) error

	//Processes a query grid of projects applying permissions on the connected user
	Grid(tx interface{}, grid *tecgrid.NgGrid, filter *ProjectFilter) error

	//Generates a simulated sequence next number for a given project
	NextNumber(tx interface{}, id string) (*domain.Project, error)

	//Checks for duplicated projects in database based in the name or pkey
	ValidateDups(tx interface{}, item *domain.Project) error

	//Returns the project data along with other dependencies for the UI
	CreateMeta(tx interface{}, id string) (*ProjectMeta, error)

	//Adds a project role to a project
	RoleAdd(tx interface{}, item *domain.ProjectRole) error

	//Removes a project role from a project
	RoleRemove(tx interface{}, item *domain.ProjectRole) error

	//Loads a project role from a project
	RoleLoad(tx interface{}, project *domain.Project, role *domain.Role) (*domain.ProjectRole, error)

	//Returns all the project roles that a given user for a given project can access
	Roles(tx interface{}, item *domain.Project) ([]domain.ProjectRole, error)

	//Generates all the rows of project roles for a given project
	RoleCreateAll(tx interface{}, item *domain.Project) error

	//Loads a member for a project role, along with its group, role or user related
	RoleMemberLoad(tx interface{}, project *domain.Project, role *domain.Role, user *domain.User, group *domain.Group) (*domain.ProjectRoleMember, error)

	//Creates a new member to the project role
	RoleMemberAdd(tx interface{}, projectRole *domain.ProjectRole, user *domain.User, group *domain.Group) error

	//Removes a member from the project role
	RoleMemberRemove(tx interface{}, item *domain.ProjectRoleMember) error

	//Returns a list of project role members for a given project role
	RoleMembers(tx interface{}, item *domain.ProjectRole) ([]domain.ProjectRoleMember, error)

	//Returns a list of project role members for a given project
	RoleProjectMembers(tx interface{}, item *domain.Project) ([]domain.ProjectRoleMember, error)
}

type ProjectQuery struct {
	domain.Project
	ProjectLead         string  `json:"projectLead"`
	PercentageCompleted float64 `json:"percentageCompleted"`
	PermissionScheme    string  `json:"permissionScheme"`
}

type ProjectFilter struct {
	ProjectLead string
	Resolved    sql.NullBool
}

type ProjectMeta struct {
	Item         *domain.Project      `json:"item"`
	Users        []domain.User        `json:"users"`
	ProjectRoles []domain.ProjectRole `json:"projectRoles"`
}
