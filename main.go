package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"reflect"
	"strings"
)

func main() {
	log.SetFlags(0)

	var ctxt context

	flag.BoolVar(&ctxt.pedantic, "pedantic", false,
		`makes several diagnostics more pedantic and comprehensive`)
	flag.Parse()

	filenames := targetsToFilenames(flag.Args())

	ctxt.Init()
	if err := visitFiles(&ctxt, filenames, ctxt.InferConventions); err != nil {
		log.Fatalf("infer conventions: %v", err)
	}
	ctxt.SetupSuggestions()
	if err := visitFiles(&ctxt, filenames, ctxt.CaptureInconsistencies); err != nil {
		log.Fatalf("report inconsistent: %v", err)
	}

	for _, warn := range ctxt.warnings {
		log.Printf("%s: %s", warn.pos, warn.text)
	}
}

type context struct {
	fset *token.FileSet

	ops []*operation

	pedantic bool

	warnings []warning
}

type warning struct {
	pos  token.Position
	text string
}

type operationPrototype interface {
	New() *operation
}

type opScope int

const (
	scopeAny opScope = iota
	scopeLocal
	scopeGlobal
)

type operation struct {
	// name is a human-readable operation descriptor.
	// Initialized by a context init.
	name string

	// suggested is an op variant that is inferred as the most frequently used one.
	suggested *opVariant

	// scope determines context in which operation is checked.
	// Initialized by prototype.
	scope opScope

	// variants is a list of equivalent operation forms.
	// Initialized by prototype.
	variants []*opVariant
}

type opVariant struct {
	// name is a human-readable operation variant descriptor.
	// Initialized by prototype.
	name string

	// skip is a function that can reject recursing into node siblings
	// during AST traversal.
	// Initialized by prototype.
	// Can be nil.
	skip func(ast.Node) bool

	// match reports whether given node represents an action described by
	// this op variant.
	// Initialized by prototype. Can be re-assigned in context init.
	match func(ast.Node) bool

	// matchPedantic is an optional pedantic match variant.
	// Initialized by prototype.
	// If not nil, used instead of normal match when -pedantic=true flag is provided.
	matchPedantic func(ast.Node) bool

	// count is a counter for op variant usages.
	// Updated during the first AST traversal.
	count int
}

func (ctxt *context) Init() {
	prototypes := []operationPrototype{
		zeroValPtrAllocProto{},
		emptySliceProto{},
		nilSliceDeclProto{},
		emptyMapProto{},
		nilMapDeclProto{},
		hexLitProto{},
		rangeCheckProto{},
		andNotProto{},
	}

	for _, proto := range prototypes {
		rv := reflect.ValueOf(proto)
		typ := rv.Type()
		if !strings.HasSuffix(typ.Name(), "Proto") {
			panic(fmt.Sprintf("%s: missing Proto type name suffix", typ.Name()))
		}
		op := proto.New()
		op.name = typ.Name()[:len(typ.Name())-len("Proto")]

		for _, v := range op.variants {
			if ctxt.pedantic && v.matchPedantic != nil {
				v.match = v.matchPedantic
			}
		}

		ctxt.ops = append(ctxt.ops, op)
	}
}

func (ctxt *context) SetupSuggestions() {
	for _, op := range ctxt.ops {
		op.suggested = op.variants[0]
		// Find the most frequently used variant.
		for _, v := range op.variants[1:] {
			if v.count > op.suggested.count {
				op.suggested = v
			}
		}
		// Diagnostic: check if there were multiple candidates.
		if op.suggested.count == 0 {
			continue
		}
		for _, v := range op.variants {
			if op.suggested != v && v.count == op.suggested.count {
				log.Printf("warning: %s: can't decide between %s and %s",
					op.name, v.name, op.suggested.name)
			}
		}
	}
}

type opVisitFunc func(*operation, *opVariant, ast.Node) bool

func (ctxt *context) visitOps(f *ast.File, visit opVisitFunc) {
	for _, op := range ctxt.ops {
		switch op.scope {
		case scopeAny:
			for _, v := range op.variants {
				ast.Inspect(f, func(n ast.Node) bool {
					if n == nil {
						return false
					}
					return visit(op, v, n)
				})
			}

		case scopeLocal:
			for _, v := range op.variants {
				for _, decl := range f.Decls {
					decl, ok := decl.(*ast.FuncDecl)
					if !ok || decl.Body == nil {
						continue
					}
					ast.Inspect(decl.Body, func(n ast.Node) bool {
						if n == nil {
							return false
						}
						return visit(op, v, n)
					})
				}
			}

		case scopeGlobal:
			// TODO(quasilyte): remove later if never used.
			panic("unimplemented and unused")

		default:
			panic(fmt.Sprintf("unexpected scope: %d", op.scope))
		}
	}
}

func (ctxt *context) InferConventions(f *ast.File) {
	ctxt.visitOps(f, func(op *operation, v *opVariant, n ast.Node) bool {
		if n == nil {
			return false
		}
		if v.skip != nil && v.skip(n) {
			return false
		}
		if v.match(n) {
			v.count++
		}
		return true
	})
}

func (ctxt *context) CaptureInconsistencies(f *ast.File) {
	ctxt.visitOps(f, func(op *operation, v *opVariant, n ast.Node) bool {
		if n == nil {
			return false
		}
		if v.skip != nil && v.skip(n) {
			return false
		}
		if v.match(n) && v != op.suggested {
			ctxt.pushWarning(n, op, v)
		}
		return true
	})
}

func (ctxt *context) pushWarning(cause ast.Node, op *operation, bad *opVariant) {
	pos := ctxt.fset.Position(cause.Pos())
	text := fmt.Sprintf("%s: use %s instead of %s", op.name, op.suggested.name, bad.name)
	ctxt.warnings = append(ctxt.warnings, warning{pos: pos, text: text})
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
			log.Printf("skip target %q: %v", target, err)
			continue
		}
		filenames = append(filenames, abs)
	}

	return filenames
}

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
