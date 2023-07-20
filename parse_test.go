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

// func TestLines(t *testing.T) {
// 	const data = `line 1
// line 2
// line 3	
// `
// 	lines, err := newLines(strings.NewReader(data))
// 	Equal(t, err, nil)
// 	NotNil(t, lines)
// 	Equal(t, lines.rowNum, 0)
// 	// t.Logf("%s", lines)
// 	//    t.Log(lines.rows[0])
// 	// t.Log(lines)
// 	line, ok := lines.next()
// 	Equal(t, line, "line 1")
// 	Equal(t, ok, true)
// 	lines, err = newLinesFromFile("README.md")
// 	Equal(t, err, nil)
// 	NotNil(t, lines)
//
// 	//   t.Log(lines.rows[0])
// 	// t.Log(lines)
// 	// t.Log(lines.next())
// }

func TestParse(t *testing.T) {
	const data =`
// comments will be removed before processing	
// a property (prop) is essentially like a variable; it will be removed from the text
// but made available to the rendering program as part of the DOM
// syntax: dot. followed immediately by a letter or underscore and then colon.
// anything after the colon is assigned as a value of type string
// all props belong to the nearest section (see below) except the ones at the beginning
// of the file before any section header is declared which belong to the root section (the document)
.prop can have long name: some value
.weight: 1
.date: 12july2023
.prop with multiline value: all Text lines (see below) following the a prop are considered
		part of the value of that prop regardless of indentation level
 .this is not a prop because the first char is space

// you could refer to any of the above props anywhere in the doc like this |.date|
.today: Today is |.date|		

// markdown headers and section names start with #
## introduction
// props that belong to this section
.author: Someone
// they can be referred to as |introduction.author|
// refer to the computed .today value above
|.today|

// A list starts with . and continues until the next non-list item element 
// this is a list named "Causes of heart failure". The name is all the text
// before the colon which is mandatory
.Causes of heart failure:
//- list item starts with -
- Hypertension
- Atrial fibrillation
- Myocardial infarction
// This list can be referred to as |.Causes of heart failure|
// its first element can be referred to as |.Causes of heart failure[0]|

// backticks used to guard preformatted text like code blocks 
// everything else is regular text

// Referring to another section using []()


`
SetDebug(DebugAll)
n, err := Parse(strings.NewReader(data))
	Equal(t, err, nil)
	if err != nil {
		return
	}
	t.Logf("%+v", n)
	Equal(t,n.Kind(),"block")
}
