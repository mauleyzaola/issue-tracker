package issue

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

func (t *IssueDb) Subscribers(tx interface{}, issue *domain.Issue, excludeCurrentUser bool) (items []domain.User, err error) {
	ids := []string{}
	if excludeCurrentUser {
		_, err = t.Base.Executor(tx).Select(&ids, "select iduser from issue_subscription where idissue=$1 and iduser <> $2", issue.Id, t.Base.CurrentSession().User.Id)
	} else {
		_, err = t.Base.Executor(tx).Select(&ids, "select iduser from issue_subscription where idissue=$1", issue.Id)
	}

	for _, i := range ids {
		user, e := t.UserDb().Load(tx, i)
		if e != nil {
			err = e
			return
		}
		items = append(items, *user)
	}
	return
}

func (t *IssueDb) MySubscriptions(tx interface{}, grid *tecgrid.NgGrid) error {
	return t.SubscriptionsUser(tx, grid, t.Base.CurrentSession().User)
}

func (t *IssueDb) SubscriptionsUser(tx interface{}, grid *tecgrid.NgGrid, user *domain.User) error {

	var pars []interface{}
	query := `
		select * 
		from view_issues i 
		where exists( 	select null 
						from issue_subscription 
						where iduser=$1 
						and idissue = i.id) 
	`
	pars = append(pars, user.Id)
	if len(grid.GetQuery()) != 0 {
		query += "and lower(name) like '%" + grid.GetQuery() + "%' "
	}
	fields := strings.Split("id,pkey,name,datecreated,lastmodified,duedate,resolveddate,reporter,priority,status,project,assignee,parent", ",")
	var rows []database.IssueGrid
	grid.MainQuery = query
	return grid.ExecuteSqlParameters(t.Base.GetTransaction(tx), &rows, fields, pars)
}

func (t *IssueDb) SubscriptionToggle(tx interface{}, issue *domain.Issue) (selected bool, err error) {
	selected, err = t.SubscriptionToggleUser(tx, issue, t.Base.CurrentSession().User)
	return
}

func (t *IssueDb) SubscriptionToggleUser(tx interface{}, issue *domain.Issue, user *domain.User) (selected bool, err error) {
	var oldItem domain.IssueSubscription
	err = t.Base.Executor(tx).SelectOne(&oldItem, "select * from issue_subscription where idissue=$1 and iduser=$2", issue.Id, user.Id)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	oldIssue, err := t.Load(tx, issue.Id, "")
	if err != nil {
		return
	}
	permission := &domain.PermissionName{}

	if len(oldItem.Id) != 0 {
		if user.Id != t.Base.CurrentSession().User.Id {
			permission.Name = domain.PERMISSION_UNSUBSCRIBE_OTHERS
		} else {
			permission.Name = domain.PERMISSION_UNSUBSCRIBE_ISSUE
		}
	} else {
		if user.Id != t.Base.CurrentSession().User.Id {
			permission.Name = domain.PERMISSION_SUBSCRIBE_OTHERS
		} else {
			permission.Name = domain.PERMISSION_SUBSCRIBE_ISSUE
		}
	}
	allowed, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, oldIssue, permission)
	if !allowed {
		err = errors.New("db.Base.AccessDenied()")
		return
	}

	if len(oldItem.Id) == 0 {
		oldItem.IdUser = user.Id
		oldItem.IdIssue = issue.Id
		err = t.Base.Executor(tx).Insert(&oldItem)
		selected = true
		return
	} else {
		_, err = t.Base.Executor(tx).Delete(&oldItem)
		selected = false
		return
	}
}

func (t *IssueDb) IsSubscribed(tx interface{}, issue *domain.Issue, user *domain.User) (ok bool, err error) {
	rowCount, err := t.Base.Executor(tx).SelectInt("select count(*) from issue_subscription where idissue=$1 and iduser=$2", issue.Id, user.Id)
	if err != nil {
		return false, err
	}
	return rowCount != 0, nil
}

func (t *IssueDb) SubscriptionAdd(tx interface{}, issue *domain.Issue, user *domain.User) error {
	ok, err := t.IsSubscribed(tx, issue, user)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	item := &domain.IssueSubscription{Issue: issue, User: user}
	item.Initialize()
	return t.Base.Executor(tx).Insert(item)
}

func (t *IssueDb) SubscribersIssue(tx interface{}, issue *domain.Issue) (selected []domain.User, unselected []domain.User, err error) {
	query := `
	select	u.*
	from		users u
	join		issue_subscription s on s.iduser = u.id
	where	s.idissue = $1
	order by u.name, u.lastname
	`

	_, err = t.Base.Executor(tx).Select(&selected, query, issue.Id)
	if err != nil {
		return
	}

	query = `
	select	u.*
	from		users u
	where	not exists(	select 	null 
						from		issue_subscription
						where	idissue = $1
						and		iduser = u.id)
	order by u.name, u.lastname
	`
	_, err = t.Base.Executor(tx).Select(&unselected, query, issue.Id)
	return
}
