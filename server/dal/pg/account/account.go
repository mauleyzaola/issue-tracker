package account

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mauleyzaola/issue-tracker/server/dal/pg"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/tecutils"
)

type AccountDb struct {
	Base      *pg.Db
	userDb    *database.User
	sessionDb *database.Session
}

func New(db database.Db) *AccountDb {
	base := db.(*pg.Db)
	return &AccountDb{Base: base}
}

func (t *AccountDb) SessionDb() database.Session {
	return *t.sessionDb
}

func (t *AccountDb) SetSessionDb(item *database.Session) {
	t.sessionDb = item
}

func (t *AccountDb) UserDb() database.User {
	return *t.userDb
}

func (t *AccountDb) SetUserDb(item *database.User) {
	t.userDb = item
}

func (t *AccountDb) FullLoadUser(tx interface{}, user *domain.User) (err error) {
	if len(user.Id) != 0 {
		err = t.Base.Executor(tx).SelectOne(&user, "select * from users where id=$1", user.Id)
	} else if len(user.TokenEmail) != 0 {
		err = t.Base.Executor(tx).SelectOne(&user, "select * from users where tokenemail=$1", user.TokenEmail)
	} else if len(user.Email) != 0 {
		err = t.Base.Executor(tx).SelectOne(&user, "select * from users where email=$1", user.Email)
	} else {
		err = errors.New("No hay suficiente informacion para cargar el usuario")
	}
	return
}

//TODO: implement session.SessionCreate instead of doing this directly to sessions table
func (t *AccountDb) LoginUser(tx interface{}, email, password string, ipaddress string, validatePassword bool) (*domain.Session, error) {
	item := &domain.User{Email: email}

	if validatePassword && len(password) == 0 {
		return nil, errors.New("Password invalido")
	}

	err := t.Base.Executor(tx).SelectOne(item, "select * from users where email=$1", email)
	if err != nil {
		return nil, errors.New("Email o Password invalidos")
	}

	if len(item.Id) == 0 {
		return nil, errors.New("Email invalido")
	}

	if tecutils.Encrypt(password) != item.Password && validatePassword {
		return nil, errors.New("Email o password incorrectos")
	}

	if !item.IsActive {
		return nil, errors.New("El usuario es invalido para esta operacion")
	}

	item, err = t.UserDb().Load(tx, item.Id)
	if err != nil {
		return nil, err
	}

	//actualizar la informacion de login
	item.LoginCount++
	if item.LastLogin == nil {
		item.LastLogin = &time.Time{}
	}
	*item.LastLogin = time.Now()
	_, err = t.Base.Executor(tx).Update(item)
	if err != nil {
		return nil, err
	}

	//eliminar sesiones que hayan expirado para el usuario conectado
	_, err = t.Base.Executor(tx).Exec("delete from sessions where iduser=$1 and expires < now()", item.Id)
	if err != nil {
		return nil, err
	}

	session := &domain.Session{User: item, IpAddress: ipaddress}
	err = session.Validate()
	if err != nil {
		return nil, err
	}
	session.User = item
	err = t.SessionDb().Create(tx, session)

	if err != nil {
		return nil, err
	}
	return session, nil
}

func (t *AccountDb) LogoutUser(tx interface{}) error {
	return t.SessionDb().Remove(tx, t.Base.CurrentSession().Id, false)
}

func (t *AccountDb) ChangePassword(tx interface{}, id string, password string) error {
	user := &domain.User{Id: id}
	err := t.FullLoadUser(tx, user)
	if err != nil {
		return err
	}
	user.Password = tecutils.Encrypt(password)
	_, err = t.Base.Executor(tx).Update(user)
	return err
}

func (t *AccountDb) ChangeMyPassword(tx interface{}, password string) error {
	return t.ChangePassword(tx, t.Base.CurrentSession().User.Id, password)
}

func (t *AccountDb) LoginDelay(tx interface{}, delay time.Duration) error {
	session, err := t.SessionDb().Load(tx, t.Base.CurrentSession().Id)
	if err != nil {
		return err
	}
	session.Expires = time.Now().Add(delay)
	_, err = t.Base.Executor(tx).Update(session)
	return err
}

func (t *AccountDb) PasswordRecoverCreateToken(tx interface{}, user *domain.User) error {
	err := t.FullLoadUser(tx, user)
	if err == sql.ErrNoRows {
		return errors.New("No se encuentra el usuario solicitado")
	}
	user.TokenEmail = tecutils.UUID()
	_, err = t.Base.Executor(tx).Update(user)
	return err
}

func (t *AccountDb) VerifyTokenEmail(tx interface{}, user *domain.User, ipaddress string) error {
	requestedToken := user.TokenEmail

	if len(requestedToken) == 0 {
		return errors.New("El token proporcionado es invalido")
	}

	err := t.FullLoadUser(tx, user)
	if err != nil {
		err = errors.New("El token proporcionado es invalido")
		return err
	}

	if requestedToken != user.TokenEmail {
		err = errors.New("El token proporcionado no coincide con la base de datos")
		return err
	}

	//resetear el password y borrar el tokenemail
	user.Password = ""
	user.TokenEmail = ""
	_, err = t.Base.Executor(tx).Update(user)
	if err != nil {
		return err
	}

	//simular un login para poder asignarle un token al usuario y evitar que el password, aunque sea temporal, se muestre en el UI
	session, err := t.LoginUser(tx, user.Email, "", ipaddress, false)
	if err != nil {
		return err
	}

	user = session.User
	user.TokenEmail = session.Id
	return nil

}
