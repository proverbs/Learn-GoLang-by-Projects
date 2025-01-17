package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// original objects
	Writer http.ResponseWriter
	Req *http.Request
	// request info
	Path string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middlewares
	handlers []HandlerFunc
	idx int
	// back-ref to engine
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req: req,
		Path: req.URL.Path,
		Method: req.Method,
		idx: -1,
	}
}

func (c *Context) Next() {
	s:= len(c.handlers)
	for c.idx++; c.idx < s; c.idx++ {
		// provide more compatibility to handlers:
		// even if they don't call Next(), the chain will still work
		c.handlers[c.idx](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.idx = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	val, _ := c.Params[key]
	return val
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, val string) {
	c.Writer.Header().Set(key, val)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	// TODO: ... as arguments, ... to expand array
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	// TODO: json module
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, templateName string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, templateName, data); err != nil {
		c.Fail(500, err.Error())
	}
}
