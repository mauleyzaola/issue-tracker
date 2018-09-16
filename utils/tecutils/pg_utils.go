package tecutils

import (
	"github.com/mauleyzaola/gorp"
)

type DbTableInfo struct {
	TableName string `json:"tableName"`
	Bytes     int64  `json:"bytes"`
}

func TablesInfo(tx gorp.SqlExecutor) (items []DbTableInfo, err error) {
	query := `
	select table_name as tablename,pg_total_relation_size(table_name)  as bytes
	from information_schema.tables 
	where table_schema='public' 
	and table_type='BASE TABLE' 
	order by 1; 
	`
	_, err = tx.Select(&items, query)
	return
}

func DbSize(tx gorp.SqlExecutor, dbName string) (string, error) {
	query := `
		select t.dbsize
		from
		(select t1.datname AS dbname, 
		       pg_size_pretty(pg_database_size(t1.datname)) as dbsize 
		from pg_database t1) as t 
		where t.dbname=$1;
	`
	return tx.SelectStr(query, dbName)
}
