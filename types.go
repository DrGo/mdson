package mdson

import (
	"fmt"
	"strconv"
	"strings"
)

// An ID describes a field that can be automatically filled with MDSon block name
type ID string

type LineType int

const (
	LtReadError LineType = iota
	LtSyntaxError
	LtEOF
	LtEmpty
	LtComment
	LtList
	LtListItem
	LtBlock
	// LtLiteralString
	LtAttrib
	LtTextLine
)

func (lt LineType) String() string {
	//[...] creates an array rather than a slice
	return [...]string{"Read Error", "Syntax Error", "EOF",
		"Empty", "Comment", "List", "List item", "Block", "Attribute", "Text line"}[lt]
}

// Node implement parser's AST leaf node
type Node interface {
	String() string
	Kind() LineType
	Key() string
	// Children() []Node
	// NthChild(idx int) Node
	// ValueOf() map[string]string
	LineNum() int
	SetLineNum(value int)
	//returns the textual representation of the node as it should appear in a document
	Value() string
	// sets the textual representation of the node as it should appear in a document
	SetValue(s string) Node
}

type BlockNode interface {
	Node
	Children() []Node
	NthChild(idx int) Node
	AddChild(n Node) Node
	// returns nesting level for a node
	Level() int
}

// baseToken implements the basic token interface root of all of other tokens
type ttBase struct {
	lnum  int
	kind  LineType
	level int
	key   string
}

func newBase(kind LineType, key string) *ttBase {
	return &ttBase{kind: kind, key: key}
}

func (bt ttBase) LineNum() int {
	return bt.lnum
}

func (bt *ttBase) SetLineNum(value int) {
	bt.lnum = value
}

const nodeDescLine = "type=%s, lineNum=%d, key='%s',value='%s'"

func (bt ttBase) String() string {
	return fmt.Sprintf(nodeDescLine, bt.Kind(), bt.LineNum(), bt.Key(), bt.Value())
}

func (bt ttBase) Kind() LineType {
	return bt.kind
}

func (bt ttBase) Key() string {
	return bt.key //GetValidVarName(bt.key)
}

// func (bt ttBase) isArray() bool {
// 	return isArray(bt.key)
// }
//
// func (bt ttBase) Children() []Node {
// 	return nil
// }
//
// func (bt ttBase) NthChild(idx int) Node{
// 	return nil
// }
//
// func (bt ttBase) ValueOf() map[string]string {
// 	return nil
// }

func (bt ttBase) Value() string {
	return bt.key
}

func (bt *ttBase) SetValue(s string) Node {
	bt.key = s
	return bt
}

type ttReadError struct {
	*ttBase
	err error
}

func newReadError(value interface{}) *ttReadError {
	re := ttReadError{ttBase: newBase(LtReadError,"")}
	switch unboxed := value.(type) {
	case string:
		re.err = fmt.Errorf("%s", unboxed)
	case error:
		re.err = unboxed
	default:
		panic("unsupported argument type in newReadError")
	}
	return &re
}

func (re ttReadError) String() string {
	return re.ttBase.String() + re.err.Error()
}

type ttSyntaxError struct{ ttReadError }

func newSyntaxError(value interface{}) *ttSyntaxError {
	se := ttSyntaxError{ttReadError: *newReadError(value)}
	se.kind = LtSyntaxError

	// return fmt.Sprintf("line %d syntax error: %s)", ese.lineNum, ese.message)
	return &se
}

type ttBlock struct {
	*ttBase
	level    int
	children []Node
}

func newBlock(key string, level int) *ttBlock {
	return &ttBlock{ttBase: newBase(LtBlock, key),
		level: level,}
}

func (blk ttBlock) String() string {
	var sb strings.Builder
	sb.Grow(1024)
	indent := strings.Repeat(" ", blk.Level())
	sb.WriteString(indent + blk.ttBase.String() + " " + strconv.Itoa(blk.Level()) + "\n")
	indent = indent + "- "
	for _, t := range blk.children {
		sb.WriteString(indent + t.String() + "\n")
	}
	return sb.String()
}

//
// func (blk ttBlock) Key() string {
// 	name := GetValidVarName(blk.key)
// 	if blk.isArray() {
// 		return strings.TrimSuffix(name, "list")
// 	}
// 	return name
// }

// func (blk *ttBlock) setName(value string) *ttBlock {
// 	blk.key = value
// 	return blk
// }

func (bt *ttBlock) setLevel(value int) Node {
	bt.level = value
	return bt
}

func (bt ttBlock) Level() int {
	return bt.level
}

func (blk *ttBlock) AddChild(n Node) Node {
	blk.children = append(blk.children, n)
	return blk
}

func (blk ttBlock) getChildByName(name string) Node {
	for _, c := range blk.children {
		if c.Key() == name {
			return c
		}
	}
	return nil
}

func (blk ttBlock) Children() []Node {
	return blk.children
}

func (blk ttBlock) NthChild(idx int) Node {
	return blk.children[idx]
}

// func (blk ttBlock) ChildByName(name string) Node {
// 	return blk.getChildByName(name)
// }

type ttAttrib struct {
	*ttBase
	value string
}

func newAttrib(key, value string) *ttAttrib {
	return &ttAttrib{ttBase: newBase(LtAttrib, strings.TrimSpace(key)),
		value: strings.TrimSpace(value)}
}

func (att ttAttrib) String() string {
	return att.ttBase.String() + ": " + att.value
}

func (kvp *ttAttrib) setKey(value string) *ttAttrib {
	kvp.key = value
	return kvp
}

func (kvp *ttAttrib) setValue(value string) *ttAttrib {
	kvp.value = value
	return kvp
}

type ttList struct {
	*ttBlock
}

func newList(name string, level int) *ttList {
	list := &ttList{ttBlock: newBlock(name, level)}
	list.kind = LtList
	return list
}

func (list ttList) String() string {
	var sb strings.Builder
	sb.Grow(1024)
	indent := strings.Repeat(" ", list.Level())
	sb.WriteString(indent + list.ttBase.String() + " " + strconv.Itoa(list.Level()) + "\n")
	indent = indent + "- "
	// for k, v := range list.attribs {
	// 	sb.WriteString(indent+ k)
	// 	sb.WriteString(" = ")
	// 	sb.WriteString(v)
	// 	sb.WriteRune('\n')
	// }
	for _, t := range list.children {
		sb.WriteString(indent + t.String() + "\n")
	}
	return sb.String()
}

// func (list ttList) Children() []Node {
// 	return list.children
// }
// func (blk ttList) NthChild(idx int) Node{
// 	return blk.children[idx]
// }
//

// func GetValidVarName(s string) string {
// 	ns := []rune{}
// 	for _, c := range s {
// 		if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' {
// 			ns = append(ns, c)
// 		}
// 	}
// 	return string(ns)
// }

// func (list ttList) Key() string {
// 	return list.key//strings.TrimSuffix(GetValidVarName(list.key), "list")
// }

// func (list *ttList) setName(value string) *ttList {
// 	list.key = value
// 	return list
// }
//
// func (list *ttList) addChild(n Node) *ttList {
// 	list.children = append(list.children, n)
// 	return list
// }

type ttListItem struct {
	*ttBase
}

func newListItem(item string) *ttListItem {
	return &ttListItem{ttBase: newBase(LtListItem, item)}
}
type ttTextLine struct {
	*ttBase
}

func newTextLine(line string) *ttTextLine {
	return &ttTextLine{ttBase: newBase(LtTextLine, line)}
}


type ttEmpty struct {
	ttBase
}

type ttComment struct {
	ttBase
}

// create singlton sentinel values once for simply returning a struct
var (
	sEOF     = ttBase{kind: LtEOF}
	sEmpty   = ttEmpty{ttBase{kind: LtEmpty}}
	sComment = ttComment{ttBase{kind: LtComment}}
)

//TODO add Document interface

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
