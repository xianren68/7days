package RPC

import (
	"RPC/codec"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
)

type Server struct {
	serviceMap sync.Map
}

func (server *Server) Register(rcvr any) error {
	s := newService(rcvr)
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("rpc: service already defined: " + s.name)
	}
	return nil
}
func Register(rcvr any) error {
	return DefaultServer.Register(rcvr)
}
func (server *Server) findService(serviceMethod string) (svc *service, method *methodType, err error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc server: service/method request ill-formed: " + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	svci, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}
	svc = svci.(*service)
	method = svc.method[methodName]
	if method == nil {
		err = errors.New("rpc server: can't find method " + methodName)
	}
	return
}

// NewServer 客户端构造函数
func NewServer() *Server {
	return &Server{}
}

// DefaultServer 客户端默认实现
var DefaultServer = NewServer()

// Accept 接收客户端连接
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc listen: accept error", err)
		}
		go server.ServerConn(conn)
	}
}

// Accept 直接监听客户端连接
func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

// ServerConn 处理rpc连接
func (server *Server) ServerConn(conn io.ReadWriteCloser) {
	// 尝试解码option
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("option err", err)
		return
	}
	// 校验连接是否正确
	if opt.MagicNumber != MagicNumber {
		log.Println("rpc server:invalid magic number", opt.MagicNumber)
		return
	}
	f := codec.NewCodecMap[opt.CodecType]
	// 判断是否存在对应的编解码器
	if f == nil {
		log.Println("rpc server:invalid codec type", opt.CodecType)
		return
	}
	server.serverCodec(f(conn))

}
func (server *Server) serverCodec(cc codec.Codec) {
	// 保证数据发送时的安全，避免客户端接收的消息无序
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, sending, req.h, struct{}{})
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

// 请求结构体
type request struct {
	// 请求头
	h *codec.Header
	// 参数
	argv reflect.Value
	// 回复
	replyv reflect.Value
	// 方法
	method *methodType
	// 服务
	svc *service
}

// 获取请求头
func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		// 不是到了末尾
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server:read Header err", err)
		}
		return nil, err
	}
	return &h, nil
}

// 获取请求信息
func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	req.svc, req.method, err = server.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}
	req.argv = req.method.newArgv()
	req.replyv = req.method.newReply()
	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}
	if err = cc.ReadBody(argvi); err != nil {
		log.Println("rpc service: read body error", err)
	}
	return req, err
}

// 发送回复信息
func (server *Server) sendResponse(cc codec.Codec, sending *sync.Mutex, h *codec.Header, body any) {
	// 上锁，保证并发安全
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc service:write response error", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	err := req.svc.call(req.method, req.argv, req.replyv)
	if err != nil {
		req.h.Error = err.Error()
		server.sendResponse(cc, sending, req.h, struct{}{})
		return
	}
	server.sendResponse(cc, sending, req.h, req.replyv.Interface())
}
