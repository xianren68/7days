package main

import (
	"Gee"
	"fmt"
	"net/http"
	"text/template"
	"time"
)

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := Gee.New()
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")
	r.Get("/", func(c *Gee.Context) {
		c.Html(http.StatusOK, "test.tmpl", nil)
	})
	r.Get("/students", func(c *Gee.Context) {
		var x []int
		c.String(200, "%d", x[10])
	})

	r.Get("/date", func(c *Gee.Context) {
		c.Html(http.StatusOK, "test.tmpl", Gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":9999")
}
