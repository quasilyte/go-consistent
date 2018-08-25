package typeof

import "go/ast"

func IsSlice(x ast.Node) bool {
	return inferApproxType(x) == typSlice
}

func IsMap(x ast.Node) bool {
	return inferApproxType(x) == typMap
}

func inferApproxType(x ast.Node) approxType {
	switch x := x.(type) {
	case *ast.Ident:
		return typeOfObject(x.Obj)

	case *ast.ArrayType:
		if x.Len == nil {
			return typSlice
		}
		return typArray

	case *ast.MapType:
		return typMap

	default:
		return typUnknown
	}
}

func typeOfObject(obj *ast.Object) approxType {
	if obj.Decl == nil {
		return typUnknown
	}
	spec, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		return typUnknown
	}
	return inferApproxType(spec.Type)
}

type approxType int

const (
	typUnknown approxType = iota
	typArray
	typSlice
	typMap
)
