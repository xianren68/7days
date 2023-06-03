package Gee

import (
	"net/http"
)

// 定义结构体,结构体包含router map,用于查找路由
type Engine struct {
	router *Router
}

// 定义方法，用于处理处理请求，它接收参数为Context实例
type HandlerFunc func(*Context)

// new方法，实例化Engine结构体
func New() *Engine {
	// 实例化
	return &Engine{router: NewRouter()}
}

// @param method 请求方法
// @param pattern 请求路径
// @param handler 处理函数
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// 定义get方法
// @param pattern 请求路径
// @param handler 处理函数
func (engine *Engine) Get(pattern string, handler HandlerFunc) {
	// 将路由处理函数与其对应的模式添加到到router map中
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) Post(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 监听端口
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 定义结构体方法，满足http.Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	// 分发路由
	engine.router.Handle(c)
}
