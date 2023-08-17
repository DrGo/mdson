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
	LtCustom
)

func (lt LineType) String() string {
	//[...] creates an array rather than a slice
	if lt >=LtCustom {
		return "Unknown"
	}	
	return [...]string{"Read Error", "Syntax Error", "EOF",
		"Empty", "Comment", "List", "List item", "Block", "Attribute", "Text line",
	"Custom"}[lt]
}

// Node implement parser's AST leaf node
type Node interface {
	String() string
	Kind() LineType
	Key() string
	LineNum() int
	SetLineNum(value int)
	//returns the textual representation of the node as it should appear in a document
	Value() string
	// sets the textual representation of the node as it should appear in a document
	SetValue(s string) Node

	// returns nesting level for a node
	Level() int
	SetLevel(value int) Node
}

type BlockNode interface {
	Node
	Children() []Node
	NthChild(idx int) Node
	AddChild(n Node) BlockNode
	UpdateChild(idx int, n Node) BlockNode
}

// baseToken implements the basic token interface root of all of other tokens
type ttBase struct {
	level    int
	lnum  int
	kind  LineType
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

const nodeDescLine = "type=%s, lineNum=%d, level=%d, key='%s',value='%s'"

func (bt ttBase) String() string {
	return fmt.Sprintf(nodeDescLine, bt.Kind(), bt.LineNum(), bt.Level(), bt.Key(), bt.Value())
}

func (bt ttBase) Kind() LineType {
	return bt.kind
}

func (bt ttBase) Key() string {
	return bt.key //GetValidVarName(bt.key)
}

func (bt ttBase) Value() string {
	return bt.key
}

func (bt *ttBase) SetValue(s string) Node {
	bt.key = s
	return bt
}

func (bt *ttBase) SetLevel(value int) Node {
	if value <0 {
		panic("Level below zero is not permitted")
	}	
	bt.level = value
	return bt
}

func (bt ttBase) Level() int {
	return bt.level
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
	children []Node
}

func newBlock(key string, level int) *ttBlock {
	tb:= &ttBlock{ttBase: newBase(LtBlock, key)}
	tb.SetLevel(level)
	return tb
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



//AddChild adds a child and sets its level to parent.Level + 1
func (blk *ttBlock) AddChild(n Node) BlockNode {
	n.SetLevel(blk.Level() + 1)
	blk.children = append(blk.children, n)
	return blk
}


func (blk *ttBlock) UpdateChild(idx int, n Node) BlockNode {
	blk.children[idx] =n 
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
	for _, t := range list.children {
		sb.WriteString(indent + t.String() + "\n")
	}
	return sb.String()
}


// func GetValidVarName(s string) string {
// 	ns := []rune{}
// 	for _, c := range s {
// 		if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' {
// 			ns = append(ns, c)
// 		}
// 	}
// 	return string(ns)
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
