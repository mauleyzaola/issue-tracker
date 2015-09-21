package issue

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
)

type IssueDb struct {
	Base         *pg.Db
	fileItemDb   *database.FileItem
	permissionDb *database.Permission
	priorityDb   *database.Priority
	projectDb    *database.Project
	statusDb     *database.Status
	userDb       *database.User
}

func New(db database.Db) *IssueDb {
	base := db.(*pg.Db)
	return &IssueDb{Base: base}
}

func (t *IssueDb) FileItemDb() database.FileItem {
	return *t.fileItemDb
}

func (t *IssueDb) SetFileItemDb(item *database.FileItem) {
	t.fileItemDb = item
}

func (t *IssueDb) PermissionDb() database.Permission {
	return *t.permissionDb
}

func (t *IssueDb) SetPermissionDb(item *database.Permission) {
	t.permissionDb = item
}

func (t *IssueDb) PriorityDb() database.Priority {
	return *t.priorityDb
}

func (t *IssueDb) SetPriorityDb(item *database.Priority) {
	t.priorityDb = item
}

func (t *IssueDb) ProjectDb() database.Project {
	return *t.projectDb
}

func (t *IssueDb) SetProjectDb(item *database.Project) {
	t.projectDb = item
}

func (t *IssueDb) StatusDb() database.Status {
	return *t.statusDb
}

func (t *IssueDb) SetStatusDb(item *database.Status) {
	t.statusDb = item
}

func (t *IssueDb) UserDb() database.User {
	return *t.userDb
}

func (t *IssueDb) SetUserDb(item *database.User) {
	t.userDb = item
}

func (t *IssueDb) StatusChange(tx interface{}, issue *domain.Issue, status *domain.Status, fn database.IssueStatusFn) error {
	return nil
}

func (t *IssueDb) Create(tx interface{}, item *domain.Issue, parent string) error {
	var err error
	item.Initialize()

	if item.Reporter == nil || len(item.Reporter.Id) == 0 {
		if item.Project != nil && len(item.Project.Id) != 0 {
			item.Project, err = t.ProjectDb().Load(tx, item.Project.Id)
			if err != nil {
				return err
			}
			if item.Project.ProjectLead != nil && len(item.Project.ProjectLead.Id) != 0 {
				item.Reporter = item.Project.ProjectLead
			}
		}
		if item.Reporter == nil || len(item.Reporter.Id) == 0 {
			item.Reporter = t.Base.CurrentSession().User
		}
	}

	if item.Priority == nil || len(item.Priority.Id) == 0 {
		item.Priority, err = t.PriorityDb().GetFirst(tx)
		if err != nil {
			return err
		}
	}

	err = item.Validate()
	if err != nil {
		return err
	}

	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_CREATE_ISSUE
	ok, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, item, permission)

	if !ok {
		return errors.New("db.Base.AccessDenied()")
	}

	if len(item.IdProject.String) != 0 {
		project, e := t.ProjectDb().NextNumber(tx, item.Project.Id)
		if e != nil {
			return e
		}
		item.Pkey = project.IssueKey()
	} else {
		next, e := t.Base.SequenceNextValue(tx, operations.ISSUE_SEQUENCE_NAME)
		if e != nil {
			return e
		}
		item.Pkey = fmt.Sprintf(operations.ISSUE_NO_PROJECT_MASK, next)
	}

	//calcular el primer status que se puede asignar a este issue segun el workflow
	statuses, err := t.StatusDb().WorkflowStepAvailableStatus(tx, item.Workflow, nil)
	if err != nil {
		return err
	}

	if len(statuses) != 1 {
		return errors.New("No se encuentra el primer estado para el workflow seleccionado")
	}

	item.IdStatus = statuses[0].NextStatus.Id

	//si se pasa un parent como parametro, crear la relacion correspondiente
	if len(parent) != 0 {
		item.IdParent.Valid = true
		item.IdParent.String = parent
	}

	err = t.Base.Executor(tx).Insert(item)
	if err != nil {
		return err
	}

	err = t.SubscriptionAdd(tx, item, item.Reporter)
	if err != nil {
		return err
	}
	if item.Assignee != nil && len(item.Assignee.Id) != 0 {
		err = t.SubscriptionAdd(tx, item, item.Assignee)
		if err != nil {
			return err
		}
	}
	return err
}

func (t *IssueDb) Load(tx interface{}, id string, pkey string) (*domain.Issue, error) {
	var err error
	item := &domain.Issue{}
	if len(id) != 0 {
		err = t.Base.Executor(tx).SelectOne(item, "select * from issue where id=$1", id)
	} else {
		err = t.Base.Executor(tx).SelectOne(item, "select * from issue where pkey=$1", pkey)
	}

	if err != nil {
		return nil, err
	}

	item.Initialize()

	permission := &domain.PermissionName{Name: domain.PERMISSION_BROWSE_PROJECT}
	ok, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, item, permission)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("db.Base.AccessDenied()")
	}

	if item.IdAssignee.Valid {
		item.Assignee, err = t.UserDb().Load(tx, item.IdAssignee.String)
		if err != nil {
			return nil, err
		}
	}

	item.Priority, err = t.PriorityDb().Load(tx, item.IdPriority)
	if err != nil {
		return nil, err
	}
	item.Reporter, err = t.UserDb().Load(tx, item.IdReporter)
	if err != nil {
		return nil, err
	}
	item.Status, err = t.StatusDb().Load(tx, item.IdStatus)
	if err != nil {
		return nil, err
	}
	item.Workflow, err = t.StatusDb().WorkflowLoad(tx, item.IdWorkflow)
	if err != nil {
		return nil, err
	}

	if item.IdProject.Valid {
		item.Project, err = t.ProjectDb().Load(tx, item.IdProject.String)
		if err != nil {
			return nil, err
		}
	}

	return item, nil
}

func (t *IssueDb) Remove(tx interface{}, id string) (*domain.Issue, error) {
	item, err := t.Load(tx, id, "")
	if err != nil {
		return nil, err
	}

	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_DELETE_ISSUE
	ok, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, item, permission)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("db.Base.AccessDenied()")
	}

	children, err := t.Base.Executor(tx).SelectInt("select count(*) from issue where idparent=$1", item.Id)
	if children != 0 {
		return nil, errors.New(fmt.Sprintf("No se puede eliminar la tarea porque hay otras %d subtareas que dependen de ella", children))
	}

	children, err = t.Base.Executor(tx).SelectInt("select count(*) from issue_relation where idissue=$1", item.Id)
	if children != 0 {
		return nil, errors.New("No se puede eliminar la tarea porque hay otro proceso dependiente de ella")
	}

	query := "delete from issue_subscription where idissue=$1"
	_, err = t.Base.Executor(tx).Exec(query, item.Id)
	if err != nil {
		return nil, err
	}

	//eliminar los comentarios
	err = t.CommentRemoveAll(tx, item)
	if err != nil {
		return nil, err
	}

	//eliminar los attachments
	err = t.AttachmentRemoveAll(tx, item)
	if err != nil {
		return nil, err
	}

	//eliminar la tarea
	_, err = t.Base.Executor(tx).Delete(item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

type IssueHierrachy struct {
	Id       string
	Pkey     string
	IdParent sql.NullString
	Level    int
	Path     string
}

func (t *IssueDb) FindRoot(tx interface{}, id string) (root string, err error) {
	currItem := &IssueHierrachy{}
	err = t.Base.Executor(tx).SelectOne(currItem, "select * from view_issue_parents where id=$1", id)

	if err != nil {
		return
	}

	if !currItem.IdParent.Valid {
		root = id
		return
	}
	paths := strings.Split(currItem.Path, ":")
	root = paths[0]
	return
}

func (t *IssueDb) Update(tx interface{}, item *domain.Issue) error {
	err := item.Validate()
	if err != nil {
		return err
	}

	oldItem, err := t.Load(tx, item.Id, "")
	if err != nil {
		return err
	}

	if !oldItem.AcceptUpdates() {
		return errors.New("No se puede actualizar una tarea resuelta")
	}

	//validar permisos sobre esta accion
	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_EDIT_ISSUE
	ok, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, oldItem, permission)
	if !ok {
		return errors.New("db.Base.AccessDenied()")
	}
	if oldItem.IdAssignee.String != item.IdAssignee.String {
		permission.Name = domain.PERMISSION_ASSIGN_USER
		ok, err = t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, oldItem, permission)
		if !ok {
			return errors.New("db.Base.AccessDenied()")
		}
	}
	if oldItem.IdReporter != item.IdReporter {
		permission.Name = domain.PERMISSION_CHANGE_REPORTER
		ok, err = t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, oldItem, permission)
		if !ok {
			return errors.New("db.Base.AccessDenied()")
		}
	}
	if oldItem.DueDate != item.DueDate {
		permission.Name = domain.PERMISSION_CHANGE_DUEDATE
		ok, err = t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, oldItem, permission)
		if !ok {
			item.DueDate = oldItem.DueDate
		}
	}

	item.DateCreated = oldItem.DateCreated
	item.IdStatus = oldItem.IdStatus
	item.IdProject = oldItem.IdProject
	item.IdWorkflow = oldItem.IdWorkflow
	item.Pkey = oldItem.Pkey
	item.IdParent = oldItem.IdParent

	_, err = t.Base.Executor(tx).Update(item)

	if err != nil {
		return err
	}

	//generar suscripciones automaticamente si ha cambiado el assignee o el reporter
	if item.Assignee != nil && len(item.Assignee.Id) != 0 && (oldItem.Assignee == nil || oldItem.Assignee.Id != item.Assignee.Id) {
		err = t.SubscriptionAdd(tx, item, item.Assignee)
		if err != nil {
			return err
		}
	}

	if item.Reporter.Id != oldItem.Reporter.Id {
		err = t.SubscriptionAdd(tx, item, item.Reporter)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *IssueDb) Children(tx interface{}, issue *domain.Issue) ([]database.IssueGrid, error) {
	var items []database.IssueGrid
	_, err := t.Base.Executor(tx).Select(&items, "select * from view_issues where idparent=$1 order by datecreated", issue.Id)
	if err != nil {
		return nil, err
	}
	for i := range items {
		item := &items[i]
		item.Initialize()
	}
	return items, nil
}

func (t *IssueDb) MoveProject(tx interface{}, issue *domain.Issue, target *domain.Project) error {
	if target == nil || len(target.Id) == 0 {
		return errors.New("No has seleccionado ningun proyecto valido para esta operacion")
	}

	oldIssue, err := t.Load(tx, issue.Id, "")
	if err != nil {
		return err
	}
	oldProject, err := t.ProjectDb().Load(tx, target.Id)
	if err != nil {
		return err
	}

	if issue.Project != nil && oldIssue.Project != nil && oldIssue.Project.Id == target.Id {
		return errors.New("La tarea ya pertenece a este projecto, no hay nada que hacer.")
	}

	target = oldProject
	issue = oldIssue

	//validar el permiso sobre el proyecto anterior
	permission := &domain.PermissionName{}
	permission.Name = domain.PERMISSION_DELETE_ISSUE
	ok, err := t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, issue, permission)
	if !ok {
		return errors.New("db.Base.AccessDenied()")
	}

	project, err := t.ProjectDb().NextNumber(tx, target.Id)
	if err != nil {
		return err
	}
	issue.Pkey = project.IssueKey()
	issue.Project = target
	issue.IdProject.Valid = false

	//validar el permiso sobre el proyecto nuevo
	permission.Name = domain.PERMISSION_CREATE_ISSUE
	ok, err = t.PermissionDb().AllowedUser(tx, t.Base.CurrentSession().User, issue, permission)
	if !ok {
		return errors.New("db.Base.AccessDenied()")
	}

	_, err = t.Base.Executor(tx).Update(issue)
	return err
}
