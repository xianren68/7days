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

// 接口，用户只需要将操作通过回调函数传入，会自动执行事务
type TxFunc func(*session.Session) (interface{}, error)

func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	// 创建新的会话
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			// 出错，回滚
			_ = s.RollBack()
			panic(p)
		} else if err != nil {
			// 执行出错，回滚
			_ = s.RollBack()

		} else {
			// 没问题，提交
			err = s.Commit()
		}
	}()
	return f(s)

}
