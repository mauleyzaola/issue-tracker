package issue

import (
	"fmt"

	"github.com/mauleyzaola/issue-tracker/server/operations/database"
)

func (t *IssueDb) DashbordIssueGroupByAssignee(tx interface{}, filter *database.IssueGroupFilter) ([]database.IssueGroupResult, error) {
	var pars []interface{}
	var err error
	var items []database.IssueGroupResult
	query := "select	idassignee as id, assignee as name, count(*) as rowcount " +
		"from	view_issues " +
		"where	resolveddate is null " +
		"and		cancelleddate is null " +
		"and		idassignee is not null "
	if filter != nil {
		if len(filter.Project) != 0 {
			pars = append(pars, filter.Project)
			query += fmt.Sprintf("and idproject = $%v ", len(pars))
		}
		if len(filter.Assignee) != 0 {
			pars = append(pars, filter.Assignee)
			query += fmt.Sprintf("and idassignee = $%v ", len(pars))
		}
		if len(filter.Reporter) != 0 {
			pars = append(pars, filter.Reporter)
			query += fmt.Sprintf("and idreporter = $%v ", len(pars))
		}
		if len(filter.Status) != 0 {
			pars = append(pars, filter.Status)
			query += fmt.Sprintf("and status = $%v ", len(pars))
		}
	}
	query += "group by idassignee,assignee " +
		"order by 2"
	if pars == nil {
		_, err = t.Base.Executor(tx).Select(&items, query)
	} else {
		_, err = t.Base.Executor(tx).Select(&items, query, pars...)
	}
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]database.IssueGroupResult, 0)
		return items, nil
	}

	for i := range items {
		item := &items[i]
		item.DataType = "assignee"
	}
	return items, nil
}

func (t *IssueDb) DashbordIssueGroupByReporter(tx interface{}, filter *database.IssueGroupFilter) ([]database.IssueGroupResult, error) {
	var err error
	var pars []interface{}
	var items []database.IssueGroupResult

	query := "select	idreporter as id, reporter as name, count(*) as rowcount " +
		"from	view_issues " +
		"where	resolveddate is null " +
		"and		cancelleddate is null "
	if filter != nil {
		if len(filter.Project) != 0 {
			pars = append(pars, filter.Project)
			query += fmt.Sprintf("and idproject = $%v ", len(pars))
		}
		if len(filter.Assignee) != 0 {
			pars = append(pars, filter.Assignee)
			query += fmt.Sprintf("and idassignee = $%v ", len(pars))
		}
		if len(filter.Reporter) != 0 {
			pars = append(pars, filter.Reporter)
			query += fmt.Sprintf("and idreporter = $%v ", len(pars))
		}
		if len(filter.Status) != 0 {
			pars = append(pars, filter.Status)
			query += fmt.Sprintf("and status = $%v ", len(pars))
		}
	}
	query += "group by idreporter,reporter " +
		"order by 2"
	if pars == nil {
		_, err = t.Base.Executor(tx).Select(&items, query)
	} else {
		_, err = t.Base.Executor(tx).Select(&items, query, pars...)
	}
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]database.IssueGroupResult, 0)
		return items, nil
	}

	for i := range items {
		item := &items[i]
		item.DataType = "reporter"
	}
	return items, nil
}

func (t *IssueDb) DashboardIssueGroupByPriority(tx interface{}, filter *database.IssueGroupFilter) ([]database.IssueGroupResult, error) {
	var err error
	var pars []interface{}
	var items []database.IssueGroupResult

	query := "select	idpriority as id, priority as name, count(*) as rowcount " +
		"from	view_issues " +
		"where	resolveddate is null " +
		"and		cancelleddate is null "
	if filter != nil {
		if len(filter.Project) != 0 {
			pars = append(pars, filter.Project)
			query += fmt.Sprintf("and idproject = $%v ", len(pars))
		}
		if len(filter.Assignee) != 0 {
			pars = append(pars, filter.Assignee)
			query += fmt.Sprintf("and idassignee = $%v ", len(pars))
		}
		if len(filter.Reporter) != 0 {
			pars = append(pars, filter.Reporter)
			query += fmt.Sprintf("and idreporter = $%v ", len(pars))
		}
		if len(filter.Status) != 0 {
			pars = append(pars, filter.Status)
			query += fmt.Sprintf("and status = $%v ", len(pars))
		}
	}
	query += "group by idpriority,priority " +
		"order by 2"
	if pars == nil {
		_, err = t.Base.Executor(tx).Select(&items, query)
	} else {
		_, err = t.Base.Executor(tx).Select(&items, query, pars...)
	}

	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]database.IssueGroupResult, 0)
		return items, nil
	}

	for i := range items {
		item := &items[i]
		item.DataType = "priority"
	}
	return items, nil
}

func (t *IssueDb) DashboardIssueGroupByStatus(tx interface{}, filter *database.IssueGroupFilter) ([]database.IssueGroupResult, error) {
	var pars []interface{}
	var err error
	var items []database.IssueGroupResult
	query := "select	'' as id, status as name, count(*) as rowcount " +
		"from	view_issues " +
		"where	resolveddate is null " +
		"and		cancelleddate is null "
	if filter != nil {
		if len(filter.Project) != 0 {
			pars = append(pars, filter.Project)
			query += fmt.Sprintf("and idproject = $%v ", len(pars))
		}
		if len(filter.Assignee) != 0 {
			pars = append(pars, filter.Assignee)
			query += fmt.Sprintf("and idassignee = $%v ", len(pars))
		}
		if len(filter.Reporter) != 0 {
			pars = append(pars, filter.Reporter)
			query += fmt.Sprintf("and idreporter = $%v ", len(pars))
		}
		if len(filter.Status) != 0 {
			pars = append(pars, filter.Status)
			query += fmt.Sprintf("and status = $%v ", len(pars))
		}
	}
	query += "group by status " +
		"order by 2"
	if pars == nil {
		_, err = t.Base.Executor(tx).Select(&items, query)
	} else {
		_, err = t.Base.Executor(tx).Select(&items, query, pars...)
	}

	if err != nil {
		return nil, err
	}

	for i := range items {
		item := &items[i]
		item.DataType = "status"
	}
	return items, nil
}

func (t *IssueDb) DashboardIssueGroupByProject(tx interface{}, filter *database.IssueGroupFilter) ([]database.IssueGroupResult, error) {
	var pars []interface{}
	var err error
	var items []database.IssueGroupResult

	query := "select	idproject as id, project as name, count(*) as rowcount " +
		"from	view_issues " +
		"where	resolveddate is null " +
		"and		cancelleddate is null " +
		"and		idproject is not null "
	if filter != nil {
		if len(filter.Project) != 0 {
			pars = append(pars, filter.Project)
			query += fmt.Sprintf("and idproject = $%v ", len(pars))
		}
		if len(filter.Assignee) != 0 {
			pars = append(pars, filter.Assignee)
			query += fmt.Sprintf("and idassignee = $%v ", len(pars))
		}
		if len(filter.Reporter) != 0 {
			pars = append(pars, filter.Reporter)
			query += fmt.Sprintf("and idreporter = $%v ", len(pars))
		}
		if len(filter.Status) != 0 {
			pars = append(pars, filter.Status)
			query += fmt.Sprintf("and status = $%v ", len(pars))
		}
	}
	query += "group by idproject,project " +
		"order by 2"

	if pars == nil {
		_, err = t.Base.Executor(tx).Select(&items, query)
	} else {
		_, err = t.Base.Executor(tx).Select(&items, query, pars...)
	}

	if err != nil {
		return nil, err
	}

	for i := range items {
		item := &items[i]
		item.DataType = "project"
	}
	return items, nil
}

func (t *IssueDb) DashbordIssueGroupByDueDate(tx interface{}, filter *database.IssueGroupFilter) ([]database.IssueGroupResult, error) {
	var pars []interface{}
	var err error
	var items []database.IssueGroupResult
	query := "select name, count(*) as rowcount " +
		"from " +
		"(select substr(cast(date_trunc('month', duedate) as text),1,7) as name " +
		"from issue " +
		"where resolveddate is null " +
		"and cancelleddate is null "
	if filter != nil {
		if len(filter.Project) != 0 {
			pars = append(pars, filter.Project)
			query += fmt.Sprintf("and idproject = $%v ", len(pars))
		}
		if len(filter.Assignee) != 0 {
			pars = append(pars, filter.Assignee)
			query += fmt.Sprintf("and idassignee = $%v ", len(pars))
		}
		if len(filter.Reporter) != 0 {
			pars = append(pars, filter.Reporter)
			query += fmt.Sprintf("and idreporter = $%v ", len(pars))
		}
		if len(filter.Status) != 0 {
			pars = append(pars, filter.Status)
			query += fmt.Sprintf("and status = $%v ", len(pars))
		}
	}

	query += ") as t " +
		"group by name " +
		"order by name"

	if pars == nil {
		_, err = t.Base.Executor(tx).Select(&items, query)
	} else {
		_, err = t.Base.Executor(tx).Select(&items, query, pars...)
	}

	if err != nil {
		return nil, err
	}
	if items == nil {
		items = make([]database.IssueGroupResult, 0)
		return items, nil
	}

	for i := range items {
		item := &items[i]
		item.DataType = "dueDate"
	}
	return items, nil
}
