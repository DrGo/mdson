package mdson

import (
	"strings"
)

//TODO add Document interface
// Document represents parsed and evaluated mdson Document
// It implements fs.File to enable reading and seeking the evaluated text
type Document struct {
	ctx *Context
	root        BlockNode 
	attribs  map[string]string
}

func newDocument(ctx *Context) *Document{
	return &Document{
		ctx :ctx, 
		root: newBlock("root",0),

		attribs: make(map[string]string),
	}
}


func (doc Document) Attribs() map[string]string {
	return doc.attribs
}

func (doc Document) String() string {
	var sb strings.Builder
	sb.Grow(1024)
	indent := " "
	sb.WriteString(doc.root.String())
	for k, v := range doc.attribs {
		sb.WriteString(indent + k)
		sb.WriteString(" = ")
		sb.WriteString(v)
		sb.WriteRune('\n')
	}
	return sb.String()
}

