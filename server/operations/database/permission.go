package database

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

type Permission interface {
	ProjectDb() Project
	SetProjectDb(item *Project)

	UserDb() User
	SetUserDb(item *User)

	//Returns a list of all the permission names in database
	Names(tx interface{}) ([]domain.PermissionName, error)

	//Given an user and a permission name, optionaly also an Issue (evaluating the project where it belongs to)
	//Evaluates if the access is granted or not
	AllowedUser(tx interface{}, user *domain.User, issue *domain.Issue, permission *domain.PermissionName) (ok bool, err error)

	//Given an user and optionaly an Issue, returns the list of permission names the user can access
	AvailablesUser(tx interface{}, user *domain.User, issue *domain.Issue) ([]domain.PermissionName, error)

	//Creates a new permission scheme item
	ItemAdd(tx interface{}, item *domain.PermissionSchemeItem) error

	//Loads a permission scheme item
	ItemLoad(tx interface{}, item *domain.PermissionSchemeItem) (*domain.PermissionSchemeItem, error)

	//Removes a permission scheme item
	ItemRemove(tx interface{}, item *domain.PermissionSchemeItem) error

	//Returns a list of all permission scheme items that belong to a given permission scheme and user
	Items(tx interface{}, item *domain.PermissionScheme) ([]domain.PermissionSchemeItem, error)

	//Creates a permission scheme on database
	Create(tx interface{}, item *domain.PermissionScheme) error

	//Updates a permission scheme in database
	Update(tx interface{}, item *domain.PermissionScheme) error

	//Loads the permission scheme from database
	Load(tx interface{}, id string) (*domain.PermissionScheme, error)

	//Returns a list of all the permission schemes available
	List(tx interface{}) ([]domain.PermissionScheme, error)

	//Removes a permission scheme from database
	Remove(tx interface{}, id string) (*domain.PermissionScheme, error)

	//Processes a grid query for permission schemes
	Grid(tx interface{}, grid *tecgrid.NgGrid) error

	//Returns a list of all the projects which are using a given permission scheme
	Projects(tx interface{}, item *domain.PermissionScheme) ([]ProjectQuery, error)

	//Clears all the permission scheme items for a given permission scheme
	ClearAll(tx interface{}, item *domain.PermissionScheme) error
}
