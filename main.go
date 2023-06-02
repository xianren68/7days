package main

import (
	"Gee"
	"fmt"
	"net/http"
)

func main() {
	// 创建Engine实例
	e := Gee.New()
	e.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	})
	e.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "你好")
	})
	e.Run(":9999")
}
