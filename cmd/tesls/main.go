package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jszwec/tesls"
)

const usage = `%s [-f='text/template'] <packages>

Options:

	-f:
		it can be "json" or any other layout where

			{{.Name}} = test name
			{{.Pkg}}  = package
			{{.File}} = file path

		Default("%s")

tesls is looking for tests in the given list of packages.
It can also look for them recursively starting in the current directory by using: tesls ./...
`
const defaultFormat = "{{.Pkg}}.{{.Name}} {{.File}}"

const iterationTemplate = "{{range .}}%s\n{{end}}"

var defaultTemplate = template.Must(
	template.New("default").Parse(fmt.Sprintf(iterationTemplate, defaultFormat)))

var format = flag.String("f", defaultFormat, "")

type set map[string]struct{}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0], defaultFormat)
	}
}

func check(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func walkfunc(root string, dirs set) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".git") {
				return filepath.SkipDir
			}
			abs, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			dirs[abs] = struct{}{}
		}
		return nil
	})
}

func recursiveArg(arg string) (string, bool) {
	if strings.HasSuffix(arg, "/...") {
		return arg[:len(arg)-4], true
	}
	return arg, false
}

func absDir(arg string) (string, error) {
	if strings.HasPrefix(arg, ".") {
		return filepath.Abs(arg)
	}
	p, err := build.Import(arg, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}

func testDirs(args []string) (set, error) {
	var dirs = make(set)
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		arg, rec := recursiveArg(arg)
		dir, err := absDir(arg)
		if err != nil {
			return nil, err
		}
		dirs[dir] = struct{}{}
		if rec {
			if err := walkfunc(dir, dirs); err != nil {
				return nil, err
			}
		}
	}
	return dirs, nil
}

func tests(dirs set) (ts tesls.TestSlice, err error) {
	for dir := range dirs {
		t, err := tesls.Tests(dir)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t...)
	}
	if len(ts) == 0 {
		return nil, errors.New("no tests were found")
	}
	ts.Sort()
	return ts, nil
}

func getTemplate(format string) (*template.Template, error) {
	switch format {
	case "", defaultFormat:
		return defaultTemplate, nil
	default:
		return template.New("TestTemplate").Parse(fmt.Sprintf(iterationTemplate, format))
	}
}

func printTests(w io.Writer, ts tesls.TestSlice, format string, t *template.Template) error {
	switch format {
	case "json":
		b, err := json.Marshal(ts)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w, string(b)); err != nil {
			return err
		}
	default:
		return t.Execute(w, ts)
	}
	return nil
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}
	t, err := getTemplate(*format)
	check(err)
	dirs, err := testDirs(flag.Args())
	check(err)
	ts, err := tests(dirs)
	check(err)
	check(printTests(os.Stdout, ts, *format, t))
}
