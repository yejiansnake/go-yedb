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

	All()  (*sql.Rows, error)
	One() *sql.Row
	Count() (count int64, err error)
	Max(fieldName string, ptrOutValue interface{}) error
	Min(fieldName string, ptrOutValue interface{}) error

	FillRows(ptrRows interface{}) error
	FillRow(prtRow interface{}) error
}