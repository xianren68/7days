package RPC

import (
	"RPC/codec"
	"go/ast"
	"log"
	"reflect"
	"sync/atomic"
)

// MagicNumber 校验连接正确
const MagicNumber = 0x3bef5c

// Option rpc连接的配置部分，用json编解码
type Option struct {
	MagicNumber int
	CodecType   string
}

// DefaultOption 默认配置
var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

// 服务方法
type methodType struct {
	method reflect.Method
	// 参数类型
	ArgType reflect.Type
	// 第二个参数
	ReplyType reflect.Type
	numsCall  uint64
}

func (m *methodType) NumCalls() uint64 {
	return atomic.LoadUint64(&m.numsCall)
}
func (m *methodType) newArgv() reflect.Value {
	var argv reflect.Value
	if m.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgType.Elem())
	} else {
		argv = reflect.New(m.ArgType).Elem()
	}
	return argv
}
func (m *methodType) newReply() reflect.Value {
	var reply = reflect.New(m.ReplyType.Elem())
	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
		reply.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		reply.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}
	return reply
}

// service 服务
type service struct {
	name   string
	typ    reflect.Type
	rcvr   reflect.Value
	method map[string]*methodType
}

func newService(rcvr any) *service {
	s := new(service)
	s.rcvr = reflect.ValueOf(rcvr)
	s.name = reflect.Indirect(s.rcvr).Type().Name()
	s.typ = reflect.TypeOf(rcvr)
	if !ast.IsExported(s.name) {
		log.Fatalf("rpc server: %s is not a valid service name", s.name)
	}
	s.registerMethods()
	return s

}
func (service *service) call(m *methodType, argv, reply reflect.Value) error {
	atomic.AddUint64(&m.numsCall, 1)
	f := m.method.Func
	returnValues := f.Call([]reflect.Value{service.rcvr, argv, reply})
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}

func (service *service) registerMethods() {
	service.method = make(map[string]*methodType)
	for i := 0; i < service.typ.NumMethod(); i++ {
		method := service.typ.Method(i)
		mType := method.Type
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		argType, replyType := mType.In(1), mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}
		service.method[method.Name] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
		log.Printf("rpc server: register %s.%s\n", service.name, method.Name)
	}
}

// 是否可导出
func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}
