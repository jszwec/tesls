package tesls

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

const pkg = "github.com/jszwec/tesls"

const codeTest = `%spackage %s

import "testing"

func %s(t *testing.T) {}
`

var goos = []string{
	"android",
	"darwin",
	"dragonfly",
	"freebsd",
	"linux",
	"nacl",
	"netbsd",
	"openbsd",
	"plan9",
	"solaris",
	"windows",
}

func chooseOtherOS() string {
	for _, os := range goos {
		if runtime.GOOS != os {
			return os
		}
	}
	return ""
}

func createTestFiles(t *testing.T, dir string) {
	if err := ioutil.WriteFile(filepath.Join(dir, "package_test.go"),
		[]byte(fmt.Sprintf(codeTest, "", "tesls_test", "TestZPackage")), 0655); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "tag_test.go"),
		[]byte(fmt.Sprintf(codeTest,
			"// +build "+chooseOtherOS()+"\n\n", "tesls", "TestTag")), 0655); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "os_"+chooseOtherOS()+"_test.go"),
		[]byte(fmt.Sprintf(codeTest,
			"", "tesls", "TestDifferentOS")), 0655); err != nil {
		t.Fatal(err)
	}
}

func TestTests(t *testing.T) {
	p, err := build.Import(pkg, "", build.FindOnly)
	if err != nil {
		t.Fatal(err)
	}
	createTestFiles(t, p.Dir)
	expected := TestSlice{
		{
			Name: "TestTests",
			File: filepath.Join(p.Dir, "tests_test.go"),
			Pkg:  "tesls",
		},
		{
			Name: "TestZPackage",
			File: filepath.Join(p.Dir, "package_test.go"),
			Pkg:  "tesls_test",
		},
	}
	ts, err := Tests(p.Dir)
	if err != nil {
		t.Fatal(err)
	}
	expected.Sort()
	ts.Sort()
	if len(ts) != 2 {
		t.Errorf("expected len(ts) = 2; got %d", len(ts))
	}
	if !reflect.DeepEqual(expected, ts) {
		t.Errorf("expected %v; got %v", expected, ts)
	}
}
