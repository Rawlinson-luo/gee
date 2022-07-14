package gee

import (
	"log"
	"net/http"
	"strings"
)

type Router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func NewRouter() *Router {
	return &Router{
		handlers: map[string]HandlerFunc{},
		roots:    map[string]*node{},
	}
}

func (r *Router) addRoute(method string, pattern string, handlerFunc HandlerFunc) {
	log.Printf("method=%s pattern=%s", method, pattern)

	parts := parsePattern(pattern)

	key := method + "-" + pattern

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handlerFunc
}

//查询route
func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	root, ok := r.roots[method]
	params := make(map[string]string)
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for i, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[i]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}

		return n, params
	}
	return nil, nil
}

func parsePattern(pattern string) []string {
	parts := make([]string, 0)
	vs := strings.Split(pattern, "/")

	for _, v := range vs {
		if v != "" {
			parts = append(parts, v)
			//字符串取位置使用单引号转为ASCII码比较
			if v[0] == '*' {
				break
			}
		}
	}
	return parts

}

func (r *Router) handle(c *Context) {

	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
