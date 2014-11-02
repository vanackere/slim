package web

import (
	"net/http"

	"code.google.com/p/go.net/context"
)

type routeMachine struct {
	sm     stateMachine
	routes []route
}

func (rm routeMachine) route(c context.Context, w http.ResponseWriter, r *http.Request) (methodSet, context.Context, *route) {
	m := httpMethod(r.Method)
	var methods methodSet
	p := r.URL.Path

	if len(rm.sm) == 0 {
		return methods, c, nil
	}

	var i int
	for {
		sm := rm.sm[i].mode
		if sm&smSetCursor != 0 {
			si := rm.sm[i].i
			p = r.URL.Path[si:]
			i++
			continue
		}

		length := int(sm & smLengthMask)
		match := false
		if length <= len(p) {
			bs := rm.sm[i].bs
			switch length {
			case 3:
				if p[2] != bs[2] {
					break
				}
				fallthrough
			case 2:
				if p[1] != bs[1] {
					break
				}
				fallthrough
			case 1:
				if p[0] != bs[0] {
					break
				}
				fallthrough
			case 0:
				p = p[length:]
				match = true
			}
		}

		if match && sm&smRoute != 0 {
			si := rm.sm[i].i
			route := &rm.routes[si]
			if mc, ok := route.pattern.Match(r, c); ok {
				if route.method&m != 0 {
					return 0, mc, route
				}
				if m == mOPTIONS {
					methods |= methodSet(route.method)
				}
			}
			i++
		} else if match != (sm&smJumpOnMatch == 0) {
			if sm&smFail != 0 {
				return methods, c, nil
			}
			i = int(rm.sm[i].i)
		} else {
			i++
		}
	}

	return methods, c, nil
}
