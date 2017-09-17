package yedb

import (
	"database/sql"
	"fmt"
	"strings"
	"errors"
)

const (
	DB_DRIVER_MYSQL = "mysql"
)

type DbModel struct {
	config    *DbConfig
	db        *sql.DB
	wrapper   ISqlWrapper
	tableName string
}

func ModelNew(key string, tableName string) *DbModel  {
	if tableName == "" {
		return nil
	}

	config, db := DbConfigMgrInstance().Get(key)

	if config == nil {
		return nil
	}

	if config.Driver == DB_DRIVER_MYSQL {
		return &DbModel{config:config, db:db, wrapper:&MysqlWrapper{}, tableName:tableName}
	}

	return nil
}

func (ptr *DbModel) Find() IQuery {

	if ptr.config.Driver == DB_DRIVER_MYSQL {
		return mysqlQueryNew(ptr)
	}

	return nil
}

func (ptr *DbModel) Insert(params *DbParams) (int64, error) {
	if params == nil {
		return 0, errors.New("params invalid")
	}

	count := len(*params)

	if count == 0 {
		return 0, errors.New("params invalid")
	}

	var fields []string = make([]string, count)
	var fieldParams []string = make([]string, count)
	args := make([]interface{}, count)

	index := 0
	for name, value := range *params {
		fields[index] = ptr.wrapper.warpFiled(name)
		fieldParams[index] = "?"
		args[index] = value
		index++
	}

	strSql := fmt.Sprintf(
		"INSERT INTO %v(%v) VALUES(%v)",
		ptr.wrapper.warpFiled(ptr.tableName),
		strings.Join(fields, ","),
		strings.Join(fieldParams, ","))

	res, err := ptr.db.Exec(strSql, args...)

	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (ptr *DbModel) Update(setParams *DbParams, conditionParams *DbParams) (int64, error) {
	if setParams == nil {
		return 0, errors.New("params invalid")
	}

	count := len(*setParams)

	if count == 0 {
		return 0, errors.New("params invalid")
	}

	setSql, setArgs, err := ptr.buildParams(setParams)

	if err != nil {
		return 0, err
	}

	whereSql := ""

	var whereArgs []interface{}

	if conditionParams != nil && 0 != len(*conditionParams) {
		whereSql, whereArgs, err = ptr.buildParams(conditionParams)

		if err != nil {
			return 0, err
		}
	}

	return ptr.update(&setSql, setArgs, &whereSql, whereArgs)
}

func (ptr *DbModel) UpdateCounters(setParams *DbParams, conditionParams *DbParams) (int64, error) {
	if setParams == nil {
		return 0, errors.New("params invalid")
	}

	count := len(*setParams)

	if count == 0 {
		return 0, errors.New("params invalid")
	}

	setSql, setArgs, err := ptr.buildCounterParams(setParams)

	if err != nil {
		return 0, err
	}

	whereSql := ""

	var whereArgs []interface{}

	if conditionParams != nil && 0 != len(*conditionParams) {
		whereSql, whereArgs, err = ptr.buildParams(conditionParams)

		if err != nil {
			return 0, err
		}
	}

	return ptr.update(&setSql, setArgs, &whereSql, whereArgs)
}

func (ptr *DbModel) update(setSql *string,
	setArgs []interface{},
	whereSql *string,
	whereArgs []interface{}) (int64, error) {

	strSql := ""

	if nil == whereSql || "" == *whereSql {
		strSql = fmt.Sprintf("UPDATE %v SET %v", ptr.tableName, *setSql)
	} else {
		strSql = fmt.Sprintf(
			"UPDATE %v SET %v WHERE %v",
			ptr.wrapper.warpFiled(ptr.tableName),
			*setSql,
			*whereSql)
	}

	args := arrayMerge(setArgs, whereArgs)

	res, err := ptr.db.Exec(strSql, args...)

	if err != nil {
		return 0, err
	}

	affectCount, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}

	return affectCount, nil
}

func (ptr *DbModel) Delete(params *DbParams) (sql.Result, error) {
	whereSql, args, err := ptr.buildParams(params)

	if err != nil {
		return nil, err
	}

	strSql := fmt.Sprintf(
		"DELETE FROM %v WHERE %v",
		ptr.wrapper.warpFiled(ptr.tableName),
		whereSql)

	return ptr.db.Exec(strSql, args...)
}

func (ptr *DbModel) buildParams(params *DbParams) (whereSql string, args []interface{}, err error)  {
	if params == nil {
		return "", nil, errors.New("params invalid")
	}

	count := len(*params)

	if count == 0 {
		return "", nil, errors.New("params invalid")
	}

	args = make([]interface{}, count)
	var fields []string = make([]string, count)

	index := 0
	for name, value := range *params {
		fields[index] = fmt.Sprintf("%v=?", ptr.wrapper.warpFiled(name))
		args[index] = value
		index++
	}

	whereSql = strings.Join(fields, " AND ")

	return
}

func (ptr *DbModel) buildCounterParams(params *DbParams) (whereSql string, args []interface{}, err error)  {
	if params == nil {
		return "", nil, errors.New("params invalid")
	}

	count := len(*params)

	if count == 0 {
		return "", nil, errors.New("params invalid")
	}

	args = make([]interface{}, count)
	var fields []string = make([]string, count)

	index := 0
	for name, value := range *params {
		nameEx := ptr.wrapper.warpFiled(name)
		fields[index] = fmt.Sprintf("%v=%v+?", nameEx, nameEx)
		args[index] = value
		index++
	}

	whereSql = strings.Join(fields, " AND ")

	return
}