package main

import (
	"net/http"

	"proverbs.top/webserver"
)

func main() {
	e := webserver.New()

	e.GET("/index", func(c *webserver.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	g1 := e.Group("/v1")
	{
		g1.GET("/", func(c *webserver.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		g1.GET("/hello", func(c *webserver.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	g2 := e.Group("/v2")
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
