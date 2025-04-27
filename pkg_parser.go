package gotest_labels

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
)

const DefaultLabelValue = "true"

var defaultPkg = "./..."

// Get go packages in "." directory since the packages and paths are actually processed earlier than
// executing the test binaries internally by the go test command. The current package only needs to
// take care of the current directory.
func getPackages() ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles,
		Tests: true,
	}
	pkgs, err := packages.Load(cfg, defaultPkg)

	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %v", err)
	}

	// Filter out packages without go files
	var validPkgs []*packages.Package
	for _, pkg := range pkgs {
		if len(pkg.GoFiles) > 0 || len(pkg.CompiledGoFiles) > 0 {
			validPkgs = append(validPkgs, pkg)
		}
	}
	return validPkgs, nil
}

// Find all *_test.go files under the given package
func getTestFiles(pkg *packages.Package) []string {
	var testFiles []string

	// Ensure all go files are included
	allFiles := append(pkg.GoFiles, pkg.CompiledGoFiles...)
	for _, file := range allFiles {
		if strings.HasSuffix(file, "_test.go") {
			testFiles = append(testFiles, file)
		}
	}
	return testFiles
}

// Find all Test* functions with (t *testing.T) signature and matching the label filter in given test files
// It returns a map of function names to their labels
func FindTestFuncs(testFiles []string, filterAST Node) (map[string]TestLabels, error) {
	testFuncs := map[string]TestLabels{}
	fset := token.NewFileSet()

	for _, file := range testFiles {
		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s, err: %v", file, err)
		}

		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name == nil || !strings.HasPrefix(fn.Name.Name, "Test") {
				continue
			}

			if !isValidTestFunc(fn) {
				continue
			}

			labels := getFuncLabels(fn)
			if !Evaluate(filterAST, labels) {
				continue
			}

			testFuncs[fn.Name.Name] = labels
		}
	}
	return testFuncs, nil
}

// Check function signature is func Test*(t *testing.T)
func isValidTestFunc(fn *ast.FuncDecl) bool {
	if len(fn.Type.Params.List) != 1 {
		return false
	}
	param, ok := fn.Type.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := param.X.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "T"
}

func filterTestFuncs(funcs map[string]TestLabels, regex *regexp.Regexp) map[string]TestLabels {
	if regex == nil {
		return funcs
	}
	matched := make(map[string]TestLabels)
	for name, labels := range funcs {
		if regex.MatchString(name) {
			matched[name] = labels
		}
	}
	return matched
}

// Get the labels from the function comments
func getFuncLabels(fn *ast.FuncDecl) TestLabels {
	tags := make(TestLabels)
	if fn.Doc == nil {
		return tags
	}
	for _, comment := range fn.Doc.List {
		text := strings.TrimSpace(comment.Text)
		// For both styles in `// @key=value` and `/* @key=value */`
		text = strings.TrimPrefix(text, "//")
		text = strings.TrimPrefix(text, "/*")
		text = strings.TrimSuffix(text, "*/")
		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, "@") {
			parts := strings.SplitN(text[1:], "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				tags[key] = value
			} else {
				key := strings.TrimSpace(parts[0])
				tags[key] = DefaultLabelValue
			}
		}
	}
	return tags
}
