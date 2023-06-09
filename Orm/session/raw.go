package session

import (
	"Orm/log"
	"database/sql"
	"strings"
)

type Session struct {
	db *sql.DB // 数据库连接
	// sql语句
	sql     strings.Builder
	sqlVars []interface{}
}

// 创建会话实例
func New(db *sql.DB) *Session {
	return &Session{
		db: db,
	}
}

// 清除
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

// 获取数据库连接
func (s *Session) DB() *sql.DB {
	return s.db
}

// 存储sql语句与占位符
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// 执行语句
func (s *Session) Exec() (result sql.Result, err error) {
	// 清除会话存储的sql语句
	defer s.Clear()
	// 输出执行语句的日志
	log.Info(s.sql.String(), s.sqlVars)
	// 判断执行是否成功
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// 查询语句(一条)
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)

}

// 查询语句（多条）
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
