package session

import (
	"fmt"
	"reflect"
)

// 钩子函数常量
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

type IAfterQuery interface {
	AfterQuery(s *Session) error
}

type IBeforeQuery interface {
	BeforeQuery(s *Session) error
}
type IBeforeInsert interface {
	BeforeInsert(s *Session) error
}

// 调用钩子函数
func (s *Session) CallMethod(method string, value interface{}) {
	param := reflect.ValueOf(value)
	switch method {
	case AfterQuery:
		if i, ok := param.Interface().(IAfterQuery); ok {
			i.AfterQuery(s)
		}
	case BeforeQuery:
		if i, ok := param.Interface().(IBeforeQuery); ok {
			i.BeforeQuery(s)
		}
	case BeforeInsert:
		if i, ok := param.Interface().(IBeforeInsert); ok {
			i.BeforeInsert(s)
		}
	default:
		panic(fmt.Sprintf("unsupported this hook method -> %s", method))
	}

}
