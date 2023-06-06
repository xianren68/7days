package Gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 给任意类型map取别名为H
type H = map[string]interface{}

// 定义context结构体
type Context struct {
	// http 写入流
	Writer http.ResponseWriter
	// http 请求信息
	Req *http.Request
	// 请求方法
	Method string
	// 请求路径
	Path string
	// 返回状态码
	StatusCode int
	// params参数
	Params map[string]string
	// 中间件列表
	handlers []HandlerFunc
	// 当前在执行哪个中间件
	index int
}

// 创建context实例
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Method: req.Method,
		Path:   req.URL.Path,
		index:  -1,
	}
}

// 获取query参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 获取post参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 获取params参数
func (c *Context) Param(key string) string {
	value := c.Params[key]
	return value
}

// 设置请求头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 返回string字符串
func (c *Context) String(code int, format string, values ...interface{}) {
	// 设置返回数据为字符
	c.SetHeader("Content-Type", "text/plain")
	// 设置状态码
	c.Status(code)
	// 写入数据
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 返回json数据
func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Context-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 返回二进制数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// 返回html数据
func (c *Context) Html(code int, html string) {
	c.SetHeader("Context-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

// 设置next方法
func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ { //c.index会改变不用担心函数重复执行
		c.handlers[c.index](c)
	}
}

// 设置Fail方法,中止中间件的执行
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers) // 跳到最后，没有函数可执行
	c.Json(code, H{
		"err": err,
	})
}
