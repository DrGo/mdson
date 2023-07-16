//go:build go1.8
// +build go1.8

package mdson

import (
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

func TestLines(t *testing.T) {
	const data = `line 1
line 2
line 3	
`
	lines, err := newLines(strings.NewReader(data))
	Equal(t, err, nil)
	NotNil(t, lines)
	Equal(t, lines.rowNum, 0)
	// t.Logf("%s", lines)
	//    t.Log(lines.rows[0])
	// t.Log(lines)
	line, ok := lines.next()
	Equal(t, line, "line 1")
	Equal(t, ok, true)
	lines, err = newLinesFromFile("README.md")
	Equal(t, err, nil)
	NotNil(t, lines)

	//   t.Log(lines.rows[0])
	// t.Log(lines)
	// t.Log(lines.next())
}
