package RPC

import (
	"log"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)
	client, _ := Dial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

// 开启rpc服务
func startServer(addr chan string) {
	var foo Foo
	if err := Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}
	// pick a free port
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", l.Addr())
	addr <- l.Addr().String()
	Accept(l)
}

type Foo int
type Args struct {
	Num1, Num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num2 + args.Num1
	return nil
}
func TestNewService(t *testing.T) {
	var foo Foo
	service := newService(&foo)
	log.Println(len(service.method) == 1)
	log.Println(service.method["Sum"] != nil)
}
func TestService_Call(t *testing.T) {
	var foo Foo
	service := newService(&foo)
	mType := service.method["Sum"]
	argv := mType.newArgv()
	reply := mType.newReply()
	argv.Set(reflect.ValueOf(Args{Num1: 1, Num2: 2}))
	err := service.call(mType, argv, reply)
	log.Println(err == nil && *reply.Interface().(*int) == 3 && mType.NumCalls() == 1)
}
