package session

import (
	"Orm/log"
	"Orm/schema"
	"fmt"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) *Session {
	// 没有表的映射或者表对应的结构体发生改变
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// 获取会话对应的表模型
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

// 创建表
func (s *Session) CreateTable() error {
	table := s.refTable
	columns := make([]string, 0)
	// 拼接sql语句
	for _, filed := range table.Fileds {
		columns = append(columns, fmt.Sprintf("%s %s %s", filed.Name, filed.Type, filed.Tag))
	}
	desc := strings.Join(columns, ",")
	// 这里明白为什么有时候需要返回结构体本身，为了方便链式调用
	_, err := s.Raw(fmt.Sprintf("create Table %s(%s)", table.Name, desc)).Exec()
	return err
}

// 删除表
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("drop table if exists %s", s.RefTable().Name)).Exec()
	return err

}

// 判断表是否存在
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
