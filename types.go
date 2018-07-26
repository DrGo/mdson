package mds

import (
	"fmt"
	"strconv"
	"strings"
)

type errEOF struct{}

func (e *errEOF) Error() string {
	return "nothing to parse"
}

type lineType int

const (
	ltReadError lineType = iota
	ltSyntaxError
	ltEOF
	ltEmpty
	ltComment
	ltList
	ltBlock
	ltContents
	ltKVPair
)

//lt lineType, linesRead int, key, value string) {
type token interface {
	linesRead() int
	String() string
}

//baseToken implements the token interface
type ttBase struct {
	lines int
}

// func (bt *ttBase) kind() lineType {
// 	return bt.lt
// }

func (bt *ttBase) linesRead() int {
	return bt.lines
}

func (bt *ttBase) String() string {
	return ""
}

func (bt *ttBase) setLinesRead(value int) *ttBase {
	bt.lines = value
	return bt
}

// var ltDescList = []string{"read error", "syntax error", "eof", "empty", "comment",
// 	"list", "block", "contents", "kv pair"}

type (
	ttReadError struct {
		ttBase
		err error
	}
	ttSyntaxError struct{ ttReadError }
	ttEOF         struct{ ttBase }
	ttEmpty       struct{ ttBase }
	ttComment     struct{ ttBase }
	ttList        struct {
		ttBase
		name  string
		items []string
	}
	ttBlock struct {
		ttBase
		name     string
		level    int
		children []token
	}
	ttContents struct{ ttBase }
	ttKVPair   struct {
		ttBase
		key   string
		value string
	}
)

func (re *ttReadError) setError(value interface{}) *ttReadError {
	switch unboxed := value.(type) {
	case string:
		re.err = fmt.Errorf("%s", unboxed)
	case error:
		re.err = unboxed
	default:
		panic("unsupported argument type in ttReadError.setError()")
	}
	return re
}

func newTokenBlock(name string) *ttBlock {
	return &ttBlock{ttBase: ttBase{}, name: name}
}

func (blk *ttBlock) setName(value string) *ttBlock {
	blk.name = value
	return blk
}

func (blk *ttBlock) setLevel(value int) *ttBlock {
	blk.level = value
	return blk
}

func (blk *ttBlock) addChild(t token) *ttBlock {
	blk.children = append(blk.children, t)
	return blk
}

func newKVPair(key, value string) *ttKVPair {
	return &ttKVPair{ttBase: ttBase{lines: 1}, key: key, value: value}
}

func (kvp *ttKVPair) setKey(value string) *ttKVPair {
	kvp.key = value
	return kvp
}

func (kvp *ttKVPair) setValue(value string) *ttKVPair {
	kvp.value = value
	return kvp
}

func newList(item string) *ttList {
	list := &ttList{ttBase: ttBase{lines: 1}}
	list.items = append(list.items, item)
	return list
}

func (list *ttList) setName(value string) *ttList {
	list.name = value
	return list
}

const tokDescLine = "type=%s, linesread=%d, key=%s value=%s\n"

func (re ttReadError) String() string {
	return fmt.Sprintf(tokDescLine, "read error", re.linesRead(), "err", re.err)
}

func (se ttSyntaxError) String() string {
	return fmt.Sprintf(tokDescLine, "syntax error", se.linesRead(), "err", se.err)
}

func (ttEOF) String() string {
	return fmt.Sprintf(tokDescLine, "EOF", 0, "", "")
}

func (em ttEmpty) String() string {
	return fmt.Sprintf(tokDescLine, "empty", em.linesRead(), "", "")
}

func (co ttComment) String() string {
	return fmt.Sprintf(tokDescLine, "comment", co.linesRead(), "", "")
}

func (blk ttBlock) String() string {
	var sb strings.Builder
	sb.Grow(10 * 1024)
	sb.WriteString(fmt.Sprintf(tokDescLine, "block", blk.linesRead(), blk.name, strconv.Itoa(blk.level)))
	for _, t := range blk.children {
		sb.WriteString(t.String())
	}
	return sb.String()
}

func (kvp ttKVPair) String() string {
	return fmt.Sprintf(tokDescLine, "KV pair", kvp.linesRead(), kvp.key, kvp.value)
}

//create sentinel values to simply returning
var (
	sReadError   = ttReadError{ttBase: ttBase{}}
	sSyntaxError = ttSyntaxError{ttReadError: ttReadError{}}
	sEOF         = ttEOF{ttBase: ttBase{}}
	sBase        = ttBase{lines: 1}
	sEmpty       = ttEmpty{ttBase: sBase}
	sComment     = ttComment{ttBase: sBase}
	sList        = ttList{ttBase: ttBase{}}
	//	sBlock       = ttBlock{ttBase: ttBase{}}
	sContents = ttContents{ttBase: ttBase{}}

//	sKVPair      = ttKVPair{ttBase: sBase}
)

type (
	//TokenBlock a map of key to value entries
	TokenBlock map[string]string
	//TokenMap a map of blocked parsed from an htds script file
	TokenMap map[string]TokenBlock
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

func (tb TokenBlock) String() string {
	var sb strings.Builder
	sb.Grow(5 * 1024)
	/*DEBUG*/ //sb.WriteString(fmt.Sprintf("%s: [%d entries]\n", blockName, len(block)))
	for k, v := range tb {
		sb.WriteString(fmt.Sprintf("%s:%s\n", k, v))
	}
	return sb.String()
}

func newTokenMap() TokenMap {
	return make(TokenMap, 10)
}

//OrderedBlockList returns a list of blocks in their order in the script
func (tl TokenMap) OrderedBlockList() []TokenBlock {
	var lst = make([]TokenBlock, len(tl))
	for _, block := range tl {
		order, _ := strconv.Atoi(block["order"])
		lst[order-1] = block
	}
	return lst
}

//BlockNames returns a slice of block names
func (tl TokenMap) BlockNames() []string {
	var lst = make([]string, len(tl))
	for blockName, block := range tl {
		order, _ := strconv.Atoi(block["order"])
		lst[order-1] = blockName
	}
	return lst
}
func (tl TokenMap) String() string {
	var sb strings.Builder
	sb.Grow(10 * 1024)
	//sb.WriteString(fmt.Sprintf("list size:%d\n", len(tl)))
	sb.WriteString("+++")
	lst := tl.BlockNames()
	//fmt.Println("Block names:", lst)
	for _, blockName := range lst {
		block := tl[blockName]
		sb.WriteString(block.String())
		sb.WriteString("+++")
	}
	return sb.String()
}

func (tl TokenMap) addEntry(blockName, key, value string) error {
	blockName = trimLower(blockName)
	block := tl[blockName]
	if block == nil {
		block = make(map[string]string, 10)
		tl[blockName] = block
	}
	block[trimLower(key)] = trimLower(value)
	/*DEBUG*/ // fmt.Println("addEntry()", blockName, "-", key, ":", tl[blockName][key])
	return nil
}

// type node interface{}

// type block struct {
// 	name     string
// 	level    int
// 	children []*node
// }

// type list struct {
// 	children []*node
// }

// type kv struct {
// 	key   string
// 	value string
// }
