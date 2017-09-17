# go-yedb 
GO 轻量级db访问库，简单易用，较少的规范，可以减少对特定数据库SQL的混合编码，减轻原生数据库访问接口的编码工作量

### 使用的数据库驱动
Mysql : [github.com/Go-SQL-Driver/MySQL](https://github.com/Go-SQL-Driver/MySQL)

### 接口说明
DbConfigMgr strcut : 数据库配置管理，提供全局唯一实例，一次 Set 后即可通过 key 获取数据库对象
DbModel struct : 由 func ModelNew 方法创建，ModelNew 直接使用 DbConfigMgr 存储的配置生成对象，省去每次都重新获取数据库对象，结构方法实现数据表的增删改，并可以生成查寻构造器 IQuery
IQuery interface : 由 DbModel.Find() 生成，支持多种形式 where 构造，统一输出数据。

### 例子

示例数据表
```sql
CREATE TABLE `t_test` (
	`id` BIGINT(20) NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(50) NULL DEFAULT '0' COLLATE 'utf8mb4_unicode_ci',
	`value` FLOAT NULL DEFAULT '0',
	PRIMARY KEY (`id`)
)
COLLATE='utf8mb4_unicode_ci'
ENGINE=InnoDB;
```

对应数据结构（数据成员必须首字母大写，Id 对应 id）
```go
type testTable struct {
	Id int
	Name string
	Value float32
}
```

先设置配置管理对象并创建 DbModel
```go
	yedb.DbConfigMgrInstance().Set(dbKey,
		&yedb.DbConfig{Driver:dbConfig.Driver, Addr: dbConfig.Addr, Name:dbConfig.Name, User:dbConfig.User, Pwd:dbConfig.Pwd})
		
	model := ModelNew(dbKey, tableName)
```

获取多行数据
```go
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
```

获取单行数据
```go
	var obj testTable

	err := model.Find().AndWhere(&DbParams{"id" : 5}).FillRow(&obj)
	//or you can code use AndWhereEx , it is like name param
	//err := model.Find().AndWhereEx("id=:id", &DbParams{":id" : 5}).FillRow(&obj)

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("row: %v %v %v \r\n",obj.Id, obj.Name, obj.Value)
```

插入一行数据
```go
	lastID, err := model.Insert(
		&DbParams{"name" : fmt.Sprintf("test_%d", time.Now().Unix()),
			"value" : 12.3})

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("lastID: %v", lastID)
```

更新数据
```go
	affectCount, err := model.Update(
		&DbParams{"name" : fmt.Sprintf("test_%d", time.Time{}.Unix())},
		&DbParams{"id" : 1})

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("affectCount: %v", affectCount)
```

更新累加数据
```go
	affectCount, err := model.UpdateCounters(
		&DbParams{"value" : 1},
		&DbParams{"id" : 2})

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("affectCount: %v", affectCount)
```

删除数据
```go
	res, err := model.Delete(&DbParams{"cid" : 1})

	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("finish: %v", res)
```
