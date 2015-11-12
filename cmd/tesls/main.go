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

type once bool

func (o *once) Do(f func()) {
	if *o {
		return
	}
	*o = true
	f()
}

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

func walkfunc(dirs set) error {
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
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

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}
	dirs := make(set)
	var once once
	for _, arg := range flag.Args() {
		switch {
		case strings.HasPrefix(arg, "-"):
			continue
		case arg == "./...":
			once.Do(func() { check(walkfunc(dirs)) })
		default:
			p, err := build.Import(arg, "", build.FindOnly)
			check(err)
			dirs[p.Dir] = struct{}{}
		}
	}
	var ts tesls.TestSlice
	for dir := range dirs {
		t, err := tesls.Tests(dir)
		check(err)
		ts = append(ts, t...)
	}
	if len(ts) == 0 {
		check(errors.New("no tests were found"))
	}
	ts.Sort()
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
