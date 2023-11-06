package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

var _ Codec = (*GobCodec)(nil)

// GobCodec 编解码Gob
type GobCodec struct {
	conn io.ReadWriteCloser
	// 缓冲区
	buf *bufio.Writer
	// 解码器
	dec *gob.Decoder
	// 编码器
	enc *gob.Encoder
}

func (g *GobCodec) Close() error {
	return g.conn.Close()

}

func (g *GobCodec) ReadHeader(header *Header) error {
	// 解码消息头信息
	return g.dec.Decode(header)
}

func (g *GobCodec) ReadBody(body any) error {
	// 解码消息体信息
	return g.dec.Decode(body)
}

func (g *GobCodec) Write(header *Header, body any) (err error) {
	defer func() {
		_ = g.buf.Flush()
		if err != nil {
			_ = g.Close()
		}
	}()
	// 尝试编码消息头
	if err = g.enc.Encode(header); err != nil {
		log.Println("rpc codec: gob error encoding header")
		return
	}
	// 尝试编码消息体
	if err = g.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body")
		return
	}
	return
}

// NewGobCodec 构造函数
func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(conn),
	}
}
