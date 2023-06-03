package Gee

import (
	"log"
	"net/http"
)

// 定义router结构体
type Router struct {
	handlers map[string]HandlerFunc
}

// 创建实例
func NewRouter() *Router {
	return &Router{handlers: make(map[string]HandlerFunc)}
}

// 添加路由
func (router *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	// 添加日志
	log.Printf("Router %4s-%s", method, pattern)
	key := method + "-" + pattern
	router.handlers[key] = handler
}

// 处理路由
func (router *Router) Handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handle, ok := router.handlers[key]; ok {
		handle(c)
	} else {
		c.String(http.StatusNotFound, "404 Not Found:%s\n", c.Path)
	}
}
