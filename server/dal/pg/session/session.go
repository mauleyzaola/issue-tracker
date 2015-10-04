package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
)

type SessionDb struct {
	Base   *pg.Db
	userDb *database.User
}

func New(db database.Db) *SessionDb {
	base := db.(*pg.Db)
	return &SessionDb{Base: base}
}

func (t *SessionDb) UserDb() database.User {
	return *t.userDb
}

func (t *SessionDb) SetUserDb(item *database.User) {
	t.userDb = item
}

func (t *SessionDb) ValidateToken(tx interface{}, id string) (*domain.Session, error) {
	session, err := t.Load(tx, id)
	if err != nil {
		return nil, err
	}
	if !session.User.IsActive {
		return nil, fmt.Errorf("invalid user")
	}
	if session.Expires.Before(time.Now()) {
		return nil, errors.New("session has expired")
	}
	return session, err
}

func (t *SessionDb) List(tx interface{}, user *domain.User) ([]domain.Session, error) {
	var items []domain.Session
	_, err := t.Base.Executor(tx).Select(&items, "select * from sessions where iduser=$1 and id <> $2", user.Id, t.Base.CurrentSession().Id)
	if err != nil {
		return nil, err
	}
	for i := range items {
		item := &items[i]
		item.Initialize()
	}
	return items, err
}

func (t *SessionDb) Load(tx interface{}, id string) (*domain.Session, error) {
	item := &domain.Session{}

	err := t.Base.Executor(tx).SelectOne(item, "select * from sessions where id=$1", id)
	if err != nil {
		return nil, err
	}
	item.User, err = t.UserDb().Load(tx, item.IdUser)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (t *SessionDb) Remove(tx interface{}, id string, validateCurrentSession bool) error {
	if validateCurrentSession && id == t.Base.CurrentSession().Id {
		return fmt.Errorf("cannot revome the current session")
	}

	session, err := t.Load(tx, id)
	if err != nil {
		return err
	}

	_, err = t.Base.Executor(tx).Delete(session)
	return err
}

func (t *SessionDb) Update(tx interface{}, session *domain.Session) error {
	oldItem, err := t.Load(tx, session.Id)
	if err != nil {
		return err
	}
	session.Id = oldItem.Id
	session.DateCreated = oldItem.DateCreated
	_, err = t.Base.Executor(tx).Update(session)
	return err
}

func (t *SessionDb) Create(tx interface{}, session *domain.Session) error {
	session.Initialize()
	return t.Base.Executor(tx).Insert(session)
}
