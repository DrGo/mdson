package mdson

import (
	"fmt"
	"io"
	"strings"

	"github.com/drgo/core/ui"
)

const EOL = "\r\n"

type Transformer interface {
	Transform(w io.Writer, n Node) error
}

type MDTransformer struct {
	w io.Writer
	ui.UI
}

func newMDTransformer(w io.Writer, options *ParserOptions) MDTransformer {
	return MDTransformer{
		w:  w,
		UI: ui.NewUI(options.Debug),
	}
}
func (m MDTransformer) print(s string) {
	fmt.Fprint(m.w, s, EOL)
	// m.w.Write([]byte(EOL))
}

func (m MDTransformer) printNode(n Node) {
	m.Log("printNode():", n)
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
