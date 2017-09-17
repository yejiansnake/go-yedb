package yedb

import "fmt"

//各个数据库的特殊处理逻辑

type MysqlWrapper struct {

}

func (ptr *MysqlWrapper) warpFiled(name string) string {
	return fmt.Sprintf("`%s`", name);
}

func (ptr *MysqlWrapper) warpLimitOffset(limit int64, offset int64) string {

	limitStr := ""
	if limit > 0 {
		limitStr = fmt.Sprintf("LIMIT %v", limit)
	}

	offsetStr := ""
	if offset > 0 {
		offsetStr = fmt.Sprintf(" OFFSET %v", offset)
	}

	return fmt.Sprintf("%v%v", limitStr, offsetStr);
}
