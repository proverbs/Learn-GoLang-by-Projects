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

	r.GET("/hello/:name", func(c *webserver.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *webserver.Context) {
		c.JSON(http.StatusOK, webserver.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}
