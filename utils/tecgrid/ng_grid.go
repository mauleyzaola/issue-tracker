// This package is a helper for parsing http request from AngularJs application
//to GO api's I have developed
package tecgrid

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mauleyzaola/glog"
	"github.com/mauleyzaola/gorp"
)

//The main structure for parsing objects from json client
type NgGrid struct {
	PageNumber     int64       `json:"pageNumber"`
	PageSize       int64       `json:"pageSize"`
	Query          string      `json:"query"`
	SortDirection  string      `json:"sortDirection"`
	SortField      string      `json:"sortField"`
	TotalCount     int64       `json:"totalCount"`
	FilterCount    int64       `json:"filterCount"`
	Rows           interface{} `json:"rows"`
	GeneratedQuery string      `json:"-"`
	MainQuery      string      `json:"-"`
}

const (
	QS_SORT_DIRECTION = "sortDirection"
	QS_SORT_FIELD     = "sortField"
	QS_PAGE_NUMBER    = "pageNumber"
	QS_PAGE_SIZE      = "pageSize"
	QS_QUERY          = "query"
)

//Parses a query string from url into NgGrid
//Useful when you are dealing with GET methods instead of POST and decoding
func ParseQueryString(values map[string][]string) *NgGrid {
	g := &NgGrid{}
	if value := values[QS_SORT_DIRECTION]; len(value) != 0 {
		g.SortDirection = value[0]
	}
	if value := values[QS_SORT_FIELD]; len(value) != 0 {
		g.SortField = value[0]
	}
	if value := values[QS_QUERY]; len(value) != 0 {
		g.Query = value[0]
	}
	if value := values[QS_PAGE_NUMBER]; len(value) != 0 {
		pn, err := strconv.Atoi(value[0])
		if err == nil {
			g.PageNumber = int64(pn)
		}
	}
	if value := values[QS_PAGE_SIZE]; len(value) != 0 {
		pn, err := strconv.Atoi(value[0])
		if err == nil {
			g.PageSize = int64(pn)
		}
	}

	if g.PageNumber < 1 {
		g.PageNumber = 1
	}

	return g
}

//Returns the actual query to lowercase
func (g *NgGrid) GetQuery() string {
	return strings.ToLower(g.Query)
}

//Calculates the page number
func (g *NgGrid) GetPageNumber() int64 {
	if g.PageNumber == 0 {
		g.PageNumber = 1
	}
	return g.PageNumber
}

//Calculates the page size
func (g *NgGrid) GetPageSize() int64 {
	if g.PageSize <= 0 {
		g.PageSize = 10
	}
	return g.PageSize
}

//Helper for getting the number of rows I need to skip based on
//pageNumber and pageSize values
func (g *NgGrid) Skip() int64 {
	return (g.GetPageNumber() - 1) * g.GetPageSize()
}

//The same as pageSize
func (g *NgGrid) Limit() int64 {
	return g.GetPageSize()
}

func (g *NgGrid) FromRow() int64 {
	return g.ToRow() - g.PageSize + 1
}

func (g *NgGrid) ToRow() int64 {
	return g.GetPageNumber() * g.GetPageSize()
}

//Generates the sql and calculates the number of rows from the result query
func (g *NgGrid) GenerateSql(tx gorp.SqlExecutor, fields []string) error {

	var sq bytes.Buffer
	sq.WriteString("select count(*) ")
	sq.WriteString("from ( %s ) as t")
	query := fmt.Sprintf(sq.String(), g.MainQuery)

	totalCount, err := tx.SelectInt(query)
	g.TotalCount = totalCount
	g.FilterCount = g.TotalCount

	if g.TotalCount == 0 {
		return err
	}

	sq.Reset()
	if fields == nil {
		sq.WriteString("select * ")
	} else {
		sq.WriteString(fmt.Sprintf("select %s ", strings.Join(fields, ",")))
	}

	sq.WriteString("from ")
	sq.WriteString("(select row_number() over(order by %s %s) as rownum, ")
	sq.WriteString("t.* ")
	sq.WriteString("from ")
	sq.WriteString("( %s ) as t) as t ")
	sq.WriteString("where t.rownum between %d and %d ")

	if g.SortField == "" {
		return errors.New("Query sortField parameter is missing")
	}

	sortDirection := "asc"
	if strings.ToLower(g.SortDirection) == "desc" {
		sortDirection = "desc"
	}
	g.SortDirection = sortDirection

	//if the sort field refers to an inner object e.g. bomItem.line.name, then we should consider only the latest part
	sf := strings.Split(g.SortField, ".")
	g.SortField = sf[len(sf)-1]

	g.GeneratedQuery = fmt.Sprintf(sq.String(), strings.ToLower(g.SortField), g.SortDirection, g.MainQuery, g.FromRow(), g.ToRow())

	glog.V(4).Infoln(g.GeneratedQuery)

	return err
}

//Executes sql statements along with its parameters
func (g *NgGrid) ExecuteSqlParameters(tx gorp.SqlExecutor, target interface{}, fields []string, pars []interface{}) (err error) {
	if len(g.MainQuery) == 0 {
		return fmt.Errorf("missing main query")
	}

	toRow := g.GetPageNumber() * g.GetPageSize()
	fromRow := toRow - g.PageSize + 1

	var sq bytes.Buffer
	sq.WriteString("select count(*) ")
	sq.WriteString("from ( %s ) as t")
	query := fmt.Sprintf(sq.String(), g.MainQuery)

	var totalCount int64

	if pars != nil && len(pars) != 0 {
		totalCount, err = tx.SelectInt(query, pars...)
	} else {
		totalCount, err = tx.SelectInt(query)
	}
	if err != nil {
		return
	}
	if totalCount == 0 {
		return
	}

	g.TotalCount = totalCount
	g.FilterCount = g.TotalCount

	sq.Reset()
	if fields == nil || len(fields) == 0 {
		if len(g.SortField) == 0 {
			err = errors.New("Query sortField parameter is missing")
			return
		}
		sq.WriteString("select * ")
	} else {
		if len(g.SortField) == 0 {
			g.SortField = fields[0]
		}
		sq.WriteString(fmt.Sprintf("select %s ", strings.Join(fields, ",")))
	}

	sq.WriteString("from ")
	sq.WriteString("(select row_number() over(order by %s %s nulls last) as rownum, ")
	sq.WriteString("t.* ")
	sq.WriteString("from ")
	sq.WriteString("( %s ) as t) as t ")
	sq.WriteString("where t.rownum between %d and %d ")

	sortDirection := "asc"
	if strings.ToLower(g.SortDirection) == "desc" {
		sortDirection = "desc"
	}
	g.SortDirection = sortDirection

	//if the sort field refers to an inner object e.g. bomItem.line.name, then we should consider only the latest part
	sf := strings.Split(g.SortField, ".")
	g.SortField = sf[len(sf)-1]

	g.GeneratedQuery = fmt.Sprintf(sq.String(), strings.ToLower(g.SortField), g.SortDirection, g.MainQuery, fromRow, toRow)

	glog.V(4).Infoln(g.GeneratedQuery)

	if pars != nil && len(pars) != 0 {
		_, err = tx.Select(target, g.GeneratedQuery, pars...)
	} else {
		_, err = tx.Select(target, g.GeneratedQuery)
	}

	if err != nil {
		return
	}

	g.Rows = target

	return
}
