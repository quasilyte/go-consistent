package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-toolsmith/astinfo"
	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/packages"
)

var generatedFileCommentRE = regexp.MustCompile("Code generated .* DO NOT EDIT.")

func main() {
	log.SetFlags(0)
	var ctxt context

	steps := []struct {
		name string
		fn   func() error
	}{
		{"parse flags", ctxt.parseFlags},
		{"resolve targets", ctxt.resolveTargets},
		{"init checkers", ctxt.initCheckers},
		{"collect candidates", ctxt.collectAllCandidates},
		{"assign suggestions", ctxt.assignSuggestions},
		{"print warnings", ctxt.printWarnings},
	}

	for _, step := range steps {
		if err := step.fn(); err != nil {
			log.Fatalf("%s: %v", step.name, err)
		}
	}
}

type context struct {
	// flags is an (effectively) immutable struct that holds all command-line
	// arguments as they were passed to the program.
	//
	// For per-argument documentation see context.parseFlags.
	flags struct {
		pedantic bool
		verbose  bool
		debug    bool
		targets  []string
		exclude  string
	}

	paths []string

	locs *locationMap

	fset    *token.FileSet
	info    *types.Info
	astinfo astinfo.Info

	checkers []checker

	candidates []candidate
}

func (ctxt *context) parseFlags() error {
	flag.BoolVar(&ctxt.flags.pedantic, "pedantic", false,
		`makes several diagnostics more pedantic and comprehensive`)
	flag.BoolVar(&ctxt.flags.verbose, "v", false,
		`turn on additional info message printing`)
	flag.BoolVar(&ctxt.flags.debug, "debug", false,
		`turn on detailed program execution info printing`)
	flag.StringVar(&ctxt.flags.exclude, "exclude", `^unsafe$|^builtin$`,
		`import path excluding regexp`)

	flag.Parse()

	ctxt.flags.targets = flag.Args()
	if len(ctxt.flags.targets) == 0 {
		return fmt.Errorf("not enough positional args (empty targets list)")
	}

	return nil
}

func (ctxt *context) resolveTargets() error {
	ctxt.paths = gotool.ImportPaths(ctxt.flags.targets)
	if len(ctxt.paths) == 0 {
		return fmt.Errorf("targets resolved to an empty import paths list")
	}

	// Filter-out packages using the exclude pattern.
	excludeRE, err := regexp.Compile(ctxt.flags.exclude)
	if err != nil {
		return fmt.Errorf("compiling -exclude regexp: %v", err)
	}
	paths := ctxt.paths[:0]
	for _, path := range ctxt.paths {
		if !excludeRE.MatchString(path) {
			paths = append(paths, path)
		}
	}
	ctxt.paths = paths

	if len(paths) == 0 {
		ctxt.infoPrintf("import paths list is empty after filtering")
	}

	return nil
}

func (ctxt *context) initCheckers() error {
	checkers := []checker{
		newUnitImportChecker(ctxt),
		newZeroValPtrAllocChecker(ctxt),
		newEmptySliceChecker(ctxt),
		newEmptyMapChecker(ctxt),
		newHexLitChecker(ctxt),
		newRangeCheckChecker(ctxt),
		newAndNotChecker(ctxt),
		newFloatLitChecker(ctxt),
		newLabelCaseChecker(ctxt),
		newUntypedConstCoerceChecker(ctxt),
		newArgListParensChecker(ctxt),
		newNonZeroLenTestChecker(ctxt),
	}

	variantID := 0
	for _, c := range checkers {
		op := c.Operation()
		if op.name == "" {
			panic(fmt.Sprintf("%T: empty operation name", c))
		}
		for i, v := range op.variants {
			if v.warning == "" {
				panic(fmt.Sprintf("%T: empty warning for variant#%d", c, i))
			}
			v.op = op
			v.id = variantID
			variantID++
		}
	}

	ctxt.locs = newLocationMap()
	ctxt.checkers = checkers

	return nil
}

func (ctxt *context) collectAllCandidates() error {
	for _, path := range ctxt.paths {
		ctxt.infoPrintf("check %q", path)
		if err := ctxt.collectPathCandidates(path); err != nil {
			return fmt.Errorf("%s: %v", path, err)
		}
	}
	return nil
}

func (ctxt *context) collectPackageCandidates(pkg *packages.Package) {
	ctxt.info = pkg.TypesInfo
	for _, f := range pkg.Syntax {
		isGenerated := len(f.Comments) != 0 &&
			generatedFileCommentRE.MatchString(f.Comments[0].Text())
		if isGenerated {
			continue
		}
		ctxt.collectFileCandidates(f)
	}
}

func (ctxt *context) collectPathCandidates(path string) error {
	ctxt.fset = token.NewFileSet()

	conf := &packages.Config{
		Mode:  packages.LoadSyntax,
		Fset:  ctxt.fset,
		Tests: true,
	}

	// TODO(Quasilyte): current approach is memory-efficient
	// and does scale well with huge amounts of targets to check,
	// but it's not very fast. Might want to optimize it a little bit.
	pkgs, err := packages.Load(conf, path)
	if err != nil {
		return err
	}
	if len(pkgs) == 0 {
		ctxt.infoPrintf("got 0 packages for %q path", path)
		return nil
	}

	seenTests := false
	for _, pkg := range pkgs {
		// For some patterns Load returns 4 packages.
		// We need at most 2 and both of them should
		// have [$pkg.test] parts in their ID.
		if !strings.Contains(pkg.ID, ".test]") {
			continue
		}
		ctxt.collectPackageCandidates(pkg)
		if !strings.HasSuffix(pkg.Name, "_test") {
			seenTests = true
		}
	}
	if !seenTests {
		// Use the standard package if there were no tests.
		ctxt.collectPackageCandidates(pkgs[0])
	}

	return nil
}

func (ctxt *context) collectFileCandidates(f *ast.File) {
	ctxt.astinfo = astinfo.Info{
		Parents: make(map[ast.Node]ast.Node),
	}
	ctxt.astinfo.Origin = f
	ctxt.astinfo.Resolve()

	for _, c := range ctxt.checkers {
		for _, decl := range f.Decls {
			ast.Inspect(decl, func(n ast.Node) bool {
				return c.Visit(n)
			})
		}
	}
}

func (ctxt *context) assignSuggestions() error {
	for _, c := range ctxt.checkers {
		op := c.Operation()
		op.suggested = op.variants[0]
		for _, v := range op.variants[1:] {
			if v.count > op.suggested.count {
				op.suggested = v
			}
		}
	}
	return nil
}

func (ctxt *context) printWarnings() error {
	exitCode := 0
	visitWarings(ctxt, func(pos token.Position, v *opVariant) {
		exitCode = 1
		fmt.Printf("%s: %s: %s\n", pos, v.op.name, v.op.suggested.warning)
	})
	os.Exit(exitCode)
	return nil
}

func visitWarings(ctxt *context, visit func(pos token.Position, v *opVariant)) {
	// Build variant map which is accessed by variantID.
	vcount := 0
	for _, c := range ctxt.checkers {
		vcount += len(c.Operation().variants)
	}
	variants := make([]*opVariant, vcount)
	for _, c := range ctxt.checkers {
		for _, v := range c.Operation().variants {
			variants[v.id] = v
		}
	}

	for _, c := range ctxt.candidates {
		v := variants[c.variantID]
		if v.op.suggested == v {
			continue // OK, everything is consistent
		}
		pos := ctxt.locs.Get(c.locationID)
		visit(pos, v)
	}
}

func (ctxt *context) debugPrintf(format string, args ...interface{}) {
	if ctxt.flags.debug {
		log.Printf("\tdebug: "+format, args...)
	}
}

func (ctxt *context) infoPrintf(format string, args ...interface{}) {
	if ctxt.flags.verbose {
		log.Printf("\tinfo: "+format, args...)
	}
}
