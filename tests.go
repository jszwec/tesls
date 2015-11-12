package tesls

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
)

// Test describes a single test found in the *_test.go file
type Test struct {
	Name string `json:"name"`
	File string `json:"file"`
	Pkg  string `json:"pkg"`
}

// Format returns Test struct as a string in the given format.
// It can be any string in which %P will turn into package name,
// %T into test name, and %F into file path in which test was found.
// Note that it is not required to use all three placeholders.
func (t *Test) Format(format string) string {
	return strings.NewReplacer("%P", t.Pkg, "%T", t.Name, "%F", t.File).Replace(format)
}

// String returns a string representation of the Test
// in the form of 'package.Test filename'
func (t *Test) String() string {
	return t.Format("%P.%T %F")
}

// TestSlice attaches the methods of sort.Interface to []Test.
// Sorting in increasing order comparing package+testname.
type TestSlice []Test

func (s TestSlice) Len() int           { return len(s) }
func (s TestSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TestSlice) Less(i, j int) bool { return s[i].Pkg+s[i].Name < s[j].Pkg+s[j].Name }

// Sort is a convenience method.
func (s TestSlice) Sort() { sort.Sort(s) }

func filter(fi os.FileInfo) bool {
	return strings.HasSuffix(fi.Name(), "_test.go")
}

// Tests function searches for test function declarations in the given directory.
func Tests(dir string) (tests TestSlice, err error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, filter, parser.Mode(0))
	if err != nil {
		return nil, err
	}
	for _, pkg := range pkgs {
		for filename, file := range pkg.Files {
			for _, decl := range file.Decls {
				fdec, ok := decl.(*ast.FuncDecl)
				if ok && strings.HasPrefix(fdec.Name.String(), "Test") {
					tests = append(tests, Test{
						Name: fdec.Name.String(),
						File: filename,
						Pkg:  pkg.Name,
					})
				}
			}
		}
	}
	return tests, nil
}
