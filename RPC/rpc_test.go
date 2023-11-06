package RPC

import (
	"RPC/codec"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	addr := make(chan string)
	go startService(addr)
	dial, _ := net.Dial("tcp", <-addr)
	defer func() { _ = dial.Close() }()
	time.Sleep(time.Second)
	// 发送option
	_ = json.NewEncoder(dial).Encode(DefaultOption)
	cc := codec.NewGobCodec(dial)
	for i := 0; i < 5; i++ {

		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply:", reply)
	}
}

// 开启rpc服务
func startService(addr chan string) {
	// 寻找一个空闲的接口
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Println("network err ", err)
		return
	}
	log.Println("start rpc service ", l.Addr())
	addr <- l.Addr().String()
	// 开启rpc端口监听
	Accept(l)
}
