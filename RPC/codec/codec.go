package codec

import "io"

// Header RPC消息头部
type Header struct {
	// 服务与方法名
	ServiceMethod string
	// 请求序列号
	Seq uint64
	// 错误信息
	Error string
}

// Codec 消息编解码接口
type Codec interface {
	io.Closer
	// ReadHeader 读取消息头
	ReadHeader(header *Header) error
	// ReadBody 读取消息体
	ReadBody(any) error
	// Write 写入消息头和消息体
	Write(header *Header, body any) error
}

// NewCodecFunc 构造函数
type NewCodecFunc func(io.ReadWriteCloser) Codec

const (
	GobType  = "application/gob"
	JsonType = "application/json"
)

// NewCodecMap 实现的编解码类型
var NewCodecMap map[string]NewCodecFunc

func init() {
	NewCodecMap = make(map[string]NewCodecFunc)
	NewCodecMap[GobType] = NewGobCodec
}
