package mdson

import (
	"fmt"
	"io"
	"strings"
)

const EOL = "\r\n"

type Transformer interface {
	Transform(w io.Writer, d *Document) error
}

type printer struct{
	w io.Writer
}

func (p printer) print(s string) {
	fmt.Fprint(p.w,  s)
}

func (p printer) println(s string) {
	fmt.Fprint(p.w, s, EOL)
}

//cfg holds basic output control vars 
type cfg struct{
	ctx *Context
	indent string
	listMarker string 
	TabWidth int
}

func newCfg(ctx *Context) cfg{
	cfg := cfg {
		ctx: ctx,
		TabWidth: 2,
	}
	return cfg 
}

type MDTransformer struct {
	printer
	cfg 
}

func newMDTransformer(w io.Writer, ctx *Context) MDTransformer {
	mdt:= MDTransformer{
		cfg :newCfg(ctx),
		printer:  printer{w: w},
	}
	mdt.listMarker= "-"
	if mdt.ctx.DefaultListStyle=="ol" {
		mdt.listMarker="1."
	}	
	return mdt 
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
