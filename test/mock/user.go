package mock

import (
	"fmt"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/tecutils"
)

func User() *domain.User {
	t := &domain.User{}
	t.Email = fmt.Sprintf("%s@test.com", tecutils.UUID())
	t.Name = tecutils.UUID()
	t.LastName = tecutils.UUID()
	t.IsActive = true
	return t
}

func UserCreate(tx interface{}, op *database.DbOperations, item *domain.User) error {
	return op.UserDb.Create(tx, item)
}

func Group() *domain.Group {
	return &domain.Group{Name: tecutils.UUID()}
}

func Role() *domain.Role {
	return &domain.Role{Name: tecutils.UUID()}
}

func RoleCreate(tx interface{}, op *database.DbOperations, item *domain.Role) error {
	return op.UserDb.RoleCreate(tx, item)
}

func GroupCreate(tx interface{}, op *database.DbOperations, item *domain.Group) error {
	return op.UserDb.GroupCreate(tx, item)
}
