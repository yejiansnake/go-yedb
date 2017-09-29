package yedb

import (
	"container/list"
	"database/sql"
	"strings"
	"fmt"
	"regexp"
	"reflect"
	"errors"
)

type MysqlQuery struct {
	model         *DbModel
	selectFields  []string
	whereParams   list.List
	groupByFields []string
	havingParams  list.List
	orderByParams map[string]int
	limit         int64
	offset        int64
}

func mysqlQueryNew(model *DbModel) IQuery {
	return &MysqlQuery{model: model, orderByParams: make(map[string]int)}
}

func (ptr *MysqlQuery) Select(fields ...string) IQuery{
	if fields == nil {
		return ptr
	}

	count := len(fields)

	if count == 0 {
		return ptr
	}

	ptr.selectFields = make([]string, count)

	for index,value := range fields {
		ptr.selectFields[index] = ptr.model.wrapper.warpFiled(value);
	}

	return ptr
}

func (ptr *MysqlQuery) AndWhere(params *DbParams) IQuery{

	if params == nil || 0 == len(*params) {
		return ptr
	}

	count := len(*params)

	param := dbWhereParam{}
	param.params = make([]interface{}, count)

	tempArr := make([]string, count)
	index := 0
	for key, value := range *params {
		tempArr[index] = fmt.Sprintf("%v=?", ptr.model.wrapper.warpFiled(key))
		param.params[index] = value
		index++
	}

	param.strSql = strings.Join(tempArr, " AND ")

	ptr.whereParams.PushBack(param)

	return ptr
}

func (ptr *MysqlQuery) AndWhereEx(whereSql string, params *DbParams) IQuery{
	//使用 param_name=:param_name 的方式做参数替换
	//构造符合个数的 ?
	//构造后提交

	param := dbWhereParam{}

	reg := regexp.MustCompile(`:\w+`)
	names := reg.FindAllString(whereSql, -1)

	count := len(names)

	if count != 0 {
		args := list.New()

		for _, name := range names {
			value, isOn := (*params)[name[1:]]

			if isOn {
				args.PushBack(value)
			} else {
				return ptr
			}
		}

		param.params = make([]interface{}, args.Len())
		index := 0
		for item := args.Front(); item != nil ; item = item.Next() {
			param.params[index] = item.Value
			index++
		}

		param.strSql = reg.ReplaceAllString(whereSql, "?")

	} else {
		param.strSql = whereSql
	}

	ptr.whereParams.PushBack(param)

	return ptr
}

func (ptr *MysqlQuery) AndWhereIn(name string, params ...interface{}) IQuery{

	if params == nil || 0 == len(params) {
		return ptr
	}

	count := len(params)
	tempArr := make([]string, count)

	for index := 0; index < count; index++ {
		tempArr[index] = "?"
	}

	param := dbWhereParam{}
	param.strSql = fmt.Sprintf("%v IN (%v)", ptr.model.wrapper.warpFiled(name), strings.Join(tempArr, ","))
	param.params = params

	ptr.whereParams.PushBack(param)

	return ptr
}

func (ptr *MysqlQuery) AndGroupBy(params ...string) IQuery{
	if params == nil {
		return ptr
	}

	count := len(params)

	if count == 0 {
		return ptr
	}

	ptr.groupByFields = make([]string, count)

	for index, value := range params {
		ptr.groupByFields[index] = ptr.model.wrapper.warpFiled(value)
	}

	return ptr
}

//暂时不支持
func (ptr *MysqlQuery) AndHaving(params ...string) IQuery{
	if params == nil {
		return ptr
	}

	count := len(params)

	if count == 0 {
		return ptr
	}

	for _,value := range params {
		ptr.havingParams.PushBack(value)
	}

	return ptr
}

//0:升序 其他:倒序
func (ptr *MysqlQuery) AndOrderBy(params *DbSortParams) IQuery{
	if params == nil {
		return ptr
	}

	count := len(*params)

	if count == 0 {
		return ptr
	}

	for key,value := range *params {
		ptr.orderByParams[ptr.model.wrapper.warpFiled(key)] = int(value)
	}

	return ptr
}

func (ptr *MysqlQuery) Limit(value int64) IQuery{
	ptr.limit = value
	return ptr
}

func (ptr *MysqlQuery) Offset(value int64) IQuery{
	ptr.offset = value
	return ptr
}

func (ptr *MysqlQuery) RawSql()  (*string)  {
	strSql, _ := ptr.build()
	return strSql
}

func (ptr *MysqlQuery) All()  (*sql.Rows, error)  {
	strSql, args := ptr.build()
	return ptr.model.db.Query(*strSql, args...)
}

func (ptr *MysqlQuery) One() *sql.Row {
	strSql, args := ptr.build()
	return ptr.model.db.QueryRow(*strSql, args...)
}

func (ptr *MysqlQuery) Count() int64 {
	return ptr.queryScalar("COUNT(0)")
}

func (ptr *MysqlQuery) queryScalar(selectExpression string) (count int64) {
	selectFields := ptr.selectFields
	ptr.selectFields = []string{selectExpression}
	strSql, args := ptr.build()
	ptr.selectFields = selectFields
	row := ptr.model.db.QueryRow(*strSql, args...)
	row.Scan(&count)
	return
}

func (ptr *MysqlQuery) FillRows(rowsSlicePtr interface{}) error {
	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))

	if sliceValue.Kind() != reflect.Slice {
		return errors.New("needs a pointer to a slice")
	}

	rows, err := ptr.All()

	if err != nil {
		return err
	}

	err = fillModels(rowsSlicePtr, rows)

	if err != nil {
		return err
	}

	return nil
}

func (ptr *MysqlQuery) FillRow(rowPtr interface{}) error {

	obj := reflect.ValueOf(rowPtr)

	if obj.Kind() != reflect.Ptr {
		return errors.New("needs a pointer")
	}

	limit := ptr.limit
	defer ptr.Limit(limit)
	ptr.Limit(1)

    rows, err := ptr.All()

	if err != nil {
		return err
	}

	err = fillModel(rowPtr, rows)

	if err != nil {
		return err
	}

	return nil
}

func (ptr *MysqlQuery) build() (*string, []interface{}) {

	selectStr := ptr.buildSelect()

	whereStr, args := ptr.buildWhere()

	groupByStr := ptr.buildGroup()

	havingStr := ptr.buildHaving()

	orderByStr := ptr.buildOrder()

	limitOffsetStr := ptr.model.wrapper.warpLimitOffset(ptr.limit, ptr.offset)

	//				 输出列 表名 where子句 组子句 组合子句 排序子句 LIMIT OFFSET
	strSql := fmt.Sprintf("SELECT %v FROM %v %v %v %v %v %v",
		*selectStr,
		ptr.model.tableName,
		*whereStr,
		*groupByStr,
		*havingStr,
		*orderByStr,
		limitOffsetStr)

	return &strSql, args
}

func (ptr *MysqlQuery) buildSelect() (*string) {
	retStr := "*"

	if 0 != len(ptr.selectFields) {
		retStr = strings.Join(ptr.selectFields, ",")
	}

	return &retStr
}

func (ptr *MysqlQuery) buildWhere() (*string, []interface{}) {
	retStr := ""
	var retArgs []interface{} = nil
	whereCount := ptr.whereParams.Len()

	if 0 != whereCount {
		tempArr := make([]string, whereCount)
		argsCount := 0
		index := 0

		for item := ptr.whereParams.Front(); item != nil; item = item.Next() {
			param := item.Value.(dbWhereParam)
			tempArr[index] = param.strSql
			argsCount += len(param.params)
			index++
		}

		retArgs = make([]interface{}, argsCount)

		curCount := 0
		for item := ptr.whereParams.Front(); item != nil; item = item.Next() {
			param := item.Value.(dbWhereParam)
			copy(retArgs[curCount:], param.params)
			curCount += len(param.params)
		}

		retStr = fmt.Sprintf("WHERE %v", strings.Join(tempArr, " AND "))
	}

	return &retStr, retArgs
}

func (ptr *MysqlQuery) buildGroup() (*string) {
	retStr := ""
	count := len(ptr.groupByFields)

	if 0 != count {
		retStr = fmt.Sprintf("GROUP BY %v", strings.Join(ptr.groupByFields, ","))
	}

	return &retStr
}

func (ptr *MysqlQuery) buildHaving() (*string) {
	retStr := ""
	count := ptr.havingParams.Len()

	if 0 != count {
		tempArr := make([]string, count)

		index := 0
		for item := ptr.havingParams.Front(); item != nil; item = item.Next() {
			tempArr[index] = item.Value.(string)
			index++
		}
		retStr = fmt.Sprintf("HAVING %v", strings.Join(tempArr, " AND "))
	}

	return &retStr
}

func (ptr *MysqlQuery) buildOrder() (*string) {
	retStr := ""
	count := len(ptr.orderByParams)

	if 0 != count {
		tempArr := make([]string, count)

		index := 0
		for key,value := range ptr.orderByParams {
			tempStr := ""

			if 0 == value {
				tempStr = fmt.Sprintf("%v", key)
			} else {
				tempStr = fmt.Sprintf("%v DESC", key)
			}

			tempArr[index] = tempStr
			index++
		}

		retStr = fmt.Sprintf("ORDER BY %v", strings.Join(tempArr, ","))
	}

	return &retStr
}