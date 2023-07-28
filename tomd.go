package mdson

import (
	"fmt"
	"io"
	"strings"
)

const EOL = "\r\n"

type Transformer interface {
	Transform(w io.Writer, n Node) error
}

type MDTransformer struct {
	w io.Writer
	ctx *Context
	indent string
	listMarker string 
	TabWidth int
}

func newMDTransformer(w io.Writer, ctx *Context) MDTransformer {
	mdt:= MDTransformer{
		ctx :ctx,
		w:  w,
		TabWidth: 2,
	}
	mdt.listMarker= "-"
	if mdt.ctx.DefaultListStyle=="ol" {
		mdt.listMarker="1."
	}	
	return mdt 
}

func (m MDTransformer) print(s string) {
	fmt.Fprint(m.w,  s)
}

func (m MDTransformer) println(s string) {
	fmt.Fprint(m.w, s, EOL)
}

func (m MDTransformer) printNode(n Node) {
	m.ctx.Log("printNode():", n)
	switch n := n.(type) {
	case *ttBlock:
		if n.Level() > 0 { // donot print root's title
			m.println(strings.Repeat("#", n.Level()) + " " + n.Value())
		}
		m.indent = strings.Repeat(" ", n.Level()* m.TabWidth)
		for _, n := range n.Children() {
			m.printNode(n)
		}
	case *ttList:
		m.println(n.Value())
		m.indent = strings.Repeat(" ", n.Level()* m.TabWidth)
		for _, n := range n.Children() {
			m.printNode(n)
		}
	case *ttListItem:
		m.print(m.indent+ m.listMarker)
		m.println(n.Value())
	default:
		m.println(n.Value())
	}
}

// TODO: check for writing errors
func (m MDTransformer) Transform(doc *Document) error {
	m.printNode(doc.root)
	// m.w.Flush()
	return nil
}
