package issue

import (
	"fmt"
	"strings"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

func (t *IssueDb) Grid(tx interface{}, grid *tecgrid.NgGrid, filter *database.IssueFilter) error {
	var query string
	var pars []interface{}

	if t.Base.CurrentSession().User.IsSystemAdministrator {
		query = "select * from view_issues where 1=1 "
	} else {
		query = `select	* 
			from 
			(select	i.* 
			from	view_issues i 
			where	i.idproject is null 
			union all 
			select	i.* 
			from	view_issues i 
			join	project p on p.id = i.idproject 
			where	p.idpermissionscheme is null 
			union all 
			select	i.* 
			from	view_issues i 
			join	project p on p.id = i.idproject 
			join	permission_scheme ps on ps.id = p.idpermissionscheme 
			where	exists(	
					/* user direct access */ 
					select	null 
					from	permission_scheme_item i 
					join	permission_name n on n.id = i.idpermissionname 
					where	n.name = $1 
					and	i.iduser = $2 
					and	i.idpermissionscheme = ps.id 
					union all 
					
					/* user / group direct access */ 
					select	null 
					from	permission_scheme_item i 
					join	permission_name n on n.id = i.idpermissionname 
					join	user_group g on g.id = i.idgroup 
					join	groups gu on gu.id = g.idgroup
					where	n.name = $1 
					and		g.iduser = $2 
					and	i.idpermissionscheme = ps.id 
					union all 
					
					/* user / role access */ 
					select	null 
					from	permission_scheme_item i 
					join	permission_name n on n.id = i.idpermissionname 
					join	roles r on r.id = i.idrole 
					join	project_role pr on pr.idrole = r.id 
					join	project_role_member pm on pm.idprojectrole = pr.id 
					where	n.name = $1 
					and		pm.iduser = $2 
					and	pr.idproject = p.id 
					union all 
					
					/* group / role access */ 
					select	null 
					from	permission_scheme_item i 
					join	permission_name n on n.id = i.idpermissionname 
					join	roles r on r.id = i.idrole 
					join	project_role pr on pr.idrole = r.id 
					join	project_role_member pm on pm.idprojectrole = pr.id 
					join	groups g on g.id = pm.idgroup 
					join	user_group gu on gu.idgroup = g.id
					where	n.name = $1 
					and		gu.iduser = $2 
					and	pr.idproject = p.id)) as t 
			where	1 = 1 		`
		pars = append(pars, domain.PERMISSION_BROWSE_PROJECT)
		pars = append(pars, t.Base.CurrentSession().User.Id)
	}

	if len(grid.Query) != 0 {
		query += "and (lower(name) like '%" + grid.GetQuery() + "%' or lower(pkey) like '%" + grid.GetQuery() + "%') "
	}
	if filter.Assignee != nil && len(filter.Assignee.Id) != 0 {
		pars = append(pars, filter.Assignee.Id)
		query += fmt.Sprintf(" and idassignee=$%v ", len(pars))
	}
	if filter.Project != nil && len(filter.Project.Id) != 0 {
		pars = append(pars, filter.Project.Id)
		query += fmt.Sprintf(" and idproject=$%v ", len(pars))
	}
	if filter.Reporter != nil && len(filter.Reporter.Id) != 0 {
		pars = append(pars, filter.Reporter.Id)
		query += fmt.Sprintf(" and idreporter=$%v ", len(pars))
	}
	if filter.Priority != nil && len(filter.Priority.Id) != 0 {
		pars = append(pars, filter.Priority.Id)
		query += fmt.Sprintf(" and idpriority=$%v ", len(pars))
	}
	if filter.Status != nil && len(filter.Status.Name) != 0 {
		pars = append(pars, filter.Status.Name)
		query += fmt.Sprintf(" and status=$%v ", len(pars))
	}

	if filter.Resolved.Valid {
		if filter.Resolved.Bool {
			query += " and cancelleddate is null " +
				"and resolveddate is not null "
		} else {
			query += " and cancelleddate is null " +
				"and resolveddate is null "
		}
	}

	if filter.Due.Valid {
		if filter.Due.Bool {
			query += " and duedate < now() "
		} else {
			query += " and duedate >= now() "
		}
	}

	if filter.DueDate != nil {
		pars = append(pars, filter.DueDate.Year())
		query += fmt.Sprintf(" and extract(year from duedate)=$%v ", len(pars))
		pars = append(pars, int(filter.DueDate.Month()))
		query += fmt.Sprintf(" and extract(month from duedate)=$%v ", len(pars))
	}

	grid.MainQuery = query

	var rows []database.IssueGrid
	fields := strings.Split("id,pkey,name,datecreated,lastmodified,duedate,resolveddate,reporter,priority,status,project,assignee,parent", ",")
	return grid.ExecuteSqlParameters(t.Base.GetTransaction(tx), &rows, fields, pars)
}
