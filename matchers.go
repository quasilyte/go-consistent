package main

import (
	"go/ast"
	"go/token"

	"github.com/Quasilyte/go-consistent/internal/typeof"
)

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

type nilSliceVarMatcher struct{ matcherBase }

func (nilSliceVarMatcher) Match(n ast.Node) bool {
	// TODO(quasilyte): can this be simplified with type assertion to GenDecl?
	s, ok := n.(*ast.DeclStmt)
	if !ok {
		return false
	}
	d := s.Decl.(*ast.GenDecl)
	if d.Tok != token.VAR {
		return false
	}
	// TODO(quasilyte): handle multi-spec var decls.
	if len(d.Specs) != 1 {
		return false
	}
	spec := d.Specs[0].(*ast.ValueSpec)
	// TODO(quasilyte): handle multi-name var decls.
	if len(spec.Names) != 1 {
		return false
	}
	if spec.Values != nil {
		return false
	}
	return typeof.IsSlice(spec.Type)
}

type nilSliceLitMatcher struct{ matcherBase }

func (nilSliceLitMatcher) Match(n ast.Node) bool {
	e, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}
	return len(e.Args) == 1 &&
		valueOf(e.Args[0]) == "nil" &&
		typeof.IsSlice(e.Fun)
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
