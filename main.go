package main

import (
	"gee"
	"net/http"
)

func main() {
	g := gee.New()
	g.Use(gee.Logger())
	g.Use(gee.Recovery())

	g.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "URL.PATH=%q\n", c.Reader.URL.Path)
	})

	g.GET("/ye", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you are at %s", c.Query("name"), c.Path)
	})

	g.POST("/login", func(c *gee.Context) {
		c.Json(http.StatusOK, map[string]interface{}{
			"userName": c.PostForm("userName"),
			"password": c.PostForm("password"),
		})
	})

	g.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello name = %s\n", c.Param("name"))
	})

	g.GET("/hello/name/*", func(c *gee.Context) {
		c.String(http.StatusOK, "dup request\n")
	})

	g.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	//分组
	//g1 := g.Group("/luo")
	//{
	//	g1.GET("/hello", func(c *gee.Context) {
	//		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
	//	})
	//}

	g.Run(":9999")

}
