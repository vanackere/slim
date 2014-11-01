package middleware

import (
	"net/http"

	"code.google.com/p/go.net/context"

	"github.com/vanackere/slim/web"
)

// SubRouter is a helper middleware that makes writing sub-routers easier.
//
// If you register a sub-router under a key like "/admin/*", Goji's router will
// automatically set c.URLParams["*"] to the unmatched path suffix. This
// middleware will help you set the request URL's Path to this unmatched suffix,
// allowing you to write sub-routers with no knowledge of what routes the parent
// router matches.
func SubRouter(ctx context.Context, w http.ResponseWriter, r *http.Request, next web.Handler) {
	if p := web.URLParams(ctx); p != nil {
		if path, ok := p["*"]; ok {
			oldpath := r.URL.Path
			r.URL.Path = path
			defer func() {
				r.URL.Path = oldpath
			}()
		}
	}
	next.ServeHTTPC(ctx, w, r)
}
