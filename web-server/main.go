package main

import (
	"net/http"

	"proverbs.top/webserver"
)

func main() {
	r := webserver.New()

	r.GET("/", func(c *webserver.Context) {
		c.HTML(http.StatusOK, "<h1>Hello, Proverbs!</h1>")
	})

	r.GET("/hello", func(c *webserver.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *webserver.Context) {
		c.JSON(http.StatusOK, webserver.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
