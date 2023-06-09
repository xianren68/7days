package example

import (
	"Orm"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func OrmExample() {
	// 创建一个连接
	collect, err := Orm.NewEngine("sqlite3", "gee.db")
	defer collect.Close()
	if err != nil {
		return
	}
	// 创建一个会话
	s := collect.NewSession()
	// 执行语句
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
