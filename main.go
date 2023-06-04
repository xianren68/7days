package main

import (
	"Gee"
)

func main() {
	// 创建Engine实例
	e := Gee.New()
	e.Get("/", func(ctx *Gee.Context) {
		ctx.String(200, "%s", ctx.Query("name"))
	})
	e.Post("/hello", func(ctx *Gee.Context) {
		ctx.String(200, "%s", ctx.PostForm("name"))
	})
	e.Run(":9999")

}
