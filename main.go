package main

import (
	"Gee"
	"net/http"
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
	e.Get("/xianren/:name", func(ctx *Gee.Context) {
		ctx.String(200, "%s,%s", ctx.Param("name"), ctx.Path)
	})
	e.Get("/assets/*filepath", func(ctx *Gee.Context) {
		ctx.Json(http.StatusOK, Gee.H{
			"filepath": ctx.Param("filepath"),
		})
	})
	e.Run(":9999")

}
