package Gee

import (
	"log"
	"net/http"
	"strings"
)

// 定义router结构体
type router struct {
	roots    map[string]*node       // key为方法，
	handlers map[string]HandlerFunc // key为路由路径
}

// 创建实例
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 生成parts，通配符后的路径无效
func parsePattern(pattern string) []string {
	// 按照‘/’分割
	vs := strings.Split(pattern, "/")
	// 返回的parts数组
	parts := make([]string, 0)
	for _, val := range vs {
		if val != "" {
			parts = append(parts, val)
			if val[0] == '*' { // 通配符
				break
			}
		}
	}
	return parts
}

// 添加路由
func (router *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	_, ok := router.roots[method] // 判断此类型的请求是否有前缀树
	if !ok {
		router.roots[method] = &node{}
	}

	key := method + "-" + pattern
	router.handlers[key] = handler
	router.roots[method].insert(pattern, parts, 0)
	// 输出日志
	log.Printf("Router %4s-%s", method, pattern)
}

// 获取对应的路由路径,并返回param参数
func (router *router) getRouter(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := router.roots[method]
	if !ok { // 如果没有此方法对应的前缀树
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index] // 获取param参数
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil

}

// 匹配路由处理函数
func (router *router) Handle(c *Context) {
	node, params := router.getRouter(c.Method, c.Path)
	// 输出日志
	log.Printf("request:%s-%s", c.Method, c.Path)
	if node != nil {
		c.Params = params
		key := c.Method + "-" + node.pattern
		router.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 Not Found:%s", c.Path)
	}
}
