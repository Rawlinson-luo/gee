package main

import (
	"fmt"
	"net/http"
)

type Engine struct{}

func (e Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.PATH=%q\n", r.URL.Path)
	case "/hello":
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND")
	}
}

func main() {
	e := Engine{}
	http.ListenAndServe(":9999", e)
}
