package yedb

import (
	"fmt"
	"time"
)

type testTable struct {
	Id int
	Name string
	Value float32
}

var dbKey = "testDB"
var dbConfig = DbConfig{Driver:"mysql", Addr:"127.0.0.1:3306", Name:"testDB", User:"user",Pwd:"pwd"}
var tableName = "t_test"

func testFillRows()  {
	DbConfigMgrInstance().Set(dbKey,
		&DbConfig{Driver:dbConfig.Driver, Addr: dbConfig.Addr, Name:dbConfig.Name, User:dbConfig.User, Pwd:dbConfig.Pwd})

	model := ModelNew(dbKey, tableName)

	inParams := make([]interface{}, 3)
	inParams[0] = 5
	inParams[1] = 3
	inParams[2] = 6

	var objArray []testTable

	err := model.Find().AndWhereIn("id", inParams...).FillRows(&objArray)

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	for index, obj := range objArray {
		fmt.Printf("row[%v]: %v %v %v \r\n", index, obj.Id, obj.Name, obj.Value)
	}
}

func testFillRow()  {
	DbConfigMgrInstance().Set(dbKey,
		&DbConfig{Driver:dbConfig.Driver, Addr: dbConfig.Addr, Name:dbConfig.Name, User:dbConfig.User, Pwd:dbConfig.Pwd})

	model := ModelNew(dbKey, tableName)

	var obj testTable

	err := model.Find().AndWhere(&DbParams{"id" : 5}).FillRow(&obj)
	//or you can code use AndWhereEx , it is like name param
	//err := model.Find().AndWhereEx("id=:id", &DbParams{":id" : 5}).FillRow(&obj)

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("row: %v %v %v \r\n",obj.Id, obj.Name, obj.Value)
}

func testInsert()  {
	DbConfigMgrInstance().Set(dbKey,
		&DbConfig{Driver:dbConfig.Driver, Addr: dbConfig.Addr, Name:dbConfig.Name, User:dbConfig.User, Pwd:dbConfig.Pwd})

	model := ModelNew(dbKey, tableName)

	res, err := model.Insert(
		&DbParams{"name" : fmt.Sprintf("test_%d", time.Now().Unix()),
			"value" : 12.3})

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("finish: %v", res)
}

func testUpdate()  {
	DbConfigMgrInstance().Set(dbKey,
		&DbConfig{Driver:dbConfig.Driver, Addr: dbConfig.Addr, Name:dbConfig.Name, User:dbConfig.User, Pwd:dbConfig.Pwd})

	model := ModelNew(dbKey, tableName)

	res, err := model.Update(
		&DbParams{"name" : fmt.Sprintf("test_%d", time.Time{}.Unix())},
		&DbParams{"id" : 1})

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("finish: %v", res)
}

func testUpdateCounters()  {
	DbConfigMgrInstance().Set(dbKey,
		&DbConfig{Driver:dbConfig.Driver, Addr: dbConfig.Addr, Name:dbConfig.Name, User:dbConfig.User, Pwd:dbConfig.Pwd})

	model := ModelNew(dbKey, tableName)

	res, err := model.UpdateCounters(
		&DbParams{"value" : 1},
		&DbParams{"id" : 2})

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("finish: %v", res)
}

func testDelete()  {
	DbConfigMgrInstance().Set(dbKey,
		&DbConfig{Driver:dbConfig.Driver, Addr: dbConfig.Addr, Name:dbConfig.Name, User:dbConfig.User, Pwd:dbConfig.Pwd})

	model := ModelNew(dbKey, tableName)

	res, err := model.Delete(&DbParams{"cid" : 1})

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("finish: %v", res)
}