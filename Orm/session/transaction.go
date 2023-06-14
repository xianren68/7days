package session

import "Orm/log"

// 封装事务方法
func (s *Session) Begin() (err error) {
	log.Info("transaction begin")
	// 开始事务，并将事务对象赋值给s.tx
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

//
func (s *Session) Commit() (err error) {
	log.Info("transcation commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session) RollBack() (err error) {
	log.Info("transcation rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
		return
	}
	return
}
