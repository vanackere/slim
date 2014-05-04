package web

import (
	"crypto/rand"
	"encoding/base64"
	mrand "math/rand"
	"net/http"
	"testing"
	"time"
)

func init() {
	mrand.Seed(time.Now().Unix())
}

/*
The core benchmarks here are based on cypriss's mux benchmarks, which can be
found here:
https://github.com/cypriss/golang-mux-benchmark

They happen to play very well into Goji's router's strengths.
*/

type nilRouter struct{}

var helloWorld = []byte("Hello world!\n")

func (_ nilRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write(helloWorld)
}

type nilResponse struct{}

func (_ nilResponse) Write(buf []byte) (int, error) {
	return len(buf), nil
}
func (_ nilResponse) Header() http.Header {
	return nil
}
func (_ nilResponse) WriteHeader(code int) {
}

var w nilResponse

func addRoutes(m *Mux, prefix string) {
	m.Get(prefix, nilRouter{})
	m.Post(prefix, nilRouter{})
	m.Get(prefix+"/:id", nilRouter{})
	m.Put(prefix+"/:id", nilRouter{})
	m.Delete(prefix+"/:id", nilRouter{})
}

func randString() string {
	var buf [6]byte
	rand.Reader.Read(buf[:])
	return base64.URLEncoding.EncodeToString(buf[:])
}

func genPrefixes(n int) []string {
	p := make([]string, n)
	for i := range p {
		p[i] = "/" + randString()
	}
	return p
}

func genRequests(prefixes []string) []*http.Request {
	rs := make([]*http.Request, 5*len(prefixes))
	for i, prefix := range prefixes {
		rs[5*i+0], _ = http.NewRequest("GET", prefix, nil)
		rs[5*i+1], _ = http.NewRequest("POST", prefix, nil)
		rs[5*i+2], _ = http.NewRequest("GET", prefix+"/foo", nil)
		rs[5*i+3], _ = http.NewRequest("PUT", prefix+"/foo", nil)
		rs[5*i+4], _ = http.NewRequest("DELETE", prefix+"/foo", nil)
	}
	return rs
}

func permuteRequests(reqs []*http.Request) []*http.Request {
	out := make([]*http.Request, len(reqs))
	perm := mrand.Perm(len(reqs))
	for i, req := range reqs {
		out[perm[i]] = req
	}
	return out
}

func testingMux(n int) (*Mux, []*http.Request) {
	m := New()
	prefixes := genPrefixes(n)
	for _, prefix := range prefixes {
		addRoutes(m, prefix)
	}
	reqs := permuteRequests(genRequests(prefixes))
	return m, reqs
}

func BenchmarkRoute5(b *testing.B) {
	m, reqs := testingMux(1)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(w, reqs[i%len(reqs)])
	}
}
func BenchmarkRoute50(b *testing.B) {
	m, reqs := testingMux(10)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(w, reqs[i%len(reqs)])
	}
}
func BenchmarkRoute500(b *testing.B) {
	m, reqs := testingMux(100)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(w, reqs[i%len(reqs)])
	}
}
func BenchmarkRoute5000(b *testing.B) {
	m, reqs := testingMux(1000)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.ServeHTTP(w, reqs[i%len(reqs)])
	}
}