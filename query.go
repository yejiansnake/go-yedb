package yedb

import "database/sql"

type SortType int

const (
	SORT_ASC SortType = iota
	SORT_DESC
)

type DbSortParams map[string]SortType

type dbWhereParam struct {
	strSql string
	params []interface{}
}

type IQuery interface {
	Select(fields ...string) IQuery
	AndWhere(params *DbParams) IQuery
	AndWhereEx(whereSql string, params *DbParams) IQuery
	AndWhereIn(name string, params ...interface{}) IQuery
	AndGroupBy(params ...string) IQuery
	AndHaving(params ...string) IQuery
	AndOrderBy(params *DbSortParams) IQuery
	Limit(value int64) IQuery
	Offset(value int64) IQuery

	RawSql() (*string)

	Query()  (*sql.Rows, error)
	QueryOne() *sql.Row

	FillRows(rowsPtr interface{}) error
	FillRow(rowPtr interface{}) error
}