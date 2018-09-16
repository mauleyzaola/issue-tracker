package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mauleyzaola/issue-tracker/server/domain"
	"github.com/mauleyzaola/issue-tracker/server/operations"
	"github.com/mauleyzaola/issue-tracker/server/operations/database"
	"github.com/mauleyzaola/issue-tracker/utils/tecweb/setup"
	"github.com/zenazn/goji/web"
)

type ApiBase struct {
	Database *database.DbOperations
	Setup    *setup.Setup
}

type HttpMessage struct {
	StatusCode int
	Message    string
}

func New(db *database.DbOperations, setup *setup.Setup) *ApiBase {
	return &ApiBase{Database: db, Setup: setup}
}

func (base *ApiBase) SuccessResponse(tx interface{}, w http.ResponseWriter) {
	base.Response(tx, nil, w, nil)
}

func (base *ApiBase) ErrorResponse(tx interface{}, err error, w http.ResponseWriter) {
	base.Response(tx, err, w, nil)
}

func (base *ApiBase) Response(tx interface{}, err error, w http.ResponseWriter, data interface{}) {
	var result []byte
	msg := &HttpMessage{StatusCode: http.StatusOK}
	if err != nil {
		if tx != nil {
			base.Database.Db.Rollback(tx)
		}

		data = nil
		msg.StatusCode = http.StatusBadRequest
		msg.Message = err.Error()
	} else {
		if tx != nil {
			base.Database.Db.Commit(tx)
		}
	}
	w.WriteHeader(msg.StatusCode)
	w.Header().Set(operations.HEADER_STATUS, strconv.Itoa(msg.StatusCode))
	if data != nil {
		result, _ = json.Marshal(data)
	} else {
		result, _ = json.Marshal(msg)
	}
	w.Write(result)
}

func (base *ApiBase) Decode(v interface{}, r *http.Request) error {
	return json.NewDecoder(r.Body).Decode(v)
}

//gets the parameter value from the query string or from the url
func (base *ApiBase) ParamValue(name string, c web.C, r *http.Request) string {
	value := c.URLParams[name]
	if len(value) != 0 {
		return value
	} else {
		value = r.URL.Query().Get(name)
	}

	return value
}

func (base *ApiBase) CurrentSession(c web.C) *domain.Session {
	if val, ok := c.Env[operations.MIDDLEWARE_CURRENT_SESSION]; ok {
		return val.(*domain.Session)
	}
	return nil
}

func (base *ApiBase) IpAddress(r *http.Request) string {
	return r.RemoteAddr
}
