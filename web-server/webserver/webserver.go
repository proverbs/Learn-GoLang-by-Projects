package webserver

import (
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix string
	engine *Engine
}

type Engine struct {
	*RouterGroup
	router *router
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	return engine
}

func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	engine := rg.engine
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix,
		engine: engine,
	}
	return newGroup
}

func (rg *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	fullPattern := rg.prefix + pattern
	log.Printf("Add route %4s - %s", method, fullPattern)
	rg.engine.router.addRoute(method, fullPattern, handler)
}

func (rg *RouterGroup) GET(pattern string, handler HandlerFunc) {
	rg.addRoute("GET", pattern, handler)
}

func (rg *RouterGroup) POST(pattern string, handler HandlerFunc) {
	rg.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	// use our engine
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
