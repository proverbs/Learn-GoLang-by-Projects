package main

import (
	"html/template"
	"log"
	"fmt"
	"net/http"
	"time"

	"proverbs.top/webserver"
)

type student struct {
	Name string
	Age  int8
}

func onlyForV2() webserver.HandlerFunc {
	return func(c *webserver.Context) {
		t := time.Now()
		// if something goes wrong:
		// c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func formatAsDate(t time.Time, more string) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d, %s", year, month, day, more)
}

func main() {
	e := webserver.New()
	e.Use(webserver.Logger())
	e.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	e.LoadHTMLGlob("templates/*")
	e.Static("/assets", "./static")

	stu1 := &student{Name: "Fxxk", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	e.GET("/", func(c *webserver.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	e.GET("/students", func(c *webserver.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", webserver.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	e.GET("/date", func(c *webserver.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", webserver.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
			"more": "lalala",
		})
	})

	g1 := e.Group("/v1")
	{
		g1.GET("/hello", func(c *webserver.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	g2 := e.Group("/v2")
	g2.Use(onlyForV2())
	{
		g2.GET("/hello/:name", func(c *webserver.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		g2.GET("/login", func(c *webserver.Context) {
			c.JSON(http.StatusOK, webserver.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	e.Run(":9999")
}
