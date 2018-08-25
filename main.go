package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"strings"

	"github.com/Quasilyte/go-consistent/internal/typeof"
)

func main() {
	log.SetFlags(0)

	flag.Parse()

	filenames := targetsToFilenames(flag.Args())

	var ctxt context
	ctxt.SetupOpsTable()
	if err := visitFiles(&ctxt, filenames, ctxt.InferConventions); err != nil {
		log.Fatalf("infer conventions: %v", err)
	}
	ctxt.SetupSuggestions()
	if err := visitFiles(&ctxt, filenames, ctxt.ReportInconsistent); err != nil {
		log.Fatalf("report inconsistent: %v", err)
	}
}

type context struct {
	ops  []*operation
	fset *token.FileSet
}

type operation struct {
	name     string
	suggest  *opVariant
	variants []*opVariant
}

type opVariant struct {
	name    string
	count   int
	matcher opMatcher
}

type opMatcher interface {
	Skip(ast.Node) bool
	Match(ast.Node) bool
}

type matcherBase struct{}

func (matcherBase) Skip(ast.Node) bool {
	return false
}

type newMatcher struct{ matcherBase }

func (newMatcher) Match(n ast.Node) bool {
	e, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}
	fn, ok := e.Fun.(*ast.Ident)
	return ok && fn.Name == "new"
}

type addressOfLitMatcher struct{ matcherBase }

func (addressOfLitMatcher) Match(n ast.Node) bool {
	e, ok := n.(*ast.UnaryExpr)
	if !ok {
		return false
	}
	_, ok = e.X.(*ast.CompositeLit)
	return ok
}

type emptySliceMakeMatcher struct{ matcherBase }

func (emptySliceMakeMatcher) Match(n ast.Node) bool {
	e, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}
	if len(e.Args) != 2 {
		return false // Requires {T, len} arguments
	}
	fn, ok := e.Fun.(*ast.Ident)
	return ok && fn.Name == "make" &&
		typeof.IsSlice(e.Args[0]) &&
		valueOf(e.Args[1]) == "0"
}

type emptySliceLitMatcher struct{ matcherBase }

func (emptySliceLitMatcher) Skip(n ast.Node) bool {
	// Don't consider slice literals like &T{}.
	e, ok := n.(*ast.UnaryExpr)
	return ok && e.Op == token.AND
}

func (emptySliceLitMatcher) Match(n ast.Node) bool {
	e, ok := n.(*ast.CompositeLit)
	if !ok {
		return false
	}
	return typeof.IsSlice(e.Type) && len(e.Elts) == 0
}

type emptyMapMakeMatcher struct{ matcherBase }

func (emptyMapMakeMatcher) Match(n ast.Node) bool {
	e, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}
	fn, ok := e.Fun.(*ast.Ident)
	return ok && fn.Name == "make" &&
		typeof.IsMap(e.Args[0]) &&
		(len(e.Args) == 1 || len(e.Args) == 2 && valueOf(e.Args[1]) == "0")
}

type emptyMapLitMatcher struct{ matcherBase }

func (emptyMapLitMatcher) Match(n ast.Node) bool {
	e, ok := n.(*ast.CompositeLit)
	if !ok {
		return false
	}
	return typeof.IsMap(e.Type) && len(e.Elts) == 0
}

func (ctxt *context) SetupOpsTable() {
	ctxt.ops = []*operation{
		{
			name: "zero value pointer allocation",
			variants: []*opVariant{
				{name: "new", matcher: newMatcher{}},
				{name: "address-of-lit", matcher: addressOfLitMatcher{}},
			},
		},

		{
			name: "empty slice",
			variants: []*opVariant{
				{name: "empty-slice-make", matcher: emptySliceMakeMatcher{}},
				{name: "empty-slice-lit", matcher: emptySliceLitMatcher{}},
			},
		},

		{
			name: "empty map",
			variants: []*opVariant{
				{name: "empty-map-make", matcher: emptyMapMakeMatcher{}},
				{name: "empty-map-lit", matcher: emptyMapLitMatcher{}},
			},
		},

		// TODO(quasilyte): nil map
		// TODO(quasilyte): nil slice
	}
}

func (ctxt *context) SetupSuggestions() {
	for _, op := range ctxt.ops {
		op.suggest = op.variants[0]
		// Find the most frequently used variant.
		for _, v := range op.variants[1:] {
			if v.count > op.suggest.count {
				op.suggest = v
			}
		}
		// Diagnostic: check if there were multiple candidates.
		if op.suggest.count == 0 {
			continue
		}
		for _, v := range op.variants {
			if v != op.suggest && v.count == op.suggest.count {
				log.Printf("warning: %s: can't decide between %s and %s",
					op.name, v.name, op.suggest.name)
			}
		}
	}
}

func (ctxt *context) InferConventions(f *ast.File) {
	for _, op := range ctxt.ops {
		for _, v := range op.variants {
			ast.Inspect(f, func(n ast.Node) bool {
				if n == nil {
					return false
				}
				if v.matcher.Skip(n) {
					return false
				}
				if v.matcher.Match(n) {
					v.count++
				}
				return true
			})
		}
	}
}

func (ctxt *context) ReportInconsistent(f *ast.File) {
	for _, op := range ctxt.ops {
		for _, v := range op.variants {
			ast.Inspect(f, func(n ast.Node) bool {
				if n == nil {
					return false
				}
				if v.matcher.Skip(n) {
					return false
				}
				if v.matcher.Match(n) && v != op.suggest {
					ctxt.printWarning(n, op, v)
				}
				return true
			})
		}
	}
}

func (ctxt *context) printWarning(cause ast.Node, op *operation, bad *opVariant) {
	// TODO(quasilyte): figure out a better message format.
	pos := ctxt.fset.Position(cause.Pos())
	log.Printf("%s: %s: use %s instead of %s",
		pos, op.name, op.suggest.name, bad.name)
}

func visitFiles(ctxt *context, filenames []string, visit func(*ast.File)) error {
	fset := token.NewFileSet()
	ctxt.fset = fset
	for _, filename := range filenames {
		f, err := parser.ParseFile(fset, filename, nil, 0)
		if err != nil {
			return err
		}
		visit(f)
	}
	return nil
}

func targetsToFilenames(targets []string) []string {
	var filenames []string

	for _, target := range targets {
		if !strings.HasSuffix(target, ".go") {
			// TODO(quasilyte): add package targets support.
			log.Printf("skip target %q: not a Go file", target)
			continue
		}
		abs, err := filepath.Abs(target)
		if err != nil {
			log.Printf("skip target %q: %v", err)
			continue
		}
		filenames = append(filenames, abs)
	}

	return filenames
}

func valueOf(x ast.Expr) string {
	switch x := x.(type) {
	case *ast.BasicLit:
		return x.Value
	default:
		return ""
	}
}
