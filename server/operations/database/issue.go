package database

import (
	"database/sql"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/utils/tecgrid"
)

//Lambda that can be passed as required when an Issue status has changed
type IssueStatusFn func(tx interface{}, oldItem *domain.Issue, nextStep *domain.WorkflowStep, oldStatus *domain.Status, newStatus *domain.Status) error

type Issue interface {
	FileItemDb() FileItem
	SetFileItemDb(item *FileItem)

	PermissionDb() Permission
	SetPermissionDb(item *Permission)

	PriorityDb() Priority
	SetPriorityDb(item *Priority)

	ProjectDb() Project
	SetProjectDb(item *Project)

	StatusDb() Status
	SetStatusDb(item *Status)

	UserDb() User
	SetUserDb(item *User)

	//Adds an attachment to an issue
	AttachmentAdd(tx interface{}, issue *domain.Issue, item *domain.FileItem) (*domain.IssueAttachment, error)

	//Loads an Attachment from database along with inner objects
	AttachmentLoad(tx interface{}, id string) (*domain.IssueAttachment, error)

	//Removes an Attachment from database along with the FileItem it is related
	AttachmentRemove(tx interface{}, attachment *domain.IssueAttachment) error

	//Returns a list of attachments for a given issue
	AttachmentList(tx interface{}, issue *domain.Issue) ([]domain.IssueAttachment, error)

	//Removes all the attachments from a given issue
	AttachmentRemoveAll(tx interface{}, issue *domain.Issue) error

	//Returns the list of comments for a given issue
	CommentList(tx interface{}, issue *domain.Issue) ([]domain.IssueComment, error)

	//Adds a comment to a given issue
	CommentAdd(tx interface{}, comment *domain.IssueComment) error

	//Retrieves a comment from database by its id
	CommentLoad(tx interface{}, id string) (*domain.IssueComment, error)

	//Updates a comment in database
	CommentUpdate(tx interface{}, comment *domain.IssueComment) error

	//Removes a comment from database
	CommentRemove(tx interface{}, id string) (*domain.IssueComment, error)

	//Removes all the comments that are related with a given issue
	CommentRemoveAll(tx interface{}, issue *domain.Issue) error

	//Returns a query grouping pending issues by assignee
	DashbordIssueGroupByAssignee(tx interface{}, filter *IssueGroupFilter) ([]IssueGroupResult, error)

	//Returns a query grouping pending issues by reporter
	DashbordIssueGroupByReporter(tx interface{}, filter *IssueGroupFilter) ([]IssueGroupResult, error)

	//Returns a query grouping pending issues by priority
	DashboardIssueGroupByPriority(tx interface{}, filter *IssueGroupFilter) ([]IssueGroupResult, error)

	//Returns a query grouping pending issues by status
	DashboardIssueGroupByStatus(tx interface{}, filter *IssueGroupFilter) ([]IssueGroupResult, error)

	//Returns a query grouping pending issues by project
	DashboardIssueGroupByProject(tx interface{}, filter *IssueGroupFilter) ([]IssueGroupResult, error)

	//Returns a query grouping pending issues by year/month of dueDate
	DashbordIssueGroupByDueDate(tx interface{}, filter *IssueGroupFilter) ([]IssueGroupResult, error)

	//Changes the status of a given issue and takes as parameter a func to process if there are no errors
	StatusChange(tx interface{}, issue *domain.Issue, status *domain.Status, fn IssueStatusFn) error

	//Returns the list of users subscribed to a given issue
	Subscribers(tx interface{}, issue *domain.Issue, excludeCurrentUser bool) (items []domain.User, err error)

	//Returns a list of the subscriptions the current session users' has been subscribed
	MySubscriptions(tx interface{}, grid *tecgrid.NgGrid) error

	//Processes a grid with all the subscriptions for a given user
	SubscriptionsUser(tx interface{}, grid *tecgrid.NgGrid, user *domain.User) error

	//Subscribes or unsubscribes the connected user with an Issue
	SubscriptionToggle(tx interface{}, issue *domain.Issue) (selected bool, err error)

	//Subscribes or unsubscribes a given user with an Issue
	SubscriptionToggleUser(tx interface{}, issue *domain.Issue, user *domain.User) (selected bool, err error)

	//Returns true if the given user is subscribed, false if he is not
	IsSubscribed(tx interface{}, issue *domain.Issue, user *domain.User) (ok bool, err error)

	//Returns a list of users subscribed/unsubscribed to a given issue
	SubscribersIssue(tx interface{}, issue *domain.Issue) (selected []domain.User, unselected []domain.User, err error)

	//Adds a subscription to a given issue
	SubscriptionAdd(tx interface{}, issue *domain.Issue, user *domain.User) error

	//Creates an issue
	Create(tx interface{}, item *domain.Issue, parent string) error

	//Loads an issue along with its inner objects
	Load(tx interface{}, id string, pkey string) (*domain.Issue, error)

	//Removes an issue from database
	Remove(tx interface{}, id string) (*domain.Issue, error)

	//Returns the id of the root issue for a given issue
	FindRoot(tx interface{}, id string) (root string, err error)

	//Updates an issue in database
	Update(tx interface{}, u *domain.Issue) error

	//Processes the grid query for issues applying permissions when requiered
	Grid(tx interface{}, grid *tecgrid.NgGrid, filter *IssueFilter) error

	//Gets all the slibings at the first level for a given issue
	Children(tx interface{}, issue *domain.Issue) ([]IssueGrid, error)

	//Updates an Issue moving the project from where it belongs. A new pkey is assigned and new
	//permissions will apply depending on the new project configuration
	MoveProject(tx interface{}, issue *domain.Issue, target *domain.Project) (*domain.Issue, error)
}

type IssueGroupResult struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	RowCount int64  `json:"rowCount"`
	DataType string `json:"dataType"`
}

type IssueGroupFilter struct {
	Project  string
	Assignee string
	Reporter string
	Status   string
}

type IssueFilter struct {
	domain.Issue
	Parent   *domain.Issue `json:"parent"`
	Due      sql.NullBool
	Resolved sql.NullBool
	DueDate  *time.Time
}

type IssueGrid struct {
	domain.Issue
	Reporter  string         `json:"reporter"`
	Priority  string         `json:"priority"`
	Status    string         `json:"status"`
	Project   string         `json:"project"`
	Assignee  string         `json:"assignee"`
	ParentKey string         `json:"parent" db:"parent"`
	IdParent  sql.NullString `json:"-"`
}
