package Gee

import (
	"fmt"
	"net/http"
)

// 定义结构体,结构体包含router map,用于查找路由
type Engine struct {
	router map[string]http.HandlerFunc
}

// new方法，实例化Engine结构体
func New() *Engine {
	// 实例化
	return &Engine{router: make(map[string]http.HandlerFunc)}
}

// @param method 请求方法
// @param pattern 请求路径
// @param handler 处理函数
func (engine *Engine) addRoute(method string, pattern string, handler http.HandlerFunc) {
	key := method + "-" + pattern
	// 将路由添加到router map中
	engine.router[key] = handler
}

// 定义get方法
// @param pattern 请求路径
// @param handler 处理函数
func (engine *Engine) Get(pattern string, handler http.HandlerFunc) {
	// 将路由处理函数与其对应的模式添加到到router map中
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) Post(pattern string, handler http.HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 监听端口
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 定义结构体方法，满足http.Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		// 执行对应的处理函数
		handler(w, req)
	} else {
		// 返回失败状态码
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
