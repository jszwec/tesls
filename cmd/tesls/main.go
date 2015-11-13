package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/jszwec/tesls"
)

const usage = `%s [-format='format'] <packages>

Options:

	-format:
		it can be "json" or any other layout where
		%%T = test name
		%%P = package
		%%F = file path
		Default("%s")

tesls is looking for tests in the given list of packages.
It can also look for them recursively starting in the current directory by using: tesls ./...
`

const defaultFormat = "%P.%T %F"

var format = flag.String("format", defaultFormat, "")

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

func testDirs() set {
	var dirs = make(set)
	for _, arg := range flag.Args() {
		if strings.HasPrefix(arg, "-") {
			continue
		}
		arg, rec := recursiveArg(arg)
		dir, err := absDir(arg)
		check(err)
		dirs[dir] = struct{}{}
		if rec {
			check(walkfunc(dir, dirs))
		}
	}
	return dirs
}

func tests(dirs set) (ts tesls.TestSlice) {
	for dir := range dirs {
		t, err := tesls.Tests(dir)
		check(err)
		ts = append(ts, t...)
	}
	if len(ts) == 0 {
		check(errors.New("no tests were found"))
	}
	ts.Sort()
	return
}

func printTests(ts tesls.TestSlice) {
	switch *format {
	case "json":
		b, err := json.Marshal(ts)
		check(err)
		fmt.Fprintln(os.Stdout, string(b))
	case "":
		*format = defaultFormat
		fallthrough
	default:
		for _, t := range ts {
			fmt.Fprintln(os.Stdout, t.Format(*format))
		}
	}
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}
	printTests(tests(testDirs()))
}
