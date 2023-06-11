package session

import (
	"Orm/clause"
	"reflect"
)

// 插入语句
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)

	for _, value := range values {
		// 获取表
		table := s.Model(value).RefTable()
		// 添加sql语句
		s.clause.Set(clause.INSERT, table.Name, table.FiledNames)
		// 获取每个字段的值
		recordValues = append(recordValues, table.RecordVlues(value))
	}
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()

}

func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	s.clause.Set(clause.SELECT, table.Name, table.FiledNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	// 遍历返回值，构造对象
	for rows.Next() {
		// 创建对应结构体指针
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FiledNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()

}