package mdson

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/drgo/booker/tu"
)

func test_md_transform(t *testing.T, doc *Document) {
	t.Helper()
	w:= 	bufio.NewWriter(os.Stdout)
	md :=newMDTransformer(ctx) 
	md.Transform(w,doc)
	fmt.Println("********************************output************")
	w.Flush()
}

func TestTransformMD(t *testing.T) {
	doc, err := ctx.ParseFile("", specs)
	tu.Equal(t, err, nil)
	if err != nil {
		return
	}
	test_md_transform(t, doc)
}

func TestMDBlock(t *testing.T) {
	doc, err := ctx.ParseFile("test/blocks.md", nil)
	tu.Equal(t, err, nil)
	if err != nil {
		return
	}
	test_md_transform(t, doc)
}
