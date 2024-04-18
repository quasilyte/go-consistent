package end2end

import (
	"fmt"
	"go/token"
	"regexp"
	"strings"
)

func (p *TestParser) parseFile(filename string, data []byte) (*TestFile, error) {
	matchers := make(map[int][]*Matcher)
	var pending []*Matcher

	f := &TestFile{Matchers: matchers, text: string(data)}

	for i, l := range strings.Split(f.text, "\n") {
		line := i + 1
		op, text, ok := p.fetchMatcher(l)
		switch {
		case ok:
			m := Matcher{text: text, pos: token.Position{Filename: filename}}
			switch op {
			case "~":
				re, err := regexp.Compile(text)
				if err != nil {
					return nil, fmt.Errorf("%s:%d: %w", filename, line, err)
				}
				m.re = re
			case "=":
				// Do nothing.
			default:
				return nil, fmt.Errorf("%s:%d: unknown op %q", filename, line, op)
			}
			pending = append(pending, &m)

		case len(pending) != 0:
			for _, m := range pending {
				m.pos.Line = line
			}
			// Copy all matchers from the pending list to the result.
			matchers[line] = append([]*Matcher{}, pending...)
			pending = pending[:0] // Clear pending list
		}
	}

	return f, nil
}

func (p *TestParser) fetchMatcher(l string) (op, text string, ok bool) {
	re := defaultMatcherRE
	if p.MatcherRE != nil {
		re = p.MatcherRE
	}

	m := re.FindStringSubmatch(l)
	if len(m) < 2 {
		return "", "", false
	}
	return m[1], m[2], true
}
