package project

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

type ProjectDb struct {
	Base         *pg.Db
	permissionDb *database.Permission
	userDb       *database.User
}

func New(db database.Db) *ProjectDb {
	base := db.(*pg.Db)
	return &ProjectDb{Base: base}
}

func (t *ProjectDb) PermissionDb() database.Permission {
	return *t.permissionDb
}
func (t *ProjectDb) SetPermissionDb(item *database.Permission) {
	t.permissionDb = item
}

func (t *ProjectDb) UserDb() database.User {
	return *t.userDb
}

func (t *ProjectDb) SetUserDb(item *database.User) {
	t.userDb = item
}

func (t *ProjectDb) Create(tx interface{}, item *domain.Project) error {
	if !t.Base.CurrentSession().User.IsSystemAdministrator {
		return errors.New("db.Base.AccessDenied()")
	}

	err := item.Validate()
	if err != nil {
		return err
	}

	if len(item.IdProjectLead) == 0 {
		item.ProjectLead = t.Base.CurrentSession().User
	}
	item.Initialize()
	item.Next = 0
	item.DateCreated = time.Now()

	err = t.ValidateDups(tx, item)
	if err != nil {
		return err
	}

	err = t.Base.Executor(tx).Insert(item)
	if err != nil {
		return err
	}

	err = t.RoleCreateAll(tx, item)
	if err != nil {
		return err
	}

	item.Initialize()
	return nil
}

func (t *ProjectDb) Load(tx interface{}, id string) (*domain.Project, error) {
	item := &domain.Project{}
	err := t.Base.Executor(tx).SelectOne(item, "select * from project where id=$1", id)
	if err != nil {
		return nil, err
	}

	item.Initialize()
	item.ProjectLead, err = t.UserDb().Load(tx, item.ProjectLead.Id)
	if err != nil {
		return nil, err
	}

	item.Initialize()

	return item, nil
}

func (t *ProjectDb) LoadSimple(tx interface{}, id string) (*domain.Project, error) {
	item := &domain.Project{}
	err := t.Base.Executor(tx).SelectOne(item, "select * from project where id=$1", id)
	if err != nil {
		return nil, err
	}
	item.Initialize()
	return item, nil
}

func (t *ProjectDb) Remove(tx interface{}, id string) (*domain.Project, error) {
	if !t.Base.CurrentSession().User.IsSystemAdministrator {
		return nil, errors.New("db.Base.AccessDenied()")
	}

	item, err := t.Load(tx, id)
	if err != nil {
		return nil, err
	}

	query := "delete from project_role_member where idprojectrole in (" +
		"select id from project_role where idproject=$1)"
	_, err = t.Base.Executor(tx).Exec(query, id)
	if err != nil {
		return nil, err
	}

	query = "delete from project_role where idproject=$1"
	_, err = t.Base.Executor(tx).Exec(query, id)
	if err != nil {
		return nil, err
	}

	_, err = t.Base.Executor(tx).Delete(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (t *ProjectDb) Update(tx interface{}, item *domain.Project) error {
	err := item.Validate()
	if err != nil {
		return err
	}

	oldItem, err := t.Load(tx, item.Id)
	if err != nil {
		return err
	}

	item.DateCreated = oldItem.DateCreated
	item.LastModified = &time.Time{}
	*item.LastModified = time.Now()
	item.Pkey = oldItem.Pkey

	err = t.ValidateDups(tx, item)
	if err != nil {
		return err
	}

	if !t.Base.CurrentSession().User.IsSystemAdministrator && t.Base.CurrentSession().User.Id != oldItem.ProjectLead.Id {
		return errors.New("db.Base.AccessDenied()")
	}

	_, err = t.Base.Executor(tx).Update(item)
	if err != nil {
		return err
	}

	oldItem, err = t.Load(tx, item.Id)
	return err
}

func (t *ProjectDb) Grid(tx interface{}, grid *tecgrid.NgGrid, filter *database.ProjectFilter) error {
	var query string
	var pars []interface{}
	if t.Base.CurrentSession().User.IsSystemAdministrator {
		query = "select * from view_projects where 1 = 1 "
	} else {
		query = `select * from (
			select * 
			from		view_projects p 
			where	not exists( select null from permission_scheme where id = p.idpermissionscheme) 
			union all 
			select	p.* 
			from		view_projects p 
			join		permission_scheme ps on ps.id = p.idpermissionscheme 
			where	exists(	/* user direct access*/ 
					select	null 
					from		permission_scheme_item i 
					join		permission_name n on n.id = i.idpermissionname 
					where	n.name = $1 
					and	i.iduser = $2 
					and	i.idpermissionscheme = ps.id 
					union all 
					/* user / group direct access*/ 
					select	null 
					from		permission_scheme_item i 
					join		permission_name n on n.id = i.idpermissionname 
					join		groups g on g.id = i.idgroup 
					join		user_group gu on gu.idgroup = g.id 
					where	n.name = $1 
					and	gu.iduser = $2 
					and	i.idpermissionscheme = ps.id 
					union all 
					/* user / role access */
					select	null 
					from		permission_scheme_item i 
					join		permission_name n on n.id = i.idpermissionname 
					join		roles r on r.id = i.idrole 
					join		project_role pr on pr.idrole = r.id 
					join		project_role_member pm on pm.idprojectrole = pr.id 
					where	n.name = $1 
					and	pm.iduser = $2 
					and	pr.idproject = p.id 
					union all 
					/*  group / role access */
					select	null 
					from		permission_scheme_item i 
					join		permission_name n on n.id = i.idpermissionname 
					join		roles r on r.id = i.idrole 
					join		project_role pr on pr.idrole = r.id 
					join		project_role_member pm on pm.idprojectrole = pr.id 
					join		groups g on g.id = pm.idgroup 
					join		user_group gu on gu.idgroup = g.id 
					where	n.name = $1 
					and	gu.iduser = $2 
					and	pr.idproject = p.id) 
			) as t  
			where 1 = 1 `

		pars = append(pars, domain.PERMISSION_BROWSE_PROJECT)
		pars = append(pars, t.Base.CurrentSession().User.Id)
	}

	if len(grid.Query) != 0 {
		query += "and lower(name) like '%" + grid.GetQuery() + "%' or lower(projectlead) like '%" + grid.GetQuery() + "%' "
	}

	if filter != nil {
		if len(filter.ProjectLead) != 0 {
			pars = append(pars, filter.ProjectLead)
			query += fmt.Sprintf(" and idprojectlead=$%v ", len(pars))
		}

		if filter.Resolved.Valid {
			if filter.Resolved.Bool {
				query += " and notresolvedcount = 0 "
			} else {
				query += " and notresolvedcount <> 0 "
			}
		}
	}

	grid.MainQuery = query
	var rows []database.ProjectQuery
	fields := strings.Split("id,name,datecreated,pkey,projectlead,percentagecompleted,notresolvedcount,issuecount,begins,ends,permissionscheme", ",")
	return grid.ExecuteSqlParameters(t.Base.GetTransaction(tx), &rows, fields, pars)
}

func (t *ProjectDb) NextNumber(tx interface{}, id string) (*domain.Project, error) {
	item, err := t.LoadSimple(tx, id)
	if err != nil {
		return nil, err
	}
	item.Next++
	_, err = t.Base.Executor(tx).Update(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (t *ProjectDb) ValidateDups(tx interface{}, item *domain.Project) error {
	var (
		err   error
		count int64
	)
	if len(item.Id) != 0 {
		count, err = t.Base.Executor(tx).SelectInt("select count(*) from project where (lower(name)=$1 or lower(pkey)=$2) and id<>$3", strings.ToLower(item.Name), strings.ToLower(item.Pkey), item.Id)
	} else {
		count, err = t.Base.Executor(tx).SelectInt("select count(*) from project where lower(name)=$1 or lower(pkey)=$2", strings.ToLower(item.Name), strings.ToLower(item.Pkey))
	}
	if err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("the project %s or key %s are duplicated ", item.Name, item.Pkey)
	}
	return nil
}

func (t *ProjectDb) CreateMeta(tx interface{}, id string) (*database.ProjectMeta, error) {
	var err error
	item := &database.ProjectMeta{}
	if len(id) != 0 {
		item.Item, err = t.Load(tx, id)
		if err != nil {
			return nil, err
		}

		err = t.RoleCreateAll(tx, item.Item)
		if err != nil {
			return nil, err
		}

		item.ProjectRoles, err = t.Roles(tx, item.Item)
		if err != nil {
			return nil, err
		}
	} else {
		item.Item = &domain.Project{}
	}
	item.Users, err = t.UserDb().List(tx)
	if err != nil {
		return nil, err
	}
	return item, nil
}
