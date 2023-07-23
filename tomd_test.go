package mdson

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func test_md_transform(t *testing.T, n Node) {
	t.Helper()
	w:= 	bufio.NewWriter(os.Stdout)
	md :=newMDTransformer(w, DefaultParserOptions().SetDebug(DebugAll)) 
	md.Transform( n)
	fmt.Println("********************************output************")
	w.Flush()
}

func TestTransformMD(t *testing.T) {
	SetDebug(DebugWarning)
	n, err := Parse(strings.NewReader(specs))
	Equal(t, err, nil)
	if err != nil {
		return
	}
	test_md_transform(t, n)
}

func TestMDBlock(t *testing.T) {
	SetDebug(DebugAll)
	n, err := ParseFile("test/blocks.md")
	Equal(t, err, nil)
	if err != nil {
		return
	}
	test_md_transform(t, n)
}
