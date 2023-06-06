package main

import (
	"Gee"
	"log"
	"net/http"
)

func main() {
	// 创建Engine实例
	e := Gee.New()
	e.Use(func(c *Gee.Context) {
		log.Printf("%s", "global")
	})
	v1 := e.Group("/hello")
	{
		v1.Get("/name", func(ctx *Gee.Context) {
			ctx.String(http.StatusOK, "%s", ctx.Query("name"))
		})
		v1.Post("/name", func(ctx *Gee.Context) {
			ctx.Json(http.StatusOK, Gee.H{
				"name": ctx.PostForm("name"),
			})
		})
	}
	v2 := e.Group("/Hi")
	v2.Use(func(ctx *Gee.Context) {
		log.Printf("%s", "v2")
		ctx.Fail(404, "出错了")
	})
	{
		v2.Get("/age/:age", func(ctx *Gee.Context) {
			ctx.String(http.StatusOK, "%s", ctx.Param("age"))
		})
	}
	e.Run(":9999")

}
