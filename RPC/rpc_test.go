package RPC

import (
	"fmt"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	addr := make(chan string)
	go startService(addr)
	client, _ := Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("rpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
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
