package mdson

import (
	"fmt"
	"io"
	"strings"
)

const EOL = "\r\n"
type Transformer interface{
	Transform(w io.Writer, n Node) error 
}

//
type MDTransformer struct{
	w io.Writer
}

func (m MDTransformer) print(s string) {
		fmt.Fprint(m.w, s , EOL)
		// m.w.Write([]byte(EOL))
	}

func (m MDTransformer) printNode( n Node){
		if n.Level()> -1 {
			m.print(strings.Repeat("*", n.Level()) + " " + n.Value())
			for _, n := range n.Children() {
				m.printNode( n)
			}
			return 
		}
		fmt.Print(n.Kind(), ": ")
		print(n.Value())
	}	

// TODO: check for writing errors
func (m MDTransformer) Transform(n Node) error{	
	m.printNode(n)
	// m.w.Flush()
	return nil 
}
