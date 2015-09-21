package domain

import ()

type PermissionName struct {
	Meta *DocumentMetadata `json:"meta" db:"-"`
	Id   string            `json:"id"`
	Name string            `json:"name"`
}

const (
	PERMISSION_BROWSE_PROJECT string = "BROWSE_PROJECT"

	PERMISSION_CREATE_ISSUE  string = "CREATE_ISSUE"
	PERMISSION_EDIT_ISSUE    string = "EDIT_ISSUE"
	PERMISSION_DELETE_ISSUE  string = "DELETE_ISSUE"
	PERMISSION_RESOLVE_ISSUE string = "RESOLVE_ISSUE"

	PERMISSION_ASSIGN_USER     string = "ASSIGN_USER"
	PERMISSION_CHANGE_REPORTER string = "CHANGE_REPORTER"
	PERMISSION_CHANGE_DUEDATE  string = "CHANGE_DUEDATE"

	PERMISSION_ADD_ATTACHMENT        string = "ADD_ATTACHMENT"
	PERMISSION_DELETE_OWN_ATTACHMENT string = "DELETE_OWN_ATTACHMENT"
	PERMISSION_DELETE_ALL_ATTACHMENT string = "DELETE_ALL_ATTACHMENT"

	PERMISSION_ADD_COMMENT        string = "ADD_COMMENT"
	PERMISSION_EDIT_OWN_COMMENT   string = "EDIT_OWN_COMMENT"
	PERMISSION_DELETE_OWN_COMMENT string = "DELETE_OWN_COMMENT"

	PERMISSION_EDIT_ALL_COMMENT   string = "EDIT_ALL_COMMENT"
	PERMISSION_DELETE_ALL_COMMENT string = "DELETE_ALL_COMMENT"

	PERMISSION_ADD_WORKLOG        string = "ADD_WORKLOG"
	PERMISSION_DELETE_OWN_WORKLOG string = "DELETE_OWN_WORKLOG"
	PERMISSION_DELETE_ALL_WORKLOG string = "DELETE_ALL_WORKLOG"

	PERMISSION_SUBSCRIBE_ISSUE    string = "SUBSCRIBE_ISSUE"
	PERMISSION_UNSUBSCRIBE_ISSUE  string = "UNSUBSCRIBE_ISSUE"
	PERMISSION_SUBSCRIBE_OTHERS   string = "SUBSCRIBE_OTHERS"
	PERMISSION_UNSUBSCRIBE_OTHERS string = "UNSUBSCRIBE_OTHERS"
)

//returns all the permission names
//the key is the constant name, and the content is the db name
func (t *PermissionName) Permissions() map[string]string {
	m := make(map[string]string)
	m["PERMISSION_BROWSE_PROJECT"] = "BROWSE_PROJECT"

	m["PERMISSION_CREATE_ISSUE"] = "CREATE_ISSUE"
	m["PERMISSION_EDIT_ISSUE"] = "EDIT_ISSUE"
	m["PERMISSION_DELETE_ISSUE"] = "DELETE_ISSUE"
	m["PERMISSION_RESOLVE_ISSUE"] = "RESOLVE_ISSUE"

	m["PERMISSION_ASSIGN_USER"] = "ASSIGN_USER"
	m["PERMISSION_CHANGE_REPORTER"] = "CHANGE_REPORTER"
	m["PERMISSION_CHANGE_DUEDATE"] = "CHANGE_DUEDATE"

	m["PERMISSION_ADD_ATTACHMENT"] = "ADD_ATTACHMENT"
	m["PERMISSION_DELETE_OWN_ATTACHMENT"] = "DELETE_OWN_ATTACHMENT"
	m["PERMISSION_DELETE_ALL_ATTACHMENT"] = "DELETE_ALL_ATTACHMENT"

	m["PERMISSION_ADD_COMMENT"] = "ADD_COMMENT"
	m["PERMISSION_EDIT_OWN_COMMENT"] = "EDIT_OWN_COMMENT"
	m["PERMISSION_DELETE_OWN_COMMENT"] = "DELETE_OWN_COMMENT"

	m["PERMISSION_EDIT_ALL_COMMENT"] = "EDIT_ALL_COMMENT"
	m["PERMISSION_DELETE_ALL_COMMENT"] = "DELETE_ALL_COMMENT"

	m["PERMISSION_ADD_WORKLOG"] = "ADD_WORKLOG"
	m["PERMISSION_DELETE_OWN_WORKLOG"] = "DELETE_OWN_WORKLOG"
	m["PERMISSION_DELETE_ALL_WORKLOG"] = "DELETE_ALL_WORKLOG"

	m["PERMISSION_SUBSCRIBE_ISSUE"] = "SUBSCRIBE_ISSUE"
	m["PERMISSION_UNSUBSCRIBE_ISSUE"] = "UNSUBSCRIBE_ISSUE"
	m["PERMISSION_SUBSCRIBE_OTHERS"] = "SUBSCRIBE_OTHERS"
	m["PERMISSION_UNSUBSCRIBE_OTHERS"] = "UNSUBSCRIBE_OTHERS"

	return m
}

func (u *PermissionName) Initialize() {
	u.Meta = &DocumentMetadata{}

	if u.Name == PERMISSION_BROWSE_PROJECT {
		u.Meta.FriendName = "Browse Project"
	} else if u.Name == PERMISSION_CREATE_ISSUE {
		u.Meta.FriendName = "Create Issue"
	} else if u.Name == PERMISSION_EDIT_ISSUE {
		u.Meta.FriendName = "Edit Issue"
	} else if u.Name == PERMISSION_DELETE_ISSUE {
		u.Meta.FriendName = "Delete Issue"
	} else if u.Name == PERMISSION_RESOLVE_ISSUE {
		u.Meta.FriendName = "Resolve Issue"
	} else if u.Name == PERMISSION_ASSIGN_USER {
		u.Meta.FriendName = "Change Assignee"
	} else if u.Name == PERMISSION_CHANGE_REPORTER {
		u.Meta.FriendName = "Change Reporter"
	} else if u.Name == PERMISSION_CHANGE_DUEDATE {
		u.Meta.FriendName = "Change Due Date"
	} else if u.Name == PERMISSION_ADD_ATTACHMENT {
		u.Meta.FriendName = "Add Attachments"
	} else if u.Name == PERMISSION_DELETE_OWN_ATTACHMENT {
		u.Meta.FriendName = "Delete Own Attachments"
	} else if u.Name == PERMISSION_DELETE_ALL_ATTACHMENT {
		u.Meta.FriendName = "Delete Any Attachment"
	} else if u.Name == PERMISSION_ADD_COMMENT {
		u.Meta.FriendName = "Add Comment"
	} else if u.Name == PERMISSION_EDIT_OWN_COMMENT {
		u.Meta.FriendName = "Edit Own Comments"
	} else if u.Name == PERMISSION_DELETE_OWN_COMMENT {
		u.Meta.FriendName = "Delete Own Comments"
	} else if u.Name == PERMISSION_EDIT_ALL_COMMENT {
		u.Meta.FriendName = "Edit Any Comment"
	} else if u.Name == PERMISSION_DELETE_ALL_COMMENT {
		u.Meta.FriendName = "Delete Any Comment"
	} else if u.Name == PERMISSION_ADD_WORKLOG {
		u.Meta.FriendName = "Add Wornlog"
	} else if u.Name == PERMISSION_DELETE_OWN_WORKLOG {
		u.Meta.FriendName = "Delete Own Worklog"
	} else if u.Name == PERMISSION_DELETE_ALL_WORKLOG {
		u.Meta.FriendName = "Delete Any Worklog"
	} else if u.Name == PERMISSION_SUBSCRIBE_ISSUE {
		u.Meta.FriendName = "Subscribe to Issue"
	} else if u.Name == PERMISSION_SUBSCRIBE_OTHERS {
		u.Meta.FriendName = "Subscribe others to Issues"
	} else if u.Name == PERMISSION_UNSUBSCRIBE_ISSUE {
		u.Meta.FriendName = "Unsubscribe from Issue"
	} else if u.Name == PERMISSION_UNSUBSCRIBE_OTHERS {
		u.Meta.FriendName = "Unsubscribe others from Issues"
	}
}
