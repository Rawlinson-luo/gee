package gee

import (
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type Engine struct {
	router *Router
	groups []*RouteGroup
	*RouteGroup
}

type RouteGroup struct {
	prefix      string
	middlewares []HandlerFunc
	engine      *Engine
	parent      *RouteGroup
}

func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouteGroup = &RouteGroup{engine: engine}
	engine.groups = []*RouteGroup{engine.RouteGroup}
	return engine
}

func (group *RouteGroup) Group(prefix string) *RouteGroup {
	g := &RouteGroup{
		prefix: group.prefix + prefix,
		engine: group.engine,
		parent: group,
	}
	group.engine.groups = append(group.engine.groups, g)
	return g
}

func (group *RouteGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouteGroup) addRoute(method string, comp string, handlerFunc HandlerFunc) {
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern, handlerFunc)
}

func (group *RouteGroup) GET(pattern string, handlerFunc HandlerFunc) {
	group.addRoute("GET", pattern, handlerFunc)
}

func (group *RouteGroup) POST(pattern string, handlerFunc HandlerFunc) {
	group.addRoute("POST", pattern, handlerFunc)
}

func (e *Engine) GET(pattern string, handlerFunc HandlerFunc) {
	e.router.addRoute("GET", pattern, handlerFunc)
}

func (e *Engine) POST(pattern string, handlerFunc HandlerFunc) {
	e.router.addRoute("POST", pattern, handlerFunc)
}

func (e *Engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	middlewares := make([]HandlerFunc, 0)

	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	context := newContext(w, r)
	context.handlers = middlewares
	e.router.handle(context)
}
