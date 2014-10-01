/*
Package slim provides an out-of-box web server with reasonable defaults.

Example:
	package main

	import (
		"fmt"
		"net/http"

		"code.google.com/p/go.net/context"

		"github.com/vanackere/slim"
		"github.com/vanackere/slim/web"
	)

	func hello(c context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", web.URLParams(c)["name"])
	}

	func main() {
		slim.Get("/hello/:name", hello)
		slim.Serve()
	}

This package exists purely as a convenience to programmers who want to get
started as quickly as possible. It draws almost all of its code from slim's
subpackages, the most interesting of which is slim/web, and where most of the
documentation for the web framework lives.

A side effect of this package's ease-of-use is the fact that it is opinionated.
If you don't like (or have outgrown) its opinions, it should be straightforward
to use the APIs of slim's subpackages to reimplement things to your liking. Both
methods of using this library are equally well supported.

Slim requires Go 1.2 or newer.
*/
package slim
