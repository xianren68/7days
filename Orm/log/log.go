package log

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	// 错误与信息日志，会输出文件代码行数以及时间戳
	errorlog = log.New(os.Stdout, "Error", log.LstdFlags|log.Lshortfile)
	infolog  = log.New(os.Stdout, "Info", log.LstdFlags|log.Lshortfile)
	// 日志实例列表
	loggers = []*log.Logger{errorlog, infolog}
	// 锁
	mu = sync.Mutex{}
)

// 输出日志的方法
var (
	Error  = errorlog.Println
	Errorf = errorlog.Printf
	Info   = infolog.Println
	Infof  = infolog.Printf
)

// 设置日志层级
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()
	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}
	if ErrorLevel < level {
		errorlog.SetOutput(io.Discard)
	}
	if InfoLevel < level {
		infolog.SetOutput(io.Discard)
	}

}
