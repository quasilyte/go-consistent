package main

import (
	"go/ast"
)

var (
	sentinelBinaryExpr = &ast.BinaryExpr{}
	sentinelUnaryExpr  = &ast.UnaryExpr{}
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

func isNil(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.BinaryExpr:
		return n == nil || n == sentinelBinaryExpr
	case *ast.UnaryExpr:
		return n == nil || n == sentinelUnaryExpr
	default:
		return n == nil
	}
}
