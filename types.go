package mds

import (
	"fmt"
	"strconv"
	"strings"
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
