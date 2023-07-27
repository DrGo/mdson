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
}

func newMDTransformer(w io.Writer, ctx *Context) MDTransformer {
	return MDTransformer{
		ctx :ctx,
		w:  w,
	}
}
func (m MDTransformer) print(s string) {
	fmt.Fprint(m.w, s, EOL)
	// m.w.Write([]byte(EOL))
}

func (m MDTransformer) printNode(n Node) {
	m.ctx.Log("printNode():", n)
	switch n := n.(type) {
	case *ttBlock:
		if n.Level() > 0 { // donot print root
			m.print(strings.Repeat("#", n.Level()) + " " + n.Value())
		}	
			for _, n := range n.Children() {
				m.printNode(n)
			}
	case *ttList:
		m.print(n.Value())
			for _, n := range n.Children() {
				m.printNode(n)
			}
	default:
		m.print(n.Value())
	}
}

// TODO: check for writing errors
func (m MDTransformer) Transform(doc *Document) error {
	m.printNode(doc.root)
	// m.w.Flush()
	return nil
}
