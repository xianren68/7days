package Gee

import (
	"net/http"
	"path"
)

// 路由分组
type RouterGroup struct {
	prefix      string        // 分组前缀
	middlewares []HandlerFunc // 中间件数组
	engine      *Engine
}

// 新建分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: engine.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	// 拼接路径
	pattern := group.prefix + comp
	// 调用router添加路由
	group.engine.router.addRoute(method, pattern, handler)
}

// Get方法
func (group *RouterGroup) Get(path string, handler HandlerFunc) {
	group.addRoute("GET", path, handler)
}

// Post
func (group *RouterGroup) Post(path string, handler HandlerFunc) {
	group.addRoute("POST", path, handler)
}

// 设置use方法，用于使用中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) { // 可以传递多个中间件
	// 添加中间件到middlewares
	group.middlewares = append(group.middlewares, middlewares...)
}

// 设置静态路由对应的处理方法
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	// 拼接路径
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// 判断文件是否可以被操作
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 设置静态目录
// @params relativePath 路由路径
// @params root 真实路径
func (group *RouterGroup) Static(relativePath, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "*filepath")
	group.Get(urlPattern, handler)
}
