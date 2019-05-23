package main

import (
	"go/token"
	"path"
	"testing"

	"github.com/quasilyte/go-consistent/internal/end2end"
)

func TestEnd2End(t *testing.T) {
	filenames := []string{
		"positive_tests1.go",
		"positive_tests2.go",
		"positive_tests3.go",
		"negative_tests1.go",
		"negative_tests2.go",
		"negative_tests3.go",
		"negative_tests4.go",
	}

	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			rel := path.Join("testdata", filename)
			f, err := end2end.ParseTestFile(rel)
			if err != nil {
				t.Fatalf("parse %s: %v", rel, err)
			}

			var ctxt context
			ctxt.paths = []string{rel}
			ctxt.initCheckers()
			if err := ctxt.collectAllCandidates(); err != nil {
				t.Fatalf("collect candidates: %v", err)
			}
			ctxt.assignSuggestions()
			visitWarnings(&ctxt, func(pos token.Position, v *opVariant) {
				text := v.op.name + ": " + v.op.suggested.warning
				mlist, ok := f.Matchers[pos.Line]
				if !ok {
					t.Errorf("%s: unexpected warning: %s", pos, text)
					return
				}

				for _, m := range mlist {
					if m.Match(text) {
						m.Matches++
						break
					} else {
						t.Errorf("%s: unexpected warning: %s", m.Position(), text)
					}
				}
			})

			for _, mlist := range f.Matchers {
				for _, m := range mlist {
					if !m.IsMatched() {
						t.Errorf("%s: no matches: %s", m.Position(), m.Text())
					}
				}
			}
		})
	}
}
