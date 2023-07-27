package mdson

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func test_md_transform(t *testing.T, doc *Document) {
	t.Helper()
	w:= 	bufio.NewWriter(os.Stdout)
	md :=newMDTransformer(w, ctx) 
	md.Transform( doc)
	fmt.Println("********************************output************")
	w.Flush()
}

func TestTransformMD(t *testing.T) {
	doc, err := ctx.ParseFile("", specs)
	Equal(t, err, nil)
	if err != nil {
		return
	}
	test_md_transform(t, doc)
}

func TestMDBlock(t *testing.T) {
	doc, err := ctx.ParseFile("test/blocks.md", nil)
	Equal(t, err, nil)
	if err != nil {
		return
	}
	test_md_transform(t, doc)
}
