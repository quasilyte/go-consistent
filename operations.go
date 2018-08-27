package main

import (
	"go/ast"
	"go/token"

	"github.com/Quasilyte/go-consistent/internal/typeof"
)

type zeroValPtrAllocProto struct{}

func (p zeroValPtrAllocProto) New() *operation {
	return &operation{
		scope: scopeAny,
		variants: []*opVariant{
			{name: "new", match: p.matchNew},
			{name: "address-of-literal", match: p.matchAddressOfLiteral},
		},
	}
}

func (zeroValPtrAllocProto) matchNew(n ast.Node) bool {
	e, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}
	fn, ok := e.Fun.(*ast.Ident)
	return ok && fn.Name == "new"
}

func (zeroValPtrAllocProto) matchAddressOfLiteral(n ast.Node) bool {
	e, ok := n.(*ast.UnaryExpr)
	if !ok {
		return false
	}
	_, ok = e.X.(*ast.CompositeLit)
	return ok
}

type emptySliceProto struct{}

func (p emptySliceProto) New() *operation {
	return &operation{
		scope: scopeAny,
		variants: []*opVariant{
			{name: "make", match: p.matchMake},
			{name: "literal", skip: p.skipLiteral, match: p.matchLiteral},
		},
	}
}

func (emptySliceProto) matchMake(n ast.Node) bool {
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

func (emptySliceProto) skipLiteral(n ast.Node) bool {
	// Don't consider slice literals like &T{}.
	e, ok := n.(*ast.UnaryExpr)
	return ok && e.Op == token.AND
}

func (emptySliceProto) matchLiteral(n ast.Node) bool {
	e, ok := n.(*ast.CompositeLit)
	if !ok {
		return false
	}
	return typeof.IsSlice(e.Type) && len(e.Elts) == 0
}

type nilSliceDeclProto struct{}

func (p nilSliceDeclProto) New() *operation {
	return &operation{
		scope: scopeLocal,
		variants: []*opVariant{
			{name: "var", match: p.matchVar},
			{name: "literal", match: p.matchLiteral},
		},
	}
}

func (nilSliceDeclProto) matchVar(n ast.Node) bool {
	d, ok := n.(*ast.GenDecl)
	if !ok || d.Tok != token.VAR {
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
	return spec.Values == nil && typeof.IsSlice(spec.Type)
}

func (nilSliceDeclProto) matchLiteral(n ast.Node) bool {
	assign, ok := n.(*ast.AssignStmt)
	if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return false
	}
	e, ok := assign.Rhs[0].(*ast.CallExpr)
	return ok && assign.Tok == token.DEFINE &&
		len(e.Args) == 1 &&
		valueOf(e.Args[0]) == "nil" &&
		typeof.IsSlice(e.Fun)
}

type emptyMapProto struct{}

func (p emptyMapProto) New() *operation {
	return &operation{
		variants: []*opVariant{
			{name: "make", match: p.matchMake},
			{name: "literal", match: p.matchLiteral},
		},
	}
}

func (emptyMapProto) matchMake(n ast.Node) bool {
	e, ok := n.(*ast.CallExpr)
	if !ok {
		return false
	}
	fn, ok := e.Fun.(*ast.Ident)
	return ok && fn.Name == "make" &&
		typeof.IsMap(e.Args[0]) &&
		(len(e.Args) == 1 || len(e.Args) == 2 && valueOf(e.Args[1]) == "0")
}

func (emptyMapProto) matchLiteral(n ast.Node) bool {
	e, ok := n.(*ast.CompositeLit)
	if !ok {
		return false
	}
	return typeof.IsMap(e.Type) && len(e.Elts) == 0
}
