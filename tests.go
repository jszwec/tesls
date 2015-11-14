package tesls

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/types"
)

// Test describes a single test found in the *_test.go file
type Test struct {
	Name string `json:"name"`
	File string `json:"file"`
	Pkg  string `json:"pkg"`
}

// String returns a string representation of the Test
// in the form of 'package.Test filename'
func (t *Test) String() string {
	return fmt.Sprintf("%s.%s %s", t.Pkg, t.Name, t.File)
}

// TestSlice attaches the methods of sort.Interface to []Test.
// Sorting in increasing order comparing package+testname.
type TestSlice []Test

func (s TestSlice) Len() int           { return len(s) }
func (s TestSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TestSlice) Less(i, j int) bool { return s[i].Pkg+s[i].Name < s[j].Pkg+s[j].Name }

// Sort is a convenience method.
func (s TestSlice) Sort() { sort.Sort(s) }

func isTest(fdecl *ast.FuncDecl) bool {
	return strings.HasPrefix(fdecl.Name.String(), "Test") &&
		fdecl.Type != nil &&
		fdecl.Type.Params != nil &&
		len(fdecl.Type.Params.List) == 1 &&
		types.ExprString(fdecl.Type.Params.List[0].Type) == "*testing.T"
}

func isNoGoError(err error) bool {
	_, ok := err.(*build.NoGoError)
	return ok
}

// Tests function searches for test function declarations in the given directory.
func Tests(dir string) (tests TestSlice, err error) {
	pkg, err := build.ImportDir(dir, build.ImportMode(0))
	if err != nil && !isNoGoError(err) {
		return nil, err
	}
	fset := token.NewFileSet()
	for _, filename := range append(pkg.TestGoFiles, pkg.XTestGoFiles...) {
		filename = filepath.Join(dir, filename)
		f, err := parser.ParseFile(fset, filename, nil, parser.Mode(0))
		if err != nil {
			return nil, err
		}
		for _, decl := range f.Decls {
			fdecl, ok := decl.(*ast.FuncDecl)
			if ok && isTest(fdecl) {
				tests = append(tests, Test{
					Name: fdecl.Name.String(),
					File: filename,
					Pkg:  f.Name.String(),
				})
			}
		}
	}
	return tests, nil
}
