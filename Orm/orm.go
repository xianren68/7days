package Orm

import (
	"Orm/dialect"
	"Orm/log"
	"Orm/session"
	"database/sql"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(dirver, source string) (e *Engine, err error) {
	db, err := sql.Open(dirver, source)
	if err != nil {
		log.Error(err)
		return
	}
	// 判断数据库是否能够正常连接
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	//
	dial, ok := dialect.GetDialect(dirver)
	if !ok {
		log.Errorf("dialect %s not Found", dirver)
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Collection database success")
	return
}

// 关闭数据库连接
func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// 创建一个新的会话
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}
