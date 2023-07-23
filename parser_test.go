//go:build go1.8
// +build go1.8

package mdson

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

// Equal helper function to test if any two objects of the same type are equal
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper() //report error in the file that calls this func
	if expected != actual {
		t.Errorf("wanted: %v; got: %v", expected, actual)
	}
}

// isNil gets whether the object is nil or not.
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}
	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false
}

func NotNil(t *testing.T, obj any) {
	t.Helper()
	if isNil(obj) {
		t.Errorf("%v is nil", obj)
	}
}

func mustReadFile(filename string) string {
	
content, err := os.ReadFile(filename)
	if err != nil {
		panic("test setup failed" + err.Error())
	}
	return string(content)
}

var blocks, lists, specs string 
func init(){
	blocks= mustReadFile("test/blocks.md")
	lists = mustReadFile("test/lists.md")
	specs = mustReadFile("test/specs.md") 
}

func TestParseFileBlock(t *testing.T) {
	SetDebug(DebugAll)
	n, err := ParseFile("test/blocks.md")
	Equal(t, err, nil)
	if err != nil {
		return
	}
	t.Logf("%+v", n)
	Equal(t, n.Kind(), LtBlock)
}


func TestParseFileLists(t *testing.T) {
	SetDebug(DebugAll)
	n, err := ParseFile("test/lists.md")
	Equal(t, err, nil)
	if err != nil {
		return
	}
	t.Logf("printout of root node: \n %+v", n)
	Equal(t, n.Kind(), LtBlock)
}

func TestParse(t *testing.T) {
	SetDebug(DebugAll)
	n, err := Parse(strings.NewReader(specs))
	Equal(t, err, nil)
	if err != nil {
		return
	}
	t.Logf("%+v", n)
	Equal(t, n.Kind(), LtBlock)
}

func TestPrinter(t *testing.T) {

	n, err := Parse(strings.NewReader(specs))
	Equal(t, err, nil)
	if err != nil {
		return
	}
	t.Logf("\n\n%s\n", newPrinter().print(n)) 
}

