package database

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

type Status interface {
	UserDb() User
	SetUserDb(item *User)

	//Loads a status object from database
	Load(tx interface{}, id string) (*domain.Status, error)

	//Creates a new status in database
	Create(tx interface{}, item *domain.Status) error

	//Updates a status in database
	Update(tx interface{}, item *domain.Status) error

	//Removes a status from database
	Remove(tx interface{}, id string) (*domain.Status, error)

	//List all of the statuses the are related with a given Workflow object
	List(tx interface{}, workflow *domain.Workflow) ([]domain.Status, error)

	//Loads a Workflow object from database
	WorkflowLoad(tx interface{}, id string) (item *domain.Workflow, err error)

	//Creates a new Workflow object in database
	WorkflowCreate(tx interface{}, item *domain.Workflow) error

	//Updates a Workflow object in database
	WorkflowUpdate(tx interface{}, item *domain.Workflow) error

	//Removes a Workflow object from database
	WorkflowRemove(tx interface{}, id string) (*domain.Workflow, error)

	//Returns a list of all the Workflows available in database
	WorkflowList(tx interface{}) ([]domain.Workflow, error)

	//Processes a grid query of Workflow objects
	WorkflowGrid(tx interface{}, grid *tecgrid.NgGrid) error

	//Checks for duplicated steps for a given status
	WorkflowStepDups(tx interface{}, step *domain.WorkflowStep) error

	WorkflowStepLoad(tx interface{}, id string) (*domain.WorkflowStep, error)

	//Retrieves a list of all the WorkflowSteps related with a given Workflow object
	WorkflowSteps(tx interface{}, workflow *domain.Workflow) ([]domain.WorkflowStep, error)

	//Creates a new WorkflowStep object in database
	WorkflowStepCreate(tx interface{}, step *domain.WorkflowStep) error

	//Updates a WorkflowStep object in database
	WorkflowStepUpdate(tx interface{}, step *domain.WorkflowStep) error

	//Removes a WorkflowStep object from database
	WorkflowStepRemove(tx interface{}, id string) (*domain.WorkflowStep, error)

	//Returns a list of all the WorkflowSteps that can be applied for a given status
	WorkflowStepAvailableStatus(tx interface{}, workflow *domain.Workflow, prevStatus *domain.Status) ([]domain.WorkflowStep, error)

	//Returns a list of all the WorkflowSteps that can be applied for a given status and a given
	//user, based on his permissions
	WorkflowStepAvailableUser(tx interface{}, workflow *domain.Workflow, prevStatus *domain.Status) ([]domain.WorkflowStep, error)

	//Loads a WorkflowStep member along with its inner objects (groups and users)
	WorkflowStepMemberLoad(tx interface{}, item *domain.WorkflowStepMember) (*domain.WorkflowStepMember, error)

	//Returns a list of all the groups that are related with a WorkflowStep
	WorkflowStepMemberGroups(tx interface{}, item *domain.WorkflowStep) (selected []domain.Group, unselected []domain.Group, err error)

	//Returns a list of all the users that are related with a WorkflowStep
	WorkflowStepMemberUsers(tx interface{}, item *domain.WorkflowStep) (selected []domain.User, unselected []domain.User, err error)

	//Adds a group or user to a WorkflowStep
	WorkflowStepMemberAdd(tx interface{}, item *domain.WorkflowStepMember) error

	//Removes a group or user from a WorkflowStep
	WorkflowStepMemberRemove(tx interface{}, item *domain.WorkflowStepMember) error

	//Returns all the members (users and groups) for a given WorkflowStep
	WorkflowStepMembers(tx interface{}, item *domain.WorkflowStep) ([]domain.WorkflowStepMember, error)

	WorkflowCreateMeta(tx interface{}, item *domain.Workflow) (*WorkflowMeta, error)
}

type WorkflowMeta struct {
	Item     *domain.Workflow      `json:"item"`
	Steps    []domain.WorkflowStep `json:"steps"`
	Statuses []domain.Status       `json:"statuses"`
}
