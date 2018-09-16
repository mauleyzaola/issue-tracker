package mock

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecutils"
)

func PermissionScheme() *domain.PermissionScheme {
	item := &domain.PermissionScheme{}
	item.Name = tecutils.UUID()
	return item
}

func PermissionSchemeCreate(op *database.DbOperations, tx interface{}, item *domain.PermissionScheme) error {
	return op.PermissionDb.Create(tx, item)
}

func PermissionSchemeItem() *domain.PermissionSchemeItem {
	item := &domain.PermissionSchemeItem{}
	item.PermissionScheme = PermissionScheme()

	return item
}

func PermissionSchemeItemCreate(op *database.DbOperations, tx interface{}, item *domain.PermissionSchemeItem) error {
	var err error
	if item.PermissionScheme == nil {
		item.PermissionScheme = PermissionScheme()
		err = PermissionSchemeCreate(op, tx, item.PermissionScheme)
		if err != nil {
			return err
		}
	}

	err = op.PermissionDb.ItemAdd(tx, item)
	return err
}
