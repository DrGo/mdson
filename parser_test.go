//go:build go1.8
// +build go1.8

package mdson

import (
	"flag"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/drgo/core/ui"
)

var (
	blocks, lists, specs io.Reader
	ctx *Context
	TestDebug=flag.Int("debug", 0, "")

)

func TestMain(m *testing.M) {
	flag.Parse()
	initTestFiles()
	ctx= NewContext(DefaultParserOptions().SetDebug(ui.Debug(*TestDebug)))
	fmt.Fprintln(os.Stderr, ctx)
	e := m.Run()
    os.Exit(e)
}

// Equal helper function to test if any two objects of the same type are equal
func initTestFiles(){
	blocks= mustReadFile("test/blocks.md")
	lists = mustReadFile("test/lists.md")
	specs = mustReadFile("test/specs.md") 
}


// const TestDebug= DebugAll

func TestParseFileBlock(t *testing.T) {
	doc, err := ctx.ParseFile("test/blocks.md", nil)
	Equal(t, err, nil)
	if err != nil {
		return
	}
	// t.Logf("%+v", doc)
	Equal(t, doc.root.Kind(), LtBlock)
	Equal(t, len(doc.root.Children()), 11)
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

