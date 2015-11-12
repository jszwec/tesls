package tesls

import (
	"fmt"
	"go/build"
	"path/filepath"
	"reflect"
	"testing"
)

const pkg = "github.com/jszwec/tesls"

func TestTests(t *testing.T) {
	p, err := build.Import(pkg, "", build.FindOnly)
	if err != nil {
		t.Fatal(err)
	}
	expectedTest := Test{
		Name: "TestTests",
		File: filepath.Join(p.Dir, "tests_test.go"),
		Pkg:  "tesls",
	}
	ts, err := Tests(p.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(ts) != 1 {
		t.Errorf("expected len(ts) = 1; got %d", len(ts))
	}
	if !reflect.DeepEqual(expectedTest, ts[0]) {
		t.Errorf("expected %v; got %v", expectedTest, ts[0])
	}
	expectedTestStr := fmt.Sprintf("test%s%s%s", expectedTest.Name,
		expectedTest.Pkg, expectedTest.File)
	TestStr := ts[0].Format(`test%T%P%F`)
	if TestStr != expectedTestStr {
		t.Errorf("expected %s; got %s", expectedTestStr, TestStr)
	}
}
