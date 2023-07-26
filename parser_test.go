//go:build go1.8
// +build go1.8

package mdson

import (
	"io"
	"testing"
)

// Equal helper function to test if any two objects of the same type are equal
var blocks, lists, specs io.Reader
func init(){
	blocks= mustReadFile("test/blocks.md")
	lists = mustReadFile("test/lists.md")
	specs = mustReadFile("test/specs.md") 
}

var ctx= NewContext(DefaultParserOptions().SetDebug(*TestDebug))

// const TestDebug= DebugAll

func TestParseFileBlock(t *testing.T) {
	doc, err := ctx.ParseFile("test/blocks.md", nil)
	Equal(t, err, nil)
	if err != nil {
		return
	}
	t.Logf("%+v", doc)
	Equal(t, doc.root.Kind(), LtBlock)
}


func TestParseFileLists(t *testing.T) {
	doc, err := ctx.ParseFile("test/lists.md", nil)
	Equal(t, err, nil)
	if err != nil {
		return
	}
	t.Logf("printout of root node: \n %+v", doc.root)
	Equal(t, doc.root.Kind(), LtBlock)
}

func TestParse(t *testing.T) {
	doc, err := ctx.ParseFile("", specs)
	Equal(t, err, nil)
	if err != nil {
		return
	}
	t.Logf("printout of doc: \n %+v", doc)
	Equal(t, doc.root.Kind(), LtBlock)
}

func TestPrinter(t *testing.T) {
	doc, err := ctx.ParseFile("", specs)
	Equal(t, err, nil)
	if err != nil {
		return
	}
	t.Logf("\n\n%s\n", newPrinter().print(doc.root)) 
}

