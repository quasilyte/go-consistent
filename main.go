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

	flag.BoolVar(&ctxt.Pedantic, "pedantic", false,
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

	for _, warn := range ctxt.Warnings {
		log.Printf("%s: %s", warn.pos, warn.text)
	}
}

type context struct {
	fset *token.FileSet

	ops []*operation

	Pedantic bool

	Warnings []warning
}

type warning struct {
	pos  token.Position
	text string
}

type operationPrototype interface {
	Variants() []opVariantPrototype
}

type opVariantPrototype struct {
	name          string
	skip          func(ast.Node) bool
	match         func(ast.Node) bool
	matchPedantic func(ast.Node) bool
}

type opScope int

const (
	scopeAny opScope = iota
	scopeLocal
	scopeGlobal
)

type operation struct {
	name      string
	scope     opScope
	suggested *opVariant
	variants  []*opVariant
}

type opVariant struct {
	name  string
	skip  func(ast.Node) bool
	match func(ast.Node) bool

	count int
}

func (ctxt *context) Init() {
	prototypes := []operationPrototype{
		zeroValPtrAllocProto{},
		emptySliceProto{},
		nilSliceDeclProto{},
		emptyMapProto{},
	}

	for _, proto := range prototypes {
		rv := reflect.ValueOf(proto)
		typ := rv.Type()
		if !strings.HasSuffix(typ.Name(), "Proto") {
			panic(fmt.Sprintf("%s: missing Proto type name suffix", typ.Name()))
		}
		op := &operation{
			name: typ.Name()[:len(typ.Name())-len("Proto")],
		}

		// Parse metadata fields, op attributes.
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if field.Type.Name() != "opAttr" {
				continue
			}

			switch field.Name {
			case "scopeAny":
				op.scope = scopeAny
			case "scopeLocal":
				op.scope = scopeLocal
			case "scopeGlobal":
				op.scope = scopeGlobal
			default:
				panic(fmt.Sprintf("%s: unexpected opAttr: %s",
					op.name, field.Name))
			}
		}

		// Instantiate op variants.
		for _, vproto := range proto.Variants() {
			match := vproto.match
			if ctxt.Pedantic && vproto.matchPedantic != nil {
				match = vproto.matchPedantic
			}
			skip := vproto.skip
			if skip == nil {
				skip = func(ast.Node) bool { return false }
			}
			op.variants = append(op.variants, &opVariant{
				name:  vproto.name,
				skip:  skip,
				match: match,
			})
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
					return visit(op, v, n)
				})
			}

		case scopeLocal:
			for _, v := range op.variants {
				for _, decl := range f.Decls {
					decl, ok := decl.(*ast.FuncDecl)
					if !ok {
						continue
					}
					ast.Inspect(decl.Body, func(n ast.Node) bool {
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
		if v.skip(n) {
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
		if v.skip(n) {
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
	ctxt.Warnings = append(ctxt.Warnings, warning{pos: pos, text: text})
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

func valueOf(x ast.Expr) string {
	switch x := x.(type) {
	case *ast.BasicLit:
		return x.Value
	case *ast.Ident:
		return x.Name
	default:
		return ""
	}
}
