//go:build go1.8
// +build go1.8

package mdson

import (
	"flag"
	"fmt"
	"github/drgo/mdson/tu"
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
	ctx= NewContext(DefaultOptions().SetDebug(ui.Debug(*TestDebug)))
	fmt.Fprintln(os.Stderr, ctx)
	e := m.Run()
    os.Exit(e)
}

// tu.Equal helper function to test if any two objects of the same type are tu.Equal
func initTestFiles(){
	blocks= tu.File.MustRead("test/blocks.md")
	lists =tu.File.MustRead("test/lists.md")
	specs =tu.File.MustRead("test/specs.md") 
}


// const TestDebug= DebugAll

func TestParseFileBlock(t *testing.T) {
	doc, err := ctx.ParseFile("test/blocks.md", nil)
	tu.Equal(t, err, nil)
	if err != nil {
		return
	}
	// t.Logf("%+v", doc)
	tu.Equal(t, doc.root.Kind(), LtBlock)
	tu.Equal(t, len(doc.root.Children()), 11)
}


func TestParseFileLists(t *testing.T) {
	doc, err := ctx.ParseFile("test/lists.md", nil)
	tu.Equal(t, err, nil)
	if err != nil {
		return
	}
	t.Logf("printout of root node: \n %+v", doc.root)
	tu.Equal(t, doc.root.Kind(), LtBlock)
}

func TestParse(t *testing.T) {
	doc, err := ctx.ParseFile("", specs)
	tu.Equal(t, err, nil)
	if err != nil {
		return
	}
	t.Logf("printout of doc: \n %+v", doc)
	tu.Equal(t, doc.root.Kind(), LtBlock)
}

func TestPrinter(t *testing.T) {
	doc, err := ctx.ParseFile("", specs)
	tu.Equal(t, err, nil)
	if err != nil {
		return
	}
	t.Logf("\n\n%s\n", newPrinter().print(doc.root)) 
}

