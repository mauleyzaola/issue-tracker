package database

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

type Priority interface {

	//Creates a new priority in database
	Create(tx interface{}, priority *domain.Priority) error

	//Loads a priority from database
	Load(tx interface{}, id string) (*domain.Priority, error)

	//Removes a priority from database
	Remove(tx interface{}, id string) (*domain.Priority, error)

	//Updates a priority in database
	Update(tx interface{}, priority *domain.Priority) error

	//Returns a list of all the priorities in database
	List(tx interface{}) ([]domain.Priority, error)

	//Returns the first priority available in database
	GetFirst(tx interface{}) (*domain.Priority, error)

	//Processes a query grid of priorities
	Grid(tx interface{}, grid *tecgrid.NgGrid) error

	ValidateDups(tx interface{}, item *domain.Priority) error
}
