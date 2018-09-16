package mock

import (
	"fmt"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecutils"
)

func User() *domain.User {
	t := &domain.User{}
	t.Email = fmt.Sprintf("%s@test.com", tecutils.UUID())
	t.Name = tecutils.UUID()
	t.LastName = tecutils.UUID()
	t.IsActive = true
	return t
}

func UserCreate(op *database.DbOperations, tx interface{}, item *domain.User) error {
	return op.UserDb.Create(tx, item)
}

func Group() *domain.Group {
	return &domain.Group{Name: tecutils.UUID()}
}

func Role() *domain.Role {
	return &domain.Role{Name: tecutils.UUID()}
}

func RoleCreate(op *database.DbOperations, tx interface{}, item *domain.Role) error {
	return op.UserDb.RoleCreate(tx, item)
}

func GroupCreate(op *database.DbOperations, tx interface{}, item *domain.Group) error {
	return op.UserDb.GroupCreate(tx, item)
}
