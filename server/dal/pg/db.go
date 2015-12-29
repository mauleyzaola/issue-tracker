package pg

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mauleyzaola/gorp"
	"github.com/mauleyzaola/issue-tracker/server/domain"
)

type Db struct {
	DbMap   *gorp.DbMap
	session *domain.Session
}

func (r *Db) CurrentSession() *domain.Session {
	if r.session == nil {
		log.Println("There is no active session to retrieve")
	}
	return r.session
}

func (r *Db) SetCurrentSession(session *domain.Session) {
	r.session = session
}

func New(dbMap *gorp.DbMap) *Db {
	base := &Db{}
	base.DbMap = dbMap
	return base
}

func (r *Db) Executor(tran interface{}) gorp.SqlExecutor {
	tx := tran.(*gorp.Transaction)
	return tx
}

func (r *Db) Commit(tran interface{}) error {
	tx := tran.(*gorp.Transaction)
	return tx.Commit()
}

func (r *Db) Rollback(tran interface{}) error {
	tx := tran.(*gorp.Transaction)
	return tx.Rollback()
}

func (r *Db) Begin() (interface{}, error) {
	tx, err := r.DbMap.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *Db) GetTransaction(tx interface{}) gorp.SqlExecutor {
	return tx.(gorp.SqlExecutor)
}

func (r *Db) SqlTraceOn() {
	r.DbMap.TraceOn("[sql]", log.New(os.Stdout, "log:", log.Lmicroseconds))
}

func (r *Db) SqlTraceOff() {
	r.DbMap.TraceOff()
}

func (r *Db) Register() {
	r.DbMap.AddTableWithName(domain.FileItem{}, "file_item").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.Group{}, "groups").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.IssueAttachment{}, "issue_attachment").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.IssueComment{}, "issue_comment").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.IssueSubscription{}, "issue_subscription").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.Issue{}, "issue").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.PermissionName{}, "permission_name").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.PermissionSchemeItem{}, "permission_scheme_item").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.PermissionScheme{}, "permission_scheme").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.Priority{}, "priority").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.ProjectRoleMember{}, "project_role_member").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.ProjectRole{}, "project_role").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.Project{}, "project").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.Role{}, "roles").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.Session{}, "sessions").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.Status{}, "status").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.UserGroup{}, "user_group").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.User{}, "users").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.UserMeta{}, "user_meta").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.WorkflowStepMember{}, "workflow_step_member").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.WorkflowStep{}, "workflow_step").SetKeys(true, "Id")
	r.DbMap.AddTableWithName(domain.Workflow{}, "workflow").SetKeys(true, "Id")
}

//Cuenta todos los registros fisicamente de una tabla
func (r *Db) CountRows(tx interface{}, tableName string) (rowCount int64, err error) {
	query := fmt.Sprintf("select count(*) from %s", tableName)
	rowCount, err = r.Executor(tx).SelectInt(query)
	return
}

func (r *Db) SequenceExists(tx interface{}, sequenceName string) bool {
	query := "select count(*) from pg_class where relkind='S' and lower(relname)=$1;"
	rowCount, _ := r.Executor(tx).SelectInt(query, strings.ToLower(sequenceName))
	return rowCount == 1
}

func (r *Db) SequenceNextValue(tx interface{}, sequenceName string) (next int64, err error) {
	next, err = r.Executor(tx).SelectInt(fmt.Sprintf("select nextval('%s')", sequenceName))
	return
}

//Devuelve el actual valor, en caso que no tenga ninguno devuelve cero
func (r *Db) SequenceLastValue(tx interface{}, sequenceName string) (curr int64, err error) {
	curr, err = r.Executor(tx).SelectInt(fmt.Sprintf("select last_value from %s", sequenceName))
	return
}

func (r *Db) SequenceRemove(tx interface{}, seqName string) (err error) {
	_, err = r.Executor(tx).Exec(fmt.Sprintf("drop sequence if exists %s", seqName))
	return
}

func (r *Db) Close() {
	r.DbMap.Db.Close()
}
