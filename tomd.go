package mdson

import (
	"fmt"
	"io"
)

const EOL = "\n"
type Transformer interface{
	Transform(w io.Writer, n Node) error 
}

//
type MDTransformer struct{}

func (m MDTransformer) Transform(w io.Writer, n Node) error{	
	for _, n := range n.Children() {
		fmt.Println("node: ", n)
		w.Write([]byte(n.Value()))
		w.Write([]byte(EOL))
	}	
	return nil 
}
