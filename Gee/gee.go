package Gee

import (
	"net/http"
)

// 定义结构体,结构体包含router map,用于查找路由
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

// 定义方法，用于处理处理请求，它接收参数为Context实例
type HandlerFunc func(*Context)

// new方法，实例化Engine结构体
func New() *Engine {
	// 实例化
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 监听端口
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 定义结构体方法，满足http.Handler接口
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	// 分发路由
	engine.router.Handle(c)
}
