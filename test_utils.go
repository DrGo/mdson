package mdson

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
)

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

func mustReadFile(filename string) io.Reader {	
	content, err := os.ReadFile(filename)
	if err != nil {
		panic("test setup failed" + err.Error())
	}
	return bytes.NewReader(content) 
}


func assert(cond bool, msg string) {
	if !cond {
		panic("assertion failed: "+ msg)
	}
}	
