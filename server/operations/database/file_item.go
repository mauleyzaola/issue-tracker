package database

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

type FileItem interface {

	//Creates a new FileItem in database
	Create(tx interface{}, item *domain.FileItem) error

	//Retrieves a FileItem object along with the data it contains
	Data(tx interface{}, id string) (*domain.FileItem, error)

	//Loads only the metadata of the FileItem without the data
	Load(tx interface{}, id string) (*domain.FileItem, error)

	//Removes the metadata and file data from database
	Remove(tx interface{}, id string) (*domain.FileItem, error)

	//Returs a grid based on information splitted by year and month simulating a directory structure
	DirectoryGrid(tx interface{}, grid *tecgrid.NgGrid) error

	//Returns a grid of files
	FileGrid(tx interface{}, grid *tecgrid.NgGrid, filter *FileFilter) error
}

//Filter used to query the DirectoryGrid
type FileFilter struct {
	YearMonth string
}
