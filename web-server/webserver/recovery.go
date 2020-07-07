package webserver

import (
	"log"
	"net/http"
	"runtime/debug"
)


func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("stacktrace from panic: \n" + string(debug.Stack()))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}