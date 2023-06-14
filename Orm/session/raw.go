package session

import (
	"Orm/clause"
	"Orm/dialect"
	"Orm/log"
	"Orm/schema"
	"database/sql"
	"strings"
)

type Session struct {
	db *sql.DB // 数据库连接
	// sql语句
	sql     strings.Builder
	sqlVars []interface{}
	tx      *sql.Tx // 数据库事务
	// 对应的数据库类型操作
	dialect  dialect.Dialect
	refTable *schema.Schema
	clause   clause.Clause
}

// 创建会话实例
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

// 清除
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

// 获取数据库连接
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
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
