package issue

import (
	"errors"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/domain"
)

func (t *IssueDb) AttachmentAdd(tx interface{}, issue *domain.Issue, item *domain.FileItem) (*domain.IssueAttachment, error) {
	oldIssue, err := t.Load(tx, issue.Id, issue.Pkey)
	if err != nil {
		return nil, err
	}

	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_ADD_ATTACHMENT
	ok, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, oldIssue, permission)
	if !ok {
		return nil, errors.New("db.Base.AccessDenied()")
	}

	if !oldIssue.AcceptUpdates() {
		return nil, errors.New("Cannot add attachments to resolved issues")
	}

	newFile, err := t.FileItemDb().Load(tx, item.Id)
	if err != nil {
		return nil, err
	}

	attachment := &domain.IssueAttachment{}
	attachment.FileItem = newFile
	attachment.DateCreated = time.Now()
	attachment.User = t.Base.CurrentSession().User
	attachment.User.Initialize()
	attachment.Issue = issue

	attachment.Issue.Initialize()
	attachment.User.Initialize()

	attachment.Initialize()
	err = t.Base.Executor(tx).Insert(attachment)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (t *IssueDb) AttachmentLoad(tx interface{}, id string) (*domain.IssueAttachment, error) {
	item := &domain.IssueAttachment{}
	err := t.Base.Executor(tx).SelectOne(item, "select * from issue_attachment where id=$1", id)
	if err != nil {
		return nil, err
	}
	item.FileItem, err = t.FileItemDb().Load(tx, item.IdFileItem)
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

func (t *IssueDb) AttachmentRemove(tx interface{}, attachment *domain.IssueAttachment) error {
	oldItem, err := t.AttachmentLoad(tx, attachment.Id)
	if err != nil {
		return err
	}

	if !oldItem.Issue.AcceptUpdates() {
		return errors.New("Cannot remove attachment from a resolved issue")
	}

	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_DELETE_ALL_ATTACHMENT
	canDeleteAll, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, oldItem.Issue, permission)
	if err != nil {
		return err
	}

	permission = &domain.PermissionName{}
	permission.Name = domain.PERMISSION_DELETE_OWN_ATTACHMENT
	canDeleteOwn, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, oldItem.Issue, permission)
	if err != nil {
		return err
	}

	if oldItem.User.Id == t.Base.CurrentSession().User.Id {
		if !canDeleteAll && !canDeleteOwn {
			return errors.New("db.Base.AccessDenied()")
		}
	} else {
		if !canDeleteAll {
			return errors.New("db.Base.AccessDenied()")
		}
	}

	_, err = t.Base.Executor(tx).Delete(oldItem)

	if err != nil {
		return err
	}

	_, err = t.FileItemDb().Remove(tx, oldItem.FileItem.Id)
	if err != nil {
		return err
	}

	attachment = oldItem

	return nil
}

func (t *IssueDb) AttachmentList(tx interface{}, issue *domain.Issue) ([]domain.IssueAttachment, error) {
	oldIssue, err := t.Load(tx, issue.Id, issue.Pkey)
	if err != nil {
		return nil, err
	}
	attachments := make([]domain.IssueAttachment, 0)
	_, err = t.Base.Executor(tx).Select(&attachments, "select * from issue_attachment where idissue=$1 order by datecreated", issue.Id)
	if err != nil {
		return nil, err
	}

	for i := range attachments {
		item := &attachments[i]
		file, _ := t.FileItemDb().Load(tx, item.IdFileItem)
		item.FileItem = file
		item.User, err = t.UserDb().Load(tx, item.IdUser)
		if err != nil {
			return nil, err
		}
		item.Issue = oldIssue
	}
	return attachments, nil
}

func (t *IssueDb) AttachmentRemoveAll(tx interface{}, issue *domain.Issue) error {
	attachments, err := t.AttachmentList(tx, issue)
	if err != nil {
		return nil
	}

	for i := range attachments {
		attachment := &attachments[i]
		err = t.AttachmentRemove(tx, attachment)
		if err != nil {
			return err
		}
	}

	return nil
}
