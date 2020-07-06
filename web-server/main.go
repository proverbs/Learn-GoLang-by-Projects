package main

import (
	"log"
	"net/http"
	"time"

	"proverbs.top/webserver"
)

func onlyForV2() webserver.HandlerFunc {
	return func(c *webserver.Context) {
		t := time.Now()
		// if something goes wrong:
		// c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	e := webserver.New()
	e.Use(webserver.Logger())

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
