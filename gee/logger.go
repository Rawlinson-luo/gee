package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		start := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v\n", c.StatusCode, c.Reader.RequestURI, time.Since(start))
	}
}
