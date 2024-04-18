package end2end

import (
	"fmt"
	"go/token"
	"io"
	"os"
	"regexp"
)

// ParseTestFile parses file at specified path using the parser with default settings.
func ParseTestFile(filename string) (*TestFile, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read test file: %w", err)
	}
	var p TestParser
	return p.parseFile(filename, data)
}

// TestFile represents parsed end2end test file.
type TestFile struct {
	// Matchers is a mapping from source code line to the
	// list of matchers for it.
	Matchers map[int][]*Matcher

	text string
}

// Text returns source code file contents that were used to create f.
func (f *TestFile) Text() string { return f.text }

// Matcher is a single parsed magic comment that is used
// to match tested program output.
type Matcher struct {
	// Matches is a counter of matches for this matcher.
	// Expected to be set by the user.
	Matches int

	text string
	pos  token.Position
	re   *regexp.Regexp // nil for "literal" (non-regexp) matchers
}

// MatchWithLine does text+position matching.
//
// Like Match, but also checks that m.Position().Line is
// equal to the provided line argument.
func (m *Matcher) MatchWithLine(s string, line int) bool {
	return m.Match(s) && m.pos.Line == line
}

// Match tries to match s string against matcher pattern.
func (m *Matcher) Match(s string) bool {
	if m.re != nil {
		return m.re.MatchString(s)
	}
	return m.text == s
}

// IsMatched reports whether m is matched at least once.
func (m *Matcher) IsMatched() bool { return m.Matches != 0 }

// Text returns the matcher comment text.
func (m *Matcher) Text() string { return m.text }

// Position returns matcher comment text position.
func (m *Matcher) Position() token.Position { return m.pos }

// TestParser is a end2end source file parser.
type TestParser struct {
	// MatcherRE is used to distinguish comments that describe matchers
	// and also to capture their text.
	//
	// Regexp must include 2 capture groups:
	//	1. Matches [=~]. This determines the matcher kind (text or regexp)
	//	2. Matches text being matched. Usually something like ".*"
	//
	// If nil, defaultMatcherRE is used. You can use it as an example.
	MatcherRE *regexp.Regexp
}

// ParseFile parses everything from r using filename to make associations.
func (p *TestParser) ParseFile(filename string, r io.Reader) (*TestFile, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read test file: %w", err)
	}
	return p.parseFile(filename, data)
}

var defaultMatcherRE = regexp.MustCompile(`^\s*//([=~]) (.*)`)
