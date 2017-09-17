package yedb

//各个数据库的特殊处理逻辑

type ISqlWrapper interface {
	warpFiled(name string) string
	warpLimitOffset(limit int64, offset int64) string
}