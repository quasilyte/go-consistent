package main

import (
	"go/ast"
)

var (
	sentinelBinaryExpr = &ast.BinaryExpr{}
)

func asBinaryExpr(n ast.Node) *ast.BinaryExpr {
	if e, ok := n.(*ast.BinaryExpr); ok {
		return e
	}
	return sentinelBinaryExpr
}

func isNil(n ast.Node) bool {
	switch n := n.(type) {
	case *ast.BinaryExpr:
		return n == nil || n == sentinelBinaryExpr
	default:
		return n == nil
	}
}
