package yedb

import (
	"reflect"
	"strings"
	"errors"
	"database/sql"
	"fmt"
)

type DbParams map[string]interface{}

func arrayMerge(args ...[]interface{}) []interface{} {
	count := len(args)
	if count == 0 {
		return nil
	}

	totalCount := 0
	for _, value := range args{
		totalCount += len(value)
	}

	res := make([]interface{}, totalCount)

	mergeCount := 0
	for _, value := range args{
		count = len(value)
		copy(res[mergeCount:], value)
		mergeCount += count
	}

	return res
}

func stringArrayContains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}
	return false
}

func getTypeName(obj interface{}) string {
	typ := reflect.TypeOf(obj)
	typeStr := typ.String()

	lastDotIndex := strings.LastIndex(typeStr, ".")

	if lastDotIndex != -1 {
		typeStr = typeStr[lastDotIndex+1:]
	}

	return typeStr
}

func fillModel(rowPtr interface{}, rows *sql.Rows) error {
	obj := reflect.Indirect(reflect.ValueOf(rowPtr))

	if obj.Kind() != reflect.Struct {
		return errors.New("needs a struct")
	}

	if !obj.CanSet() {
		return errors.New("can not set")
	}

	colNames, err := rows.Columns()

	if err != nil {
		return err
	}

	nameCount := len(colNames)

	if nameCount == 0 {
		return nil
	}

	values := make([]interface{}, nameCount)
	for index := 0; index < nameCount; index++ {
		var temp interface{}
		values[index] = &temp
	}

	if rows.Next() {
		rows.Scan(values...)
	} else {
		return nil
	}

	structNames := make([]string, nameCount)

	for index, name := range colNames {
		structNames[index] = fmt.Sprintf("%s%s", strings.ToUpper(name[0:1]), name[1:])
	}

	for index, name := range structNames {
		value := obj.FieldByName(name)
		if value.IsValid() {
			temp := *(values[index].(*interface{}))
			err := ConvertValue(&value, temp)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func fillModels(rowsSlicePtr interface{}, rows *sql.Rows) error {
	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))

	if sliceValue.Kind() != reflect.Slice {
		return errors.New("needs a slice")
	}

	if !sliceValue.CanSet() {
		return errors.New("can not set")
	}

	colNames, err := rows.Columns()

	if err != nil {
		return err
	}

	nameCount := len(colNames)

	if nameCount == 0 {
		return nil
	}

	structNames := make([]string, nameCount)

	for index, name := range colNames {
		structNames[index] = fmt.Sprintf("%s%s", strings.ToUpper(name[0:1]), name[1:])
	}

	sliceElementType := sliceValue.Type().Elem()

	for rows.Next() {
		values := make([]interface{}, nameCount)
		for index := 0; index < nameCount; index++ {
			var temp interface{}
			values[index] = &temp
		}

		rows.Scan(values...)

		newValue := reflect.New(sliceElementType)
		objElem := newValue.Elem()

		for index, name := range structNames {
			value := objElem.FieldByName(name)
			if value.IsValid() {
				temp := *(values[index].(*interface{}))
				err := ConvertValue(&value, temp)
				if err != nil {
					return err
				}
			}
		}

		sliceAppend := reflect.Append(sliceValue, reflect.Indirect(reflect.ValueOf(newValue.Interface())))
		sliceValue.Set(sliceAppend)
	}

	return nil
}