package Gee

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
