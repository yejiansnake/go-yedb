package yedb

import (
	"../common"
	"sync"
	"database/sql"
	"fmt"
)

type DbConfig struct {
	Driver string	//驱动名称（mysql）
	Addr string
	Name string
	User string
	Pwd string
}

type DbConfigMgr struct {
	configs map[string]*DbConfig
	dbs     map[string]*sql.DB
	safe    bool
	lock    sync.Mutex
}

var _instance DbConfigMgr = DbConfigMgr{
	configs : make(map[string]*DbConfig),
	dbs:make(map[string]*sql.DB),
	safe:false}

func DbConfigMgrInstance() *DbConfigMgr  {
	return &_instance
}

func (ptr *DbConfigMgr) Get(key string) (*DbConfig, *sql.DB){
	if ptr.safe {
		ptr.lock.Lock()
		defer ptr.lock.Unlock()
	}

	_config, isIn := ptr.configs[key]

	if !isIn {
		return nil, nil
	}

	return _config, ptr.dbs[key]
}

func (ptr *DbConfigMgr) Set(key string, config *DbConfig) error {
	if config == nil {
		return common.BaseErrorNew("configs prt is nil")
	}

	if config.Driver == "" || config.Addr == "" || config.Name == "" || config.User == "" || config.Pwd  == "" {
		return common.BaseErrorNew("params invalid")
	}

	//user@unix(/path/to/socket)/dbname?charset=utf8
	//user:password@tcp(localhost:5555)/dbname?charset=utf8
	//user:password@/dbname
	//user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", config.User, config.Pwd, config.Addr, config.Name)

	db, err := sql.Open(config.Driver, dataSourceName)

	if err != nil {
		return err
	}

	if ptr.safe {
		ptr.lock.Lock()
		defer ptr.lock.Unlock()
	}

	curDB, isIn := ptr.dbs[key]

	if isIn {
		curDB.Close()
	}

	ptr.configs[key] = config
	ptr.dbs[key] = db

	return nil
}

func (ptr *DbConfigMgr) Remove(key string) {
	if ptr.safe {
		ptr.lock.Lock()
		defer ptr.lock.Unlock()
	}

	db, isIn := ptr.dbs[key]

	if !isIn {
		return
	}

	db.Close()

	delete(ptr.configs, key)

	delete(ptr.dbs, key)
}


