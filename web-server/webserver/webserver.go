package webserver

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix string
	engine *Engine
	middlewares []HandlerFunc
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (rg *RouterGroup) Group(prefix string) *RouterGroup {
	engine := rg.engine
	newGroup := &RouterGroup{
		prefix: rg.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (rg *RouterGroup) addRoute(method string, pattern string, handler HandlerFunc) {
	fullPattern := rg.prefix + pattern
	log.Printf("Add route %4s - %s", method, fullPattern)
	rg.engine.router.addRoute(method, fullPattern, handler)
}

func (rg *RouterGroup) Use(middlewares ...HandlerFunc) {
	rg.middlewares = append(rg.middlewares, middlewares...)
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
	// This is not a good way to handle middlewares, because we cannot control the order
	// of the handlers.
	// My idea is to add middlewares to the trie, then we can collect all middlewares
	// from the root to the leaf(which follows first-in-first-out oder)
	// while searching the pattern in the trie.
	middlewares := make([]HandlerFunc, 0)
	for _, rg := range engine.groups {
		if strings.HasPrefix(req.URL.Path, rg.prefix) {
			middlewares = append(middlewares, rg.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
