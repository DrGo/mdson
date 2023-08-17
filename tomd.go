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

//TransformerConfig holds basic output control vars 
type TransformerConfig struct{
	Indent string
	ListMaker string 
	TabWidth int
}

func DefaultTransformerConfig() TransformerConfig{
	cfg := TransformerConfig {
		TabWidth: 2,
	}
	return cfg 
}

var _ Transformer=&MDTransformer{}

type MDTransformer struct {
	printer
	TransformerConfig 
}

func NewMDTransformer(cfg TransformerConfig) MDTransformer {
	mdt:= MDTransformer{
		TransformerConfig :cfg,
	}
	return mdt 
}

func (m MDTransformer) printNode(n Node) {
	// m.ctx.Log("printNode():", n)
	switch n := n.(type) {
	case *ttBlock:
		if n.Level() > 0 { // donot print root's title
			m.println(strings.Repeat("#", n.Level()) + " " + n.Value())
		}
		m.Indent = strings.Repeat(" ", n.Level()* m.TabWidth)
		for _, n := range n.Children() {
			m.printNode(n)
		}
	case *ttList:
		m.println(n.Value())
		m.Indent = strings.Repeat(" ", n.Level()* m.TabWidth)
		for _, n := range n.Children() {
			m.printNode(n)
		}
	case *ttListItem:
		m.print(m.Indent+ m.ListMaker)
		m.println(n.Value())
	default:
		m.println(n.Value())
	}
}

// TODO: check for writing errors
func (m MDTransformer) Transform(w io.Writer, doc *Document) error {
	m.printer = printer{w} 
	m.printNode(doc.root)
	// m.w.Flush()
	return nil
}
