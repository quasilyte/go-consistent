package main

import (
	"go/ast"
)

func valueOf(x ast.Node) string {
	switch x := x.(type) {
	case *ast.BasicLit:
		return x.Value
	case *ast.Ident:
		return x.Name
	default:
		return ""
	}
}

var (
	sentinelBinaryExpr = &ast.BinaryExpr{}
	sentinelUnaryExpr  = &ast.UnaryExpr{}
	sentinelCallExpr   = &ast.CallExpr{}
	sentinelIdent      = &ast.Ident{}
	sentinelSliceExpr  = &ast.SliceExpr{}
	sentinelGenDecl    = &ast.GenDecl{}
)

func asBinaryExpr(n ast.Node) *ast.BinaryExpr {
	if e, ok := n.(*ast.BinaryExpr); ok {
		return e
	}
	return sentinelBinaryExpr
}

func asUnaryExpr(n ast.Node) *ast.UnaryExpr {
	if e, ok := n.(*ast.UnaryExpr); ok {
		return e
	}
	return sentinelUnaryExpr
}

func asCallExpr(n ast.Node) *ast.CallExpr {
	if e, ok := n.(*ast.CallExpr); ok {
		return e
	}
	return sentinelCallExpr
}

func asIdent(n ast.Node) *ast.Ident {
	if e, ok := n.(*ast.Ident); ok {
		return e
	}
	return sentinelIdent
}

func asSliceExpr(n ast.Node) *ast.SliceExpr {
	if e, ok := n.(*ast.SliceExpr); ok {
		return e
	}
	return sentinelSliceExpr
}

func asGenDecl(n ast.Node) *ast.GenDecl {
	if decl, ok := n.(*ast.GenDecl); ok {
		return decl
	}
	return sentinelGenDecl
}

func isNil(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.BinaryExpr:
		return n == nil || n == sentinelBinaryExpr
	case *ast.UnaryExpr:
		return n == nil || n == sentinelUnaryExpr
	case *ast.CallExpr:
		return n == nil || n == sentinelCallExpr
	case *ast.Ident:
		return n == nil || n == sentinelIdent
	case *ast.SliceExpr:
		return n == nil || n == sentinelSliceExpr
	case *ast.GenDecl:
		return n == nil || n == sentinelGenDecl
	default:
		return n == nil
	}
}
