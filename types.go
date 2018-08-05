package mdson

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

//An ID describes a field that can be automatically filled with MDSon block name
type ID string

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "mdson: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "mdson: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "mdson: Unmarshal(nil " + e.Type.String() + ")"
}

// An InvalidMarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidMarshalError struct {
	Type reflect.Type
}

func (e *InvalidMarshalError) Error() string {
	if e.Type == nil {
		return "mdson: Marshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "mdson: Marshal(non-pointer " + e.Type.String() + ")"
	}
	return "mdson: Marshal(nil " + e.Type.String() + ")"
}

type errEOF struct{}

func (errEOF) Error() string {
	return "nothing to parse"
}

//TODO: add to ttSyntaxError?!

type ESyntaxError struct {
	lineNum int
	message string
}

func (ese ESyntaxError) Error() string {
	return fmt.Sprintf("line %d syntax error: %s)", ese.lineNum, ese.message)
}

type lineType int

const (
	ltReadError lineType = iota
	ltSyntaxError
	ltEOF
	ltEmpty
	ltComment
	ltList
	ltListItem
	ltBlock
	ltLiteralString
	ltKVPair
)

//lt lineType, linesRead int, key, value string) {
type Node interface {
	lineNum() int
	String() string
	Kind() string
	Name() string
}

//baseToken implements the basic token interface root of all of other tokens
type ttBase struct {
	lnum int
	kind string
	key  string
}

func (bt ttBase) lineNum() int {
	return bt.lnum
}

const nodeDescLine = "type=%s, linesread=%d, key=%s"

func (bt ttBase) String() string {
	return fmt.Sprintf(nodeDescLine, bt.Kind(), bt.lineNum(), bt.key)
}

func (bt ttBase) Kind() string {
	return bt.kind
}

func (bt ttBase) Name() string {
	return bt.key
}

func (bt *ttBase) setLineNum(value int) *ttBase {
	bt.lnum = value
	return bt
}

type ttReadError struct {
	ttBase
	err error
}

func newReadError(value interface{}) *ttReadError {
	re := ttReadError{ttBase: ttBase{kind: "ReadError"}}
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
	se.kind = "syntax error"
	return &se
}

type ttEOF struct{ ttBase }

type ttEmpty struct{ ttBase }

type ttComment struct{ ttBase }

type ttBlock struct {
	ttBase
	level    int
	children []Node
}

func (blk ttBlock) isArray() bool {
	//true if this block functions to hold other blocks
	return strings.HasSuffix(blk.key, " list")
}

func newTokenBlock(name string) *ttBlock {
	return &ttBlock{ttBase: ttBase{kind: "Block", key: name}}
}

func (blk ttBlock) String() string {
	var sb strings.Builder
	sb.Grow(10 * 1024)
	sb.WriteString(blk.ttBase.String() + " " + strconv.Itoa(blk.level) + "\n")
	for _, t := range blk.children {
		sb.WriteString(" " + t.String() + "\n")
	}
	return sb.String()
}

func (blk ttBlock) Name() string {
	if blk.isArray() {
		return strings.TrimSuffix(blk.key, " list")
	}
	return blk.key
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

type ttKVPair struct {
	ttBase
	value string
}

func newKVPair(key, value string) *ttKVPair {
	return &ttKVPair{ttBase: ttBase{kind: "KV Pair", key: key}, value: value}
}

func (kvp ttKVPair) String() string {
	return kvp.ttBase.String() + ": " + kvp.value
}

func (kvp *ttKVPair) setKey(value string) *ttKVPair {
	kvp.key = value
	return kvp
}

func (kvp *ttKVPair) setValue(value string) *ttKVPair {
	kvp.value = value
	return kvp
}

type ttList struct {
	ttBase
	items []*ttListItem
}

func newList(name string) *ttList {
	list := &ttList{ttBase: ttBase{kind: "List", key: name}}
	return list
}

func (list ttList) String() string {
	var sb strings.Builder
	sb.Grow(10 * 1024)
	sb.WriteString(list.ttBase.String())
	for _, li := range list.items {
		sb.WriteString(" " + li.String())
	}
	return sb.String()
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

func newListItem(item string) *ttListItem {
	return &ttListItem{ttBase: ttBase{kind: "ListItem", key: item}}
}

type ttLiteralString struct {
	ttKVPair
}

func newLiteralString(key, value string) *ttLiteralString {
	ls := ttLiteralString{ttKVPair: *newKVPair(key, value)}
	ls.kind = "LiteralString"
	return &ls
}

//create sentinel values once for simply returning a struct
var (
	// sReadError   = ttReadError{ttBase: ttBase{}}
	// sSyntaxError = ttSyntaxError{ttReadError: ttReadError{}}
	sEOF     = ttEOF{ttBase: ttBase{kind: "EOF"}}
	sEmpty   = ttEmpty{ttBase: ttBase{kind: "Empty"}}
	sComment = ttComment{ttBase: ttBase{kind: "Comment"}}
)

//ParserOptions holds parsing options
type ParserOptions struct {
	Debug int
}

//DefaultParserOptions returns reasonable default for parsing
func DefaultParserOptions() *ParserOptions {
	return &ParserOptions{
		Debug: DebugUpdates,
	}
}

//SetDebug sets verbosity level
func (po *ParserOptions) SetDebug(d int) *ParserOptions {
	po.Debug = d
	return po
}
