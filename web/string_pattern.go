package web

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"code.google.com/p/go.net/context"
)

type stringPattern struct {
	raw      string
	pats     []string
	literals []string
	isPrefix bool
}

func (s stringPattern) Prefix() string {
	return s.literals[0]
}
func (s stringPattern) Match(r *http.Request, c context.Context) (context.Context, bool) {
	if _, ok := s.match(r, c, true); !ok {
		return c, ok
	}
	return s.match(r, c, false)
}
func (s stringPattern) match(r *http.Request, c context.Context, dryrun bool) (context.Context, bool) {
	path := r.URL.Path
	var matches map[string]string
	if !dryrun && len(s.pats) > 0 {
		matches = make(map[string]string, len(s.pats))
	}
	for i := 0; i < len(s.pats); i++ {
		sli := s.literals[i]
		if !strings.HasPrefix(path, sli) {
			return c, false
		}
		path = path[len(sli):]

		m := 0
		for ; m < len(path); m++ {
			if path[m] == '/' {
				break
			}
		}
		if m == 0 {
			// Empty strings are not matches, otherwise routes like
			// "/:foo" would match the path "/"
			return c, false
		}
		if !dryrun {
			matches[s.pats[i]] = path[:m]
		}
		path = path[m:]
	}
	// There's exactly one more literal than pat.
	if s.isPrefix {
		if !strings.HasPrefix(path, s.literals[len(s.pats)]) {
			return c, false
		}
	} else {
		if path != s.literals[len(s.pats)] {
			return c, false
		}
	}

	if c == nil || dryrun || len(s.pats) == 0 {
		return c, true
	}

	return withURLParams(c, matches), true
}

func (s stringPattern) String() string {
	return fmt.Sprintf("stringPattern(%q, %v)", s.raw, s.isPrefix)
}

var patternRe = regexp.MustCompile(`/:([^/]+)`)

func parseStringPattern(s string) stringPattern {
	var isPrefix bool
	// Routes that end in an asterisk ("*") are prefix routes
	if len(s) > 0 && s[len(s)-1] == '*' {
		s = s[:len(s)-1]
		isPrefix = true
	}

	matches := patternRe.FindAllStringSubmatchIndex(s, -1)
	pats := make([]string, len(matches))
	literals := make([]string, len(matches)+1)
	n := 0
	for i, match := range matches {
		a, b := match[2], match[3]
		literals[i] = s[n : a-1] // Need to leave off the colon
		pats[i] = s[a:b]
		n = b
	}
	literals[len(matches)] = s[n:]
	return stringPattern{
		raw:      s,
		pats:     pats,
		literals: literals,
		isPrefix: isPrefix,
	}
}
