package RPC

import (
	"RPC/codec"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
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

// Service 服务端
type Service struct {
}

// NewService 客户端构造函数
func NewService() *Service {
	return &Service{}
}

// DefaultService 客户端默认实现
var DefaultService = NewService()

// Accept 接收客户端连接
func (service *Service) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc listen: accept error", err)
		}
		go service.ServiceConn(conn)
	}
}

// Accept 直接监听客户端连接
func Accept(lis net.Listener) {
	DefaultService.Accept(lis)
}

// ServiceConn 处理rpc连接
func (service *Service) ServiceConn(conn io.ReadWriteCloser) {
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
	service.serviceCodec(f(conn))

}
func (service *Service) serviceCodec(cc codec.Codec) {
	// 保证数据发送时的安全，避免客户端接收的消息无序
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := service.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			service.sendResponse(cc, sending, req.h, struct{}{})
			continue
		}
		wg.Add(1)
		go service.handleRequest(cc, req, sending, wg)
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
}

// 获取请求头
func (service *Service) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
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
func (service *Service) readRequest(cc codec.Codec) (*request, error) {
	h, err := service.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	// 暂时只支持string
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc service: read body error", err)
	}
	return req, err
}

// 发送回复信息
func (service *Service) sendResponse(cc codec.Codec, sending *sync.Mutex, h *codec.Header, body any) {
	// 上锁，保证并发安全
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc service:write response error", err)
	}
}

func (service *Service) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("rpc resp %d", req.h.Seq))
	service.sendResponse(cc, sending, req.h, req.replyv.Interface())
}
