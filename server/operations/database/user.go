package database

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

type User interface {
	AccountDb() Account
	SetAccountDb(item *Account)

	//Counts the number of sys admin users that are active in database
	CountSystemAdministrators(tx interface{}) (int64, error)

	//Creates a new user in database
	Create(tx interface{}, item *domain.User) error

	//Loads an user and inner objects
	Load(tx interface{}, id string) (*domain.User, error)

	//Removes an user from database
	Remove(tx interface{}, id string) (*domain.User, error)

	//Updates an user in database
	Update(tx interface{}, item *domain.User) error

	//Returns a list of all the users in database
	List(tx interface{}) ([]domain.User, error)

	//Processes a grid query of User objects
	Grid(tx interface{}, grid *tecgrid.NgGrid) error

	//Creates a new group of users
	GroupCreate(tx interface{}, item *domain.Group) (err error)

	//Loads a group from database
	GroupLoad(tx interface{}, id string) (item *domain.Group, err error)

	//Removes a group from database
	GroupRemove(tx interface{}, id string) (item *domain.Group, err error)

	//Updates a group in database
	GroupUpdate(tx interface{}, item *domain.Group) (err error)

	//Returns a list of all the groups in database
	GroupList(tx interface{}) (items []domain.Group, err error)

	//Processes a grid query of groups in database
	GroupGrid(tx interface{}, grid *tecgrid.NgGrid) (err error)

	//Return the selected and unselected groups for a given User object
	UserGroupListGroups(tx interface{}, u *domain.User) (selected []domain.Group, unselected []domain.Group, err error)

	//Returns the selected and unselected users for a given Group object
	UserGroupListUsers(tx interface{}, g *domain.Group) (selected []domain.User, unselected []domain.User, err error)

	//Adds an User to a Group or vice versa
	UserGroupAdd(tx interface{}, u *domain.UserGroup) error

	//Removes an User to a Group
	UserGroupRemove(tx interface{}, u *domain.UserGroup) error

	//Returns true if the user is a member of the group
	UserGroupIsMember(tx interface{}, group *domain.Group, user *domain.User) (ok bool, err error)

	//Creates a new role in database
	RoleCreate(tx interface{}, item *domain.Role) error

	//Updates a role in database
	RoleUpdate(tx interface{}, item *domain.Role) error

	//Loads a role from database
	RoleLoad(tx interface{}, id string) (*domain.Role, error)

	//Returns a list of roles from database
	RoleList(tx interface{}) ([]domain.Role, error)

	//Removes a role from database
	RoleRemove(tx interface{}, id string) (*domain.Role, error)

	//Processes a grid query of roles
	RoleGrid(tx interface{}, grid *tecgrid.NgGrid) error
}
