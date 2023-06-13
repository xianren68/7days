package session

import (
	"Orm/clause"
	"errors"
	"reflect"
)

// 插入语句
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)

	for _, value := range values {
		s.CallMethod("BeforeInsert", value)
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
	// 查询前钩子
	// s.CallMethod("BeforeQuery", nil)
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
		// 查询后钩子
		s.CallMethod("AfterQuery", dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()

}

// 更新操作
func (s *Session) Update(values ...interface{}) (int64, error) {
	m, ok := values[0].(map[string]interface{})
	if !ok { // 传入的不是map类型
		m = make(map[string]interface{})
		for i := 0; i < len(values); i += 2 {
			m[values[i].(string)] = values[i+1]
		}
	}

	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// 删除操作
func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	result := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := result.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

// limit语句
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s // 返回实例本身，用于链式调用
}
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	vars = append(vars, desc)
	vars = append(vars, args...)
	s.clause.Set(clause.WHERE, vars...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 { //不存在
		return errors.New("Not found")

	}
	dest.Set(destSlice.Index(0)) // 将返回结果赋给结构体
	return nil
}
