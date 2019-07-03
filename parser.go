// Package mdson a package to parse and process the contents of an MDson file
package mdson

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/drgo/core/ui"
)

//Parser type for parsing MDSon files and text into a a token tree
type Parser struct {
	ParserOptions
	ui.UI
	lineNum     int
	currentLine string
	nextLine    string
	errorState  error
	pendingLeaf Node
	scanner     *bufio.Scanner
}

//NewParser returns an initialized MDsonParser
//FIXME: no need to expose since parser's funcs are not exposed
func NewParser(r io.Reader, options *ParserOptions) (*Parser, error) {
	hp := &Parser{
		ParserOptions: *options,
		scanner:       bufio.NewScanner(r),
		UI:            ui.NewUI(options.Debug),
	}
	// fmt.Println("hp.Debug:", hp.Debug)
	bufCap := 1024 * 1024 //1 megabyte buffer
	buf := make([]byte, bufCap)
	hp.scanner.Buffer(buf, bufCap)
	//prime the scanner
	if hp.scanner.Scan() {
		hp.nextLine = hp.scanner.Text()
		return hp, nil
	}
	//error from the start
	if hp.scanner.Err() == nil { //eof reached
		return nil, &errEOF{}
	}
	//read error
	return nil, hp.scanner.Err()
}

//ParseFile parses an MDSon source file into an a
func ParseFile(fileName string) (root Node, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return throw(err)
	}
	root, err = Parse(file)
	if err != nil {
		return throw(fmt.Errorf("error parsing file '%s': %s", fileName, err))
	}
	return root, nil
}

//Parse parses an MDSon source into an an AST
func Parse(r io.Reader) (root Node, err error) {
	hp, err := NewParser(r, DefaultParserOptions().SetDebug(debug))
	if err != nil {
		return nil, err
	}
	return hp.parse()
}

//Err return parser error state after last advance() call
func (hp *Parser) Err() error {
	return hp.errorState
}

//Parse parses an MDson source
//FIXME: validate block name uniqueness
func (hp *Parser) parse() (root *ttBlock, err error) {
	root = newTokenBlock("root")
	for {
		do, err := hp.parseBlock(root)
		if err != nil {
			return throw(err)
		}
		if !do {
			break
		}
	}
	if len(root.children) != 1 {
		return throw("there must be only exactly one first-level (#) heading")
	}
	//discard the root we added above
	if root, ok := root.children[0].(*ttBlock); ok {
		return root, nil
	}
	return throw("no valid first-level (#) heading")
}

func (hp *Parser) advance() Node {
	iFace := hp.parseNextLine()
	iFace.setLineNum(hp.lineNum)
	if hp.Debug >= DebugAll {
		switch t := iFace.(type) {
		case *ttKVPair:
		default:
			/*DEBUG*/ hp.Log(hp.lineNum, t)
		}
	}
	return iFace
}

//parseBlock return values: false+nil=EOF, false+!nil=error,true+ nil continue
func (hp *Parser) parseBlock(parent *ttBlock) (bool, error) {
	for {
		iFace := hp.advance()
		hp.Log("after advance()-->line", hp.lineNum, iFace)
		switch node := iFace.(type) {
		case *ttReadError:
			return false, fmt.Errorf("line %d: %s", hp.lineNum, node.err)
		case *ttSyntaxError:
			//TODO: convert to EsynatxError
			return false, fmt.Errorf("line %d: %s", hp.lineNum, node.err)
		case *ttComment, *ttEmpty:
			continue
		case *ttEOF:
			return false, nil
		case *ttListItem:
			return false, fmt.Errorf("line %d: '-' outside a list", hp.lineNum)
		case *ttBlock:
			var ok bool
			var err error
			if hp.Debug >= DebugAll {
				fmt.Println("parseblock: parent", parent.Name(), parent.level, "current", node.Name(), node.level)
			}
		loop:
			for { //loop to process all children
				switch {
				case node.level-parent.level > 1: ////handle nesting errors
					return false, fmt.Errorf("line %d: block too deeply nested for its level", hp.lineNum)
				case node.level-parent.level == 1: // this is a child
					ok, err = hp.parseBlock(node) //parse it passing this token as a parent
					if err != nil {
						return false, err
					}
					parent.addChild(node)
					if !ok {
						break loop
					}
					if hp.pendingLeaf != nil { //do we have unhandled blocken?
						node, _ = hp.pendingLeaf.(*ttBlock)
						hp.pendingLeaf = nil
						continue
					}
				default: //sibling or higher
					hp.pendingLeaf = node //park it in pendingLeaf for ancestors to handle
					return true, nil
				}
			}
		case *ttList:
			_, err := hp.parseList(parent, node)
			if err != nil { //error
				return false, fmt.Errorf("line %d: %s", hp.lineNum, err)
			}
			// if !ok { //end of file
			// 	return false, nil
			// }
		case *ttKVPair:
			if parent.isArray() {
				return false, ESyntaxError{hp.lineNum, "key-value pairs not permitted within a list"}
			}
			hp.Log("inside parseBlock.ttkvpair:", node.key, node.value)
			if node.value != "" {
				parent.addChild(node)
				continue
			}
			//value is empty, is this the start of free text entry
			_, err := hp.parseFreeTextOrKVP(parent, node)
			if err != nil { //error
				return false, fmt.Errorf("line %d: %s", hp.lineNum, err)
			}
			// if !ok { //end of file
			// 	return false, nil
			// }
		default:
			panic(fmt.Sprintf("unhandled token type in parseBlock():line %d: %v reflect.type=%s", hp.lineNum, node, reflect.TypeOf(iFace).String()))
		}
	}
}

func (hp *Parser) parseList(block *ttBlock, list *ttList) (bool, error) {
	nextLine := trimLeftSpace(hp.nextLine)
	i := 0
	for strings.HasPrefix(nextLine, "-") {
		i++
		iFace := hp.advance()
		// fmt.Println("iteration:", i, "current line", hp.currentLine, "next line", hp.nextLine)
		li, ok := iFace.(*ttListItem)
		if !ok { //just sanity check
			return false, fmt.Errorf("something went wrong in parsing list items in line %d", hp.lineNum)
		}
		list.addItem(li)
		nextLine = trimLeftSpace(hp.nextLine)
	}
	block.addChild(list)
	return true, nil
}

func (hp *Parser) parseFreeTextOrKVP(block *ttBlock, openKVP *ttKVPair) (bool, error) {
	nextLine := trimLeftSpace(hp.nextLine)
	if strings.HasPrefix(nextLine, "<<") { //next line starts with '<<' or '<<<', parse free text
		ls, err := hp.parseFreeText(openKVP)
		if err != nil {
			return false, err
		}
		hp.Log(ls)
		block.addChild(ls)
		return true, nil
	}
	//this was just an empty kvp, add
	block.addChild(openKVP)
	return true, nil

}

func (hp *Parser) getErrorStatus() Node {
	switch hp.Err().(type) {
	case *errEOF: //eof reached
		return &sEOF
	default: //read error, return it
		return newReadError(hp.scanner.Err())
	}
}

func (hp *Parser) isEOF() bool {
	if hp.Err() == nil {
		return false
	}
	switch hp.Err().(type) {
	case *errEOF: //eof reached
		return true
	default: //read error, return it
		return false
	}
}

//parseFreeText parses free text fields
func (hp *Parser) parseFreeText(openKVP *ttKVPair) (*ttLiteralString, error) {
	contents := ""
	oTag, cTag := "<<", ">>" // we are in a block that starts with << or <<<
	eol := ""
	first := true
loop:
	for { //accumulate all text until cTag or an error
		ok := hp.readNextLine()
		// hp.Log("PFT()", hp.lineNum, "RNL() returned:", ok, hp.errorState)
		switch {
		case !ok: //error or nothing left to parse
			if !hp.isEOF() {
				return nil, hp.errorState
			}
			// end of input, break out of loop
			break loop
		case first: //first line in free text block
			if strings.HasPrefix(hp.currentLine, "<<<") { // is this a literal block
				oTag = "<<<"
				cTag = ">>>"
				eol = lineBreak
			}
			contents += strings.TrimPrefix(hp.currentLine, oTag) + eol //store rest of first line
			//handle the situation where there is only one text line  ie << one line here >>
			if strings.HasSuffix(hp.currentLine, cTag) {
				contents = strings.TrimSuffix(contents, cTag) + eol
				break loop
			}
			first = false
		case strings.HasSuffix(hp.currentLine, cTag): //the end of the text block
			contents += strings.TrimSuffix(hp.currentLine, cTag) + eol
			break loop
		default: //other lines in between
			contents += hp.currentLine + eol
		}
	}
	return newLiteralString(openKVP.key, contents), nil
}

func (hp *Parser) parseNextLine() Node {
	if !hp.readNextLine() {
		return hp.getErrorStatus()
	}
	line := hp.currentLine
	hp.Log("parseNextLine()", hp.lineNum, ":", line)
	trimmed := trimLeftSpace(line)
	//scenario 1 : empty line
	if trimmed == "" {
		return (&sEmpty)
	}
	//scenario 2: commented line
	if strings.HasPrefix(trimmed, "//") {
		return (&sComment)
	}
	//scenario 3: list item
	if trimmed[0] == '-' { //guaranteed to have >=1 char b/c of the empty check above
		item := ""
		if len(trimmed) > 1 {
			item = trimmed[1:] //skip the minus
		}
		return newListItem(item)
	}
	//scenario 4: block
	if hd := getHeading(trimmed); hd.level > 0 {
		/*DEBUG*/ hp.Log(":", line, hd)
		return newTokenBlock(hd.name).setLevel(hd.level)
	}
	//scenario 6: invalid key:value pair
	parts := strings.SplitN(line, ":", 2) //split on the first colon
	if len(parts) != 2 {
		return newSyntaxError("likely missing ':'")
	}
	//scenario 7: an array of non-block types
	key := trimLower(parts[0])
	value := parts[1]
	if isArray(key) && strings.TrimSpace(value) == "" {
		return newList(key)
	}
	//scenario 8: valid key:value pair
	return newKVPair(key, value)
}

// readNextLine advances the scanner to the next line and return false
// if EOF encountered or error occurred. Parser.Err() reports the specific error
// otherwise it return true
//first time called there is always something to read
func (hp *Parser) readNextLine() bool {
	if hp.errorState != nil { //we have reached eof or encountered an error in previous call
		hp.currentLine = ""
		hp.nextLine = ""
		return false
	}
	if hp.scanner.Scan() {
		hp.currentLine = hp.nextLine
		hp.nextLine = hp.scanner.Text()
		hp.lineNum++
		return true
	}
	//there was an error, set parser.errorState
	if hp.scanner.Err() == nil { //eof reached
		hp.errorState = &errEOF{}
	} else { //read error
		hp.errorState = hp.scanner.Err()
	}
	hp.currentLine = hp.nextLine
	hp.lineNum++
	return true //still ok for this iteration
}

// //move all these into a class in core that can be enclosed in any class
// func (hp *Parser) log(a ...interface{}) {
// 	if hp.Debug >= DebugAll {
// 		//		fmt.Printf(strings.Repeat("  ", dec.depth))
// 		fmt.Println(a...)
// 	}
// }

// func (hp *Parser) warn(a ...interface{}) {
// 	if hp.Debug >= DebugWarning {
// 		fmt.Printf("warning: ")
// 		fmt.Println(a...)
// 	}
// }
