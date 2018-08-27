package main

import (
	"log"
	"path"
	"testing"

	"github.com/Quasilyte/go-consistent/internal/end2end"
)

func TestEnd2End(t *testing.T) {
	filenames := []string{
		"positive_tests1.go",
		"positive_tests2.go",
		"negative_tests1.go",
		"negative_tests2.go",
	}

	for _, filename := range filenames {
		t.Run(filename, func(t *testing.T) {
			rel := path.Join("testdata", filename)
			f, err := end2end.ParseTestFile(rel)
			if err != nil {
				t.Fatalf("parse %s: %v", rel, err)
			}

			var ctxt context
			ctxt.Init()
			if err := visitFiles(&ctxt, []string{rel}, ctxt.InferConventions); err != nil {
				log.Fatalf("infer conventions: %v", err)
			}
			ctxt.SetupSuggestions()
			if err := visitFiles(&ctxt, []string{rel}, ctxt.CaptureInconsistencies); err != nil {
				log.Fatalf("report inconsistent: %v", err)
			}

			for _, warn := range ctxt.warnings {
				mlist, ok := f.Matchers[warn.pos.Line]
				if !ok {
					t.Errorf("%s: unexpected warning: %s", warn.pos, warn.text)
					continue
				}

				for _, m := range mlist {
					if m.Match(warn.text) {
						m.Matches++
						break
					}
				}
			}

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
