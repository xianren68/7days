package RPC

import (
	"RPC/codec"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// Call 客户端调用rpc
type Call struct {
	Seq           uint64
	ServiceMethod string
	Args          any
	Reply         any
	Error         error
	// 用于断开请求
	Done chan *Call
}

func (call *Call) done() {
	call.Done <- call
}

// Client 客户端
type Client struct {
	cc       codec.Codec
	opt      *Option
	sending  sync.Mutex
	header   codec.Header
	mu       sync.Mutex
	seq      uint64
	pending  map[uint64]*Call
	closing  bool
	shutdown bool
}

var ErrShutdown = errors.New("connection is shutdown")

// Close 关闭连接
func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing {
		return ErrShutdown
	}
	return client.cc.Close()
}

// IsAvailable 连接是否能够使用
func (client *Client) IsAvailable() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return !client.shutdown && !client.closing
}

// 注册Call
func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.shutdown || client.closing {
		return 0, ErrShutdown
	}
	call.Seq = client.seq
	client.pending[call.Seq] = call
	client.seq++
	return call.Seq, nil
}

// 删除Call
func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	call := client.pending[seq]
	delete(client.pending, seq)
	return call
}

// Call 调用
func (client *Client) terminateCalls(err error) {
	client.sending.Lock()
	defer client.sending.Unlock()
	client.mu.Lock()
	defer client.mu.Unlock()
	client.shutdown = true
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}

// 接收服务端回复
func (client *Client) recive() {
	var err error
	for err == nil {
		var h codec.Header
		// 是否获取到响应头信息
		if err = client.cc.ReadHeader(&h); err != nil {
			break
		}
		// 取出对应的调用结构体
		call := client.removeCall(h.Seq)
		switch {
		case call == nil:
			err = client.cc.ReadBody(nil)
		case h.Error != "":
			call.Error = fmt.Errorf(h.Error)
			err = client.cc.ReadBody(nil)
			call.done()
		default:
			err = client.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
		client.terminateCalls(err)
	}
}

// NewClient 创建新客户端
func NewClient(conn net.Conn, opt *Option) (*Client, error) {
	f := codec.NewCodecMap[opt.CodecType]
	if f == nil {
		err := fmt.Errorf("invalid codec type %s", opt.CodecType)
		log.Println("rpc client: codec error:", err)
		return nil, err
	}
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("rpc client: encode option error:", err)
		_ = conn.Close()
		return nil, err
	}
	return newClientCodec(f(conn), opt), nil
}

func newClientCodec(cc codec.Codec, opt *Option) *Client {
	client := &Client{
		cc:      cc,
		opt:     opt,
		pending: make(map[uint64]*Call),
		// 请求序号从1开始
		seq: 1,
	}
	// 获取响应信息
	go client.recive()
	return client
}

// 解析option，并将其设计为可选参数
func parseOption(opts ...*Option) (*Option, error) {
	if len(opts) == 0 {
		return DefaultOption, nil
	}
	if len(opts) > 1 {
		return nil, errors.New("too many options")
	}
	opt := opts[0]
	opt.MagicNumber = DefaultOption.MagicNumber

	if opt.CodecType == "" {
		opt.CodecType = DefaultOption.CodecType
	}
	return opt, nil
}

// Dial 连接服务端
func Dial(network, address string, opts ...*Option) (client *Client, err error) {
	opt, err := parseOption(opts...)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()
	return NewClient(conn, opt)
}

func (client *Client) send(call *Call) {
	client.sending.Lock()
	defer client.sending.Unlock()
	// 注册
	seq, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}
	client.header.ServiceMethod = call.ServiceMethod
	client.header.Seq = seq
	client.header.Error = ""
	if err = client.cc.Write(&client.header, call.Args); err != nil {
		ca := client.removeCall(seq)
		if ca != nil {
			ca.Error = err
			ca.done()
		}
	}
}

func (client *Client) Go(ServiceMethod string, args, reply any, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 1)
	} else if cap(done) == 0 {
		log.Panic("rpc client:done channel is unbuffered")
	}
	call := new(Call)
	call.ServiceMethod = ServiceMethod
	call.Args = args
	call.Reply = reply
	call.Error = nil
	call.Done = done
	client.send(call)
	return call
}

// Call 调用后端方法
func (client *Client) Call(ServiceMethod string, args, reply any) error {
	call := <-client.Go(ServiceMethod, args, reply, nil).Done
	return call.Error
}
