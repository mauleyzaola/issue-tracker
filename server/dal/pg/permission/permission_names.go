package permission

import (
	"github.com/mauleyzaola/issue-tracker/server/domain"
)

func (t *PermissionDb) Names(tx interface{}) ([]domain.PermissionName, error) {
	var items []domain.PermissionName
	_, err := t.Base.Executor(tx).Select(&items, "select * from permission_name")
	if err != nil {
		return nil, err
	}
	for i := range items {
		item := &items[i]
		item.Initialize()
	}
	return items, nil
}

func (t *PermissionDb) AllowedUser(tx interface{}, user *domain.User, issue *domain.Issue, permission *domain.PermissionName) (ok bool, err error) {
	if permission == nil {
		return false, nil
	}

	permissions, err := t.AvailablesUser(tx, user, issue)
	if err != nil {
		return false, err
	}

	for i := range permissions {
		p := &permissions[i]
		if p.Id == permission.Id {
			return true, nil
		} else if p.Name == permission.Name {
			return true, nil
		}
	}
	return false, nil

}

func (t *PermissionDb) AvailablesUser(tx interface{}, user *domain.User, issue *domain.Issue) ([]domain.PermissionName, error) {
	if user == nil {
		return nil, nil
	}

	if issue != nil {
		issue.Initialize()
	}

	currUser, err := t.UserDb().Load(tx, user.Id)
	if err != nil {
		return nil, err
	}

	//es un sysadmin, devolver siempre todos los permisos disponibles
	if currUser.IsSystemAdministrator {
		return t.Names(tx)
	}

	var idPermissionScheme, idProject string
	if issue != nil && len(issue.IdProject.String) != 0 {
		idProject = issue.IdProject.String
		issue.Project, err = t.ProjectDb().LoadSimple(tx, idProject)
		if err != nil {
			idProject = ""
		} else {
			idPermissionScheme = issue.Project.IdPermissionScheme.String
		}
	}

	var query string
	var items []domain.PermissionName

	if len(idProject) != 0 && len(idPermissionScheme) != 0 {
		query = `
		/* hereda la pertenencia por los usuarios de roles en un proyecto */
		select	n.* 
		from		permission_name n 
		join		permission_scheme_item i on i.idpermissionname = n.id 
		join		project_role pr on pr.idrole = i.idrole 
		join		project_role_member pi on pi.idprojectrole = pr.id 
		where	pi.iduser = $1
		and		i.idpermissionscheme = $2
		and		pr.idproject = $3
		union 
		
		/* hereda la pertenencia por los grupos de roles en un proyecto */
		select	n.* 
		from		permission_name n 
		join		permission_scheme_item i on i.idpermissionname = n.id 
		join		project_role pr on pr.idrole = i.idrole 
		join		project_role_member pi on pi.idprojectrole = pr.id 
		join		user_group gu on gu.idgroup = pi.idgroup 
		where	gu.iduser = $1
		and	i.idpermissionscheme = $2
		and	pr.idproject = $3
		union
		
		/* hereda la pertenencia por los usuarios de un permission scheme */
		select	n.* 
		from		permission_name n 
		join		permission_scheme_item i on i.idpermissionname = n.id 
		join		project p on p.idpermissionscheme = i.idpermissionscheme
		where	i.iduser = $1
		and		i.idpermissionscheme = $2
		and		p.id = $3
		union 
		
		/* hereda la pertenencia por los grupos de un permission scheme */
		select	n.* 
		from		permission_name n 
		join		permission_scheme_item i on i.idpermissionname = n.id 
		join		project p on p.idpermissionscheme = i.idpermissionscheme
		join		group_user gu on gu.idgroup = i.idgroup
		where	gu.iduser = $1
		and		i.idpermissionscheme = $2
		and		p.id = $3
		`
		_, err = t.Base.Executor(tx).Select(&items, query, currUser.Id, idPermissionScheme, idProject)
	} else if len(idPermissionScheme) != 0 {
		query = `
		select	n.* 
		from		permission_name n 
		join		permission_scheme_item i on i.idpermissionname = n.id 
		join		user_grouop gu on gu.idgroup = i.idgroup 
		where	gu.iduser = $1
		and		i.idpermissionscheme = $2
		union
		select	n.* 
		from		permission_name n 
		join		permission_scheme_item i on i.idpermissionname = n.id 
		where	i.iduser = $1
		and	i.idpermissionscheme = $2
		`
		_, err = t.Base.Executor(tx).Select(&items, query, currUser.Id, idPermissionScheme)
	} else {
		//si no se pasa un proyecto, ni esta cargando un scheme especifico, devolver todos los permisos
		//no aplica la seguridad fuera de los proyectos
		items, err = t.Names(tx)
	}

	for i := range items {
		p := &items[i]
		p.Initialize()
	}
	return items, err
}
