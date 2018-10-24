package main

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/Quasilyte/go-consistent/internal/typeof"
	"github.com/go-toolsmith/astequal"
)

type unitImportProto struct{}

func (p unitImportProto) New() *operation {
	return &operation{
		scope: scopeGlobal,
		variants: []*opVariant{
			{text: "omit parenthesis in single-package import", match: p.matchNoParens},
			{text: "wrap single-package import spec into parenthesis", match: p.matchWithParens},
		},
	}
}

func (p unitImportProto) matchNoParens(n ast.Node) bool {
	decl := asGenDecl(n)
	return decl.Tok == token.IMPORT && len(decl.Specs) == 1 &&
		decl.Lparen == 0 && decl.Rparen == 0
}

func (p unitImportProto) matchWithParens(n ast.Node) bool {
	decl := asGenDecl(n)
	return decl.Tok == token.IMPORT && len(decl.Specs) == 1 &&
		decl.Lparen != 0 && decl.Rparen != 0
}

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

type nilMapDeclProto struct{}

func (p nilMapDeclProto) New() *operation {
	return &operation{
		scope: scopeLocal,
		variants: []*opVariant{
			{name: "var", match: p.matchVar},
			{name: "literal", match: p.matchLiteral},
		},
	}
}

func (nilMapDeclProto) matchVar(n ast.Node) bool {
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
	return spec.Values == nil && typeof.IsMap(spec.Type)
}

func (nilMapDeclProto) matchLiteral(n ast.Node) bool {
	assign, ok := n.(*ast.AssignStmt)
	if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 {
		return false
	}
	e, ok := assign.Rhs[0].(*ast.CallExpr)
	return ok && assign.Tok == token.DEFINE &&
		len(e.Args) == 1 &&
		valueOf(e.Args[0]) == "nil" &&
		typeof.IsMap(e.Fun)
}

type hexLitProto struct{}

func (p hexLitProto) New() *operation {
	return &operation{
		scope: scopeAny,
		variants: []*opVariant{
			{name: "a-f", match: p.matchLowercase},
			{name: "A-F", match: p.matchUppercase},
		},
	}
}

func (hexLitProto) matchLowercase(n ast.Node) bool {
	v := valueOf(n)
	return strings.HasPrefix(v, "0x") && strings.ContainsAny(v, "abcdef")
}

func (hexLitProto) matchUppercase(n ast.Node) bool {
	v := valueOf(n)
	return strings.HasPrefix(v, "0x") && strings.ContainsAny(v, "ABCDEF")
}

type rangeCheckProto struct{}

func (p rangeCheckProto) New() *operation {
	return &operation{
		scope: scopeAny,
		variants: []*opVariant{
			{name: "align-left", match: p.matchAlignLeft},
			{name: "align-center", match: p.matchAlignCenter},
		},
	}
}

func (rangeCheckProto) matchAlignLeft(n ast.Node) bool {
	e := asBinaryExpr(n)
	lhs := asBinaryExpr(e.X)
	rhs := asBinaryExpr(e.Y)
	return !isNil(e) && !isNil(lhs) && !isNil(rhs) &&
		(e.Op == token.LAND || e.Op == token.LOR) &&
		(lhs.Op == token.GTR || lhs.Op == token.GEQ) &&
		(rhs.Op == token.LSS || rhs.Op == token.LEQ) &&
		astequal.Expr(lhs.X, rhs.X)
}

func (rangeCheckProto) matchAlignCenter(n ast.Node) bool {
	e := asBinaryExpr(n)
	lhs := asBinaryExpr(e.X)
	rhs := asBinaryExpr(e.Y)
	return !isNil(e) && !isNil(lhs) && !isNil(rhs) &&
		(e.Op == token.LAND || e.Op == token.LOR) &&
		(lhs.Op == token.LSS || lhs.Op == token.LEQ) &&
		(rhs.Op == token.LSS || lhs.Op == token.LEQ) &&
		astequal.Expr(lhs.Y, rhs.X)
}

type andNotProto struct{}

func (p andNotProto) New() *operation {
	return &operation{
		scope: scopeAny,
		variants: []*opVariant{
			{name: "&^", match: p.matchSingleTok},
			{name: "&-plus-^", match: p.matchTwoTok},
		},
	}
}

func (andNotProto) matchSingleTok(n ast.Node) bool {
	e, ok := n.(*ast.BinaryExpr)
	return ok && e.Op == token.AND_NOT
}

func (andNotProto) matchTwoTok(n ast.Node) bool {
	e := asBinaryExpr(n)
	rhs := asUnaryExpr(e.Y)
	return !isNil(e) && !isNil(rhs) &&
		e.Op == token.AND && rhs.Op == token.XOR
}

type floatLitProto struct{}

func (p floatLitProto) New() *operation {
	return &operation{
		scope: scopeAny,
		variants: []*opVariant{
			{name: "explicit-int/frac", match: p.matchExplicit},
			{name: "omitted-int/frac", match: p.matchOmitted},
		},
	}
}

func (floatLitProto) splitIntFrac(n *ast.BasicLit) (integer, frac string) {
	parts := strings.Split(n.Value, ".")
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func (p floatLitProto) matchExplicit(n ast.Node) bool {
	lit, ok := n.(*ast.BasicLit)
	if !ok || lit.Kind != token.FLOAT {
		return false
	}
	integer, frac := p.splitIntFrac(lit)
	return (integer == "0" && frac != "") ||
		(integer != "" && frac == "0")
}

func (p floatLitProto) matchOmitted(n ast.Node) bool {
	lit, ok := n.(*ast.BasicLit)
	if !ok || lit.Kind != token.FLOAT {
		return false
	}
	integer, frac := p.splitIntFrac(lit)
	return (frac == "" && integer != "") ||
		(frac != "" && integer == "")
}

type labelCaseProto struct {
	allUpperCase   *regexp.Regexp
	upperCamelCase *regexp.Regexp
	lowerCamelCase *regexp.Regexp
}

func (p labelCaseProto) New() *operation {
	allUpperCase := regexp.MustCompile(`^[A-Z][A-Z_]`)
	upperCamelCase := regexp.MustCompile(`^[A-Z][a-z]`)
	lowerCamelCase := regexp.MustCompile(`^[a-z]`)
	match := func(rx *regexp.Regexp) func(n ast.Node) bool {
		return func(n ast.Node) bool {
			stmt, ok := n.(*ast.LabeledStmt)
			return ok && rx.MatchString(stmt.Label.Name)
		}
	}
	return &operation{
		scope: scopeLocal,
		variants: []*opVariant{
			{name: "ALL_UPPER", match: match(allUpperCase)},
			{name: "UpperCamelCase", match: match(upperCamelCase)},
			{name: "lowerCamelCase", match: match(lowerCamelCase)},
		},
	}
}
