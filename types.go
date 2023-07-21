package mdson

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

//An ID describes a field that can be automatically filled with MDSon block name
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
	LtLiteralString
	LtAttrib
	LtTextLine
)

//Node implement parser's AST node
type Node interface {
	String() string
	Kind() LineType 
	Name() string
	Children() []Node
	ChildByName(name string) Node
	ValueOf() map[string]string
	lineNum() int
	setLineNum(value int)
	Value() string
}

//baseToken implements the basic token interface root of all of other tokens
type ttBase struct {
	lnum int
	kind LineType
	key  string
}

func (bt ttBase) lineNum() int {
	return bt.lnum
}

const nodeDescLine = "type=%d, lineNum=%d, key=%s"

func (bt ttBase) String() string {
	return fmt.Sprintf(nodeDescLine, bt.Kind(), bt.lineNum(), bt.key)
}

func (bt ttBase) Kind() LineType {
	return bt.kind
}

func (bt ttBase) Name() string {
	return GetValidVarName(bt.key)
}

func (bt *ttBase) setLineNum(value int) {
	bt.lnum = value
}

func (bt ttBase) isArray() bool {
	return isArray(bt.key)
}

func (bt ttBase) Children() []Node {
	return nil
}

func (bt ttBase) ChildByName(name string) Node {
	return nil
}

func (bt ttBase) ValueOf() map[string]string {
	return nil
}

func (bt ttBase) Value() string {
	return bt.key
}

type ttReadError struct {
	ttBase
	err error
}

func newReadError(value interface{}) *ttReadError {
	re := ttReadError{ttBase: ttBase{kind: LtReadError}}
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
	ttBase
	level    int
	children []Node
	attribs map[string]string 
}

func newTokenBlock(name string) *ttBlock {
	return &ttBlock{ttBase: ttBase{kind: LtBlock, key: name},
		attribs: make(map[string]string)}
}

func (blk ttBlock) String() string {
	var sb strings.Builder
	sb.Grow(10 * 1024)
	sb.WriteString(blk.ttBase.String() + " " + strconv.Itoa(blk.level) + "\n")
	for k, v := range blk.attribs {
		sb.WriteString(k)
		sb.WriteString(" = ")
		sb.WriteString(v)
		sb.WriteRune('\n')
	}
	for _, t := range blk.children {
		sb.WriteString(" " + t.String() + "\n")
	}
	return sb.String()
}

func (blk ttBlock) Name() string {
	name := GetValidVarName(blk.key)
	if blk.isArray() {
		return strings.TrimSuffix(name, "list")
	}
	return name
}

func (blk *ttBlock) setName(value string) *ttBlock {
	blk.key = value
	return blk
}

func (blk *ttBlock) setLevel(value int) *ttBlock {
	blk.level = value
	return blk
}

func (blk *ttBlock) addChild(t Node) *ttBlock {
	blk.children = append(blk.children, t)
	return blk
}

func (blk ttBlock) getChildByName(name string) Node {
	for _, c := range blk.children {
		if c.Name() == name {
			return c
		}
	}
	return nil
}

func (blk ttBlock) Children() []Node {
	return blk.children
}

func (blk ttBlock) ChildByName(name string) Node {
	return blk.getChildByName(name)
}

func (blk ttBlock) ValueOf() map[string]string {
	contents := map[string]string{}
	for _, c := range blk.children {
		switch uc := c.(type) {
		case *ttAttrib:
			contents[uc.key] = uc.value
		// case *ttTextLine:
		// 	contents[uc.key] = uc.value
		}
	}
	return contents
}

type ttAttrib struct {
	ttBase
	value string
}

func newAttrib(key, value string) *ttAttrib {
	return &ttAttrib{ttBase: ttBase{kind:LtAttrib, key: strings.TrimSpace( key)},
		value: strings.TrimSpace(value)}
}

func (kvp ttAttrib) String() string {
	return kvp.ttBase.String() + ": " + kvp.value
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
	ttBase
	items []*ttListItem
}

func newList(name string) *ttList {
	list := &ttList{ttBase: ttBase{kind: LtList, key: name}}
	return list
}

func (list ttList) String() string {
	var sb strings.Builder
	sb.Grow(10 * 1024)
	sb.WriteString(list.ttBase.String())
	sb.WriteRune('\n')
	for _, li := range list.items {
		sb.WriteString(li.String())
	sb.WriteRune('\n')
	}
	return sb.String()
}



func GetValidVarName(s string) string {
	ns := []rune{}
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' {
			ns = append(ns, c)
		}
	}
	return string(ns)
}

func (list ttList) Name() string {
	return strings.TrimSuffix(GetValidVarName(list.key), "list")
}

func (list *ttList) setName(value string) *ttList {
	list.key = value
	return list
}

func (list *ttList) addItem(li *ttListItem) *ttList {
	list.items = append(list.items, li)
	return list
}

type ttListItem struct {
	ttBase
}

type ttTextLine struct {
	ttBase
}


func newTextLine(item string) *ttTextLine {
	return &ttTextLine{ttBase: ttBase{kind:LtTextLine, key: item}}
}

func newListItem(item string) *ttListItem {
	return &ttListItem{ttBase: ttBase{kind: LtListItem, key: item}}
}

type ttEmpty struct {
	ttBase
}

type ttComment struct {
	ttBase
}
//create singlton sentinel values once for simply returning a struct
var (
	sEOF     = ttBase{kind: LtEOF}
	sEmpty   = ttEmpty{ttBase{kind: LtEmpty}}
	sComment = ttComment{ttBase{kind: LtComment}}
)
