package webserver

import (
	"html/template"
	"log"
	"net/http"
	"path"
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
	htmlTemplates *template.Template
	funcMap template.FuncMap
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

func (rg *RouterGroup) createStaticHandler(routePath string, fs http.FileSystem) HandlerFunc {
	routePath = path.Join(rg.prefix, routePath)
	fileServer := http.StripPrefix(routePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
} 

func (rg *RouterGroup) Static(routePath string, root string) {
	handler := rg.createStaticHandler(routePath, http.Dir(root))
	pattern := path.Join(routePath, "/*filepath")
	rg.GET(pattern, handler)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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
	c.engine = engine
	engine.router.handle(c)
}
