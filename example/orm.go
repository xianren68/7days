package example

import (
	"Orm"

	_ "github.com/mattn/go-sqlite3"
)

type Student struct {
	Name string `Orm:"primary key"`
	age  int    `Orm:"not null"`
}

func OrmExample() {
	// 创建一个连接
	collect, err := Orm.NewEngine("sqlite3", "gee.db")
	defer collect.Close()
	if err != nil {
		return
	}
	// 创建一个会话
	s := collect.NewSession()
	s.Model(&Student{})
	// 创建表
	s.CreateTable()

}
