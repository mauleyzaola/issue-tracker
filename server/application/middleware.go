package application

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mauleyzaola/issue-tracker/server/api"
	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/zenazn/goji/web"
)

func (a *Application) MiddlewareAttach(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		baseApi := api.New(a.Db, a.Setup)
		c.Env[operations.MIDDLEWARE_BASE_API] = baseApi
		c.Env[operations.MIDDLEWARE_DB] = a.Db
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (a *Application) MiddlewareJsonResponseHeaders(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (a *Application) MiddlewareAuth(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tokenLabel := operations.MIDDLEWARE_TOKEN_LABEL
		var token = r.Header.Get(tokenLabel)

		if len(token) == 0 {
			token = r.URL.Query().Get(tokenLabel)
		}

		db := a.Db

		tx, _ := db.Db.Begin()
		session, err := db.SessionDb.ValidateToken(tx, token)

		if err != nil {
			db.Db.Rollback(tx)
			msg := &struct {
				StatusCode int
				Message    string
			}{
				http.StatusUnauthorized,
				"Token Invalido en la peticion HTTP",
			}

			w.WriteHeader(msg.StatusCode)
			result, _ := json.Marshal(msg)
			w.Write(result)
		} else {
			c.Env[operations.MIDDLEWARE_CURRENT_SESSION] = session
			db.Db.SetCurrentSession(session)
			db.Db.Commit(tx)
			h.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}

func (a *Application) MiddlewareSysAdmin(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error
		if c.Env[operations.MIDDLEWARE_CURRENT_SESSION] != nil {
			session := c.Env[operations.MIDDLEWARE_CURRENT_SESSION].(*domain.Session)
			if session == nil || !session.User.IsSystemAdministrator {
				err = errors.New("Acceso restringido a administradores del sistema")
			}
		} else {
			err = errors.New("Acceso restringido a administradores del sistema")
		}

		if err != nil {
			var result []byte
			msg := &struct {
				StatusCode int
				Message    string
			}{
				http.StatusBadRequest,
				err.Error(),
			}
			w.WriteHeader(msg.StatusCode)
			result, _ = json.Marshal(msg)
			w.Write(result)

		} else {
			h.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
