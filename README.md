Slim
====

[![GoDoc](https://godoc.org/github.com/vanackere/slim/web?status.svg)](https://godoc.org/github.com/vanackere/slim/web) [![Build Status](https://travis-ci.org/vanackere/slim.svg)](https://travis-ci.org/vanackere/slim)

Slim is a fork of [Goji][goji], a minimalistic web framework that values composability and simplicity.

[goji]: https://github.com/zenazn/goji

Differences with Goji
---------------------

 * Slim opted to implement its request context using ["code.google.com/p/go.net/context"][context] (see [here][cblog] for an introduction).
 * Slim middlewares can also be written using the simpler form:
``` go 
func(ctx context.Context, w http.ResponseWriter, r *http.Request, next web.Handler)
```

[context]: https://code.google.com/p/go.net/context

Example
-------

```go
package main

import (
        "fmt"
        "net/http"

        "code.google.com/p/go.net/context"
        "github.com/vanackere/slim"
        "github.com/vanackere/slim/web"
)

func hello(ctx context.Context, w http.ResponseWriter, r *http.Request) {
        p := web.URLParams(ctx)
        fmt.Fprintf(w, "Hello, %s!", p["name"])
}

func main() {
        slim.Get("/hello/:name", hello)
        slim.Serve()
}
```

Slim also includes a [sample application][sample] in the `example` folder which
was artificially constructed to show off all of Slim's features. Check it out!

[sample]: https://github.com/vanackere/slim/tree/master/example


Features
--------

* Fork of the excellent [Goji framework][goji]
* Compatible with `net/http`
* URL patterns (both Sinatra style `/foo/:bar` patterns and regular expressions,
  as well as [custom patterns][pattern])
* Reconfigurable middleware stack
* Context/environment object threaded through middleware and handlers
* Context is the one from "code.google.com/p/go.net/context" (see [here][cblog] for an introduction)
* Automatic support for [Einhorn][einhorn], systemd, and [more][bind]
* [Graceful shutdown][graceful], and zero-downtime graceful reload when combined
  with Einhorn.

[cblog]: https://blog.golang.org/context
[einhorn]: https://github.com/stripe/einhorn
[bind]: https://godoc.org/github.com/vanackere/slim/bind
[graceful]: https://godoc.org/github.com/vanackere/slim/graceful
[pattern]: https://godoc.org/github.com/vanackere/slim/web#Pattern

See [Goji's README][readme] for why Goji's - and therore Slim's ! - approach is good.

[readme]: https://github.com/zenazn/goji/blob/master/README.md
