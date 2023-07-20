package mdson

import (
	"strings"
)

type Printer struct {
	sb strings.Builder
}


func newPrinter() *Printer {
	return &Printer{
		sb: strings.Builder{},
	}
}

func (w *Printer) print(n Node) string {
	w.sb.WriteString(n.String())
	return w.sb.String()
}

