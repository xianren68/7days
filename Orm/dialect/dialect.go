package dialect

import "reflect"

var dialectMap = map[string]Dialect{}

// 接口，用于兼容多种数据库操作
type Dialect interface {
	// 将go类型映射为对应数据库类型
	DataTypeOf(typ reflect.Value) string
	// 判断表是否存在
	TableExistSQL(tableName string) (string, []interface{})
}

// 添加数据库
func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

// 获取对应数据库的处理
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}
