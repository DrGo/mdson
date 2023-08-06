package mdson

import (
	"io"
	"strings"
)

// produces groff with mom macros
//TODO: move to own package after testing

type mom struct {
	printer
	cfg 
}

func newMom(w io.Writer, ctx *Context) mom {
	m:= mom{
		cfg : newCfg(ctx),
		printer:  printer{w: w},
	}
	m.listMarker= "-"
	if m.ctx.DefaultListStyle=="ol" {
		m.listMarker="1."
	}	
	return m
}


func (m mom) printNode(n Node) {
	// m.ctx.Log("printNode():", n)
//- print general preamble 
//- use root attribs to create doc-specific preamble
//- output nodes 
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

func (m *mom) Transform(w io.Writer, d *Document) error  {
	m.printNode(d.root)
	return	nil
}





