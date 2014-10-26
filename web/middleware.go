package web

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"code.google.com/p/go.net/context"
)

// mStack is an entire middleware stack. It contains a slice of middleware
// layers (outermost first) protected by a mutex, a cache of pre-built stack
// instances, and a final routing function.
type mStack struct {
	lock   sync.Mutex
	stack  []interface{}
	pool   *cPool
	router internalRouter
}

type internalRouter interface {
	route(context.Context, http.ResponseWriter, *http.Request)
}

/*
cStack is a cached middleware stack instance. Constructing a middleware stack
involves a lot of allocations: at the very least each layer will have to close
over the layer after (inside) it and a stack N levels deep will incur at least N
separate allocations. Instead of doing this on every request, we keep a pool of
pre-built stacks around for reuse.
*/
type cStack struct {
	ctx  context.Context
	m    Handler
	pool *cPool
}

func (s *cStack) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	s.m.ServeHTTPC(ctx, w, r)
}

func (s *cStack) toHTTPHandler(h HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(s.ctx, w, r)
	})
}

func (s *cStack) fromHTTPHandler(h http.HandlerFunc) HandlerFunc {
	return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		s.ctx = ctx
		h(w, r)
	})
}

func (m *mStack) appendLayer(fn interface{}) {
	switch fn.(type) {
	case func(http.Handler) http.Handler:
	case func(Handler) Handler:
	case func(context.Context, http.ResponseWriter, *http.Request, Handler /*next*/):
	default:
		log.Panicf(`Unknown middleware type %T. Expected a function `+
			`with signature "func(http.Handler) http.Handler" or `+
			`"func(web.Handler) web.Handler" or `+
			`"func(context.Context, http.ResponseWriter, *http.Request, Handler".`, fn)
	}
	m.stack = append(m.stack, fn)
}

func (m *mStack) findLayer(l interface{}) int {
	for i, middleware := range m.stack {
		if funcEqual(l, middleware) {
			return i
		}
	}
	return -1
}

func (m *mStack) invalidate() {
	m.pool = makeCPool()
}

type handlerNext struct {
	f    func(_ context.Context, _ http.ResponseWriter, _ *http.Request, next Handler)
	next HandlerFunc
}

func (h handlerNext) apply(c context.Context, w http.ResponseWriter, r *http.Request) {
	h.f(c, w, r, h.next)
}

func (m *mStack) newStack() *cStack {
	cs := cStack{}
	router := m.router

	h := HandlerFunc(router.route)

	for i := len(m.stack) - 1; i >= 0; i-- {
		switch fn := m.stack[i].(type) {
		case func(http.Handler) http.Handler:
			httphandler := cs.toHTTPHandler(h)
			h = cs.fromHTTPHandler(fn(httphandler).ServeHTTP)
		case func(Handler) Handler:
			h = fn(h).ServeHTTPC
		case func(context.Context, http.ResponseWriter, *http.Request, Handler /*next*/):
			h = HandlerFunc(handlerNext{fn, h}.apply)
		}
	}
	cs.m = h
	return &cs
}

func (m *mStack) alloc() *cStack {
	p := m.pool
	cs := p.alloc()
	if cs == nil {
		cs = m.newStack()
	}

	cs.pool = p
	return cs
}

func (m *mStack) release(cs *cStack) {
	cs.ctx = nil
	if cs.pool != m.pool {
		return
	}
	cs.pool.release(cs)
	cs.pool = nil
}

func (m *mStack) Use(middleware interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.appendLayer(middleware)
	m.invalidate()
}

func (m *mStack) Insert(middleware, before interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	i := m.findLayer(before)
	if i < 0 {
		return fmt.Errorf("web: unknown middleware %v", before)
	}

	m.appendLayer(middleware)
	inserted := m.stack[len(m.stack)-1]
	copy(m.stack[i+1:], m.stack[i:])
	m.stack[i] = inserted

	m.invalidate()
	return nil
}

func (m *mStack) Abandon(middleware interface{}) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	i := m.findLayer(middleware)
	if i < 0 {
		return fmt.Errorf("web: unknown middleware %v", middleware)
	}

	copy(m.stack[i:], m.stack[i+1:])
	m.stack = m.stack[:len(m.stack)-1 : len(m.stack)]

	m.invalidate()
	return nil
}
