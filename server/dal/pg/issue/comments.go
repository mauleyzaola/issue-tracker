package issue

import (
	"errors"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/domain"
)

func (t *IssueDb) CommentList(tx interface{}, issue *domain.Issue) ([]domain.IssueComment, error) {
	items := make([]domain.IssueComment, 0)
	type IssueComment struct {
		domain.IssueComment
		Name     string
		LastName string
	}
	var comments []IssueComment

	query := "select c.*, u.name, u.lastname " +
		"from issue_comment c " +
		"join users u on u.id = c.iduser " +
		"where c.idissue = $1 " +
		"order by c.datecreated asc"

	_, err := t.Base.Executor(tx).Select(&comments, query, issue.Id)
	if err != nil {
		return nil, err
	}
	for _, item := range comments {
		i := &domain.IssueComment{Id: item.Id, DateCreated: item.DateCreated, LastModified: item.LastModified, Body: item.Body}
		i.User = &domain.User{Id: item.IdUser, Name: item.Name, LastName: item.LastName}
		i.User.Initialize()
		items = append(items, *i)
	}
	return items, nil
}

func (t *IssueDb) CommentAdd(tx interface{}, comment *domain.IssueComment) error {
	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_ADD_COMMENT
	ok, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, comment.Issue, permission)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("db.Base.AccessDenied()")
	}

	comment.Validate()
	comment.DateCreated = time.Now()
	comment.User = t.Base.CurrentSession().User
	comment.Initialize()
	return t.Base.Executor(tx).Insert(comment)
}

func (t *IssueDb) CommentLoad(tx interface{}, id string) (*domain.IssueComment, error) {
	item := &domain.IssueComment{}
	err := t.Base.Executor(tx).SelectOne(item, "select * from issue_comment where id=$1", id)
	if err != nil {
		return nil, err
	}
	item.User, err = t.UserDb().Load(tx, item.IdUser)
	if err != nil {
		return nil, err
	}
	item.Issue, err = t.Load(tx, item.IdIssue, "")
	if err != nil {
		return nil, err
	}

	item.Initialize()
	return item, nil
}

func (t *IssueDb) CommentUpdate(tx interface{}, comment *domain.IssueComment) error {
	oldItem, err := t.CommentLoad(tx, comment.Id)
	if err != nil {
		return err
	}
	comment.Validate()

	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_EDIT_OWN_COMMENT
	editOwn, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, comment.Issue, permission)
	if err != nil {
		return err
	}

	permission = &domain.PermissionName{}
	permission.Name = domain.PERMISSION_EDIT_ALL_COMMENT
	editAll, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, comment.Issue, permission)
	if err != nil {
		return err
	}

	if oldItem.User.Id == t.Base.CurrentSession().User.Id {
		if !editAll && !editOwn {
			return errors.New("db.Base.AccessDenied()")
		}
	} else {
		if !editAll {
			return errors.New("db.Base.AccessDenied()")
		}
	}

	oldItem.Body = comment.Body
	oldItem.LastModified = &time.Time{}
	*oldItem.LastModified = time.Now()
	_, err = t.Base.Executor(tx).Update(oldItem)
	comment.LastModified = oldItem.LastModified
	if err != nil {
		return err
	}

	comment.Issue = oldItem.Issue

	return nil
}

func (t *IssueDb) CommentRemove(tx interface{}, id string) (*domain.IssueComment, error) {
	item, err := t.CommentLoad(tx, id)
	if err != nil {
		return nil, err
	}

	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_DELETE_OWN_COMMENT
	editOwn, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, item.Issue, permission)
	if err != nil {
		return nil, err
	}

	permission = &domain.PermissionName{}
	permission.Name = domain.PERMISSION_DELETE_ALL_COMMENT
	editAll, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, item.Issue, permission)
	if err != nil {
		return nil, err
	}

	if item.User.Id == t.Base.CurrentSession().User.Id {
		if !editAll && !editOwn {
			return nil, errors.New("db.Base.AccessDenied()")
		}
	} else {
		if !editAll {
			return nil, errors.New("db.Base.AccessDenied()")
		}
	}

	_, err = t.Base.Executor(tx).Delete(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (t *IssueDb) CommentRemoveAll(tx interface{}, issue *domain.Issue) error {
	comments, err := t.CommentList(tx, issue)
	if err != nil {
		return err
	}
	for i := range comments {
		comment := &comments[i]
		_, err = t.CommentRemove(tx, comment.Id)
		if err != nil {
			return err
		}
	}
	return nil
}
