// Package mdson a package to parse and process the contents of an MDson file
package mdson

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

var debug int

// SetDebug sets the default debug level for all package routines
func SetDebug(level int) {
	debug = level
}

//Parser type for parsing MDSon files and text into a a token tree
type Parser struct {
	ParserOptions
	lineNum     int
	currentLine string
	nextLine    string
	errorState  error
	pendingLeaf Node
	scanner     *bufio.Scanner
}

//NewParser returns an initialized MDsonParser
func NewParser(r io.Reader, options *ParserOptions) (*Parser, error) {
	hp := &Parser{
		ParserOptions: *options,
		scanner:       bufio.NewScanner(r),
	}
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

//Parse parses an MDsonsource
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
	return throw("no valid first-leve (#) heading")
}

func (hp *Parser) advance() Node {
	iFace := hp.parseNextLine()
	hp.lineNum += iFace.lineNum()
	if hp.Debug >= DebugAll {
		switch t := iFace.(type) {
		case *ttKVPair:
		default:
			/*DEBUG*/ hp.log(hp.lineNum, t)
		}
	}
	return iFace
}

//parseBlock return values: false+nil=EOF, false+!nil=error,true+ nil continue
func (hp *Parser) parseBlock(parent *ttBlock) (bool, error) {
	for {
		iFace := hp.advance()
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
		case *ttListItem:
			return false, fmt.Errorf("line %d: '-' outside a list", hp.lineNum)
		case *ttKVPair:
			if parent.isArray() {
				return false, ESyntaxError{hp.lineNum, "only blocks are permitted within a block list"}
			}
			if node.value != "" {
				parent.addChild(node)
				continue
			}
			//value is empty, is this the start of a list or a free text entry
			//store current token to add if if it is not
			_, err := hp.parseListOrFreeText(parent, node)
			if err != nil {
				return false, fmt.Errorf("line %d: %s", hp.lineNum, err)
			}
		default:
			panic(fmt.Sprintf("unhandled token type in parseBlock():line %d", hp.lineNum))
		}
	}
}

func (hp *Parser) parseListOrFreeText(block *ttBlock, openKVP *ttKVPair) (bool, error) {
	nextLine := strings.TrimSpace(hp.nextLine)
	switch {
	case strings.HasPrefix(nextLine, "-"): //nextline starts with '-' must be a list
		list := newList(openKVP.key) //use the key of the open kvp as the list's name
		i := 0
		for strings.HasPrefix(nextLine, "-") {
			i++
			//	fmt.Println("iteration:", i, "current line", hp.currentLine, "next line", hp.nextLine)
			iFace := hp.advance()
			//	fmt.Println("iteration:", i, "current line", hp.currentLine, "next line", hp.nextLine)
			li, ok := iFace.(*ttListItem)
			if !ok { //sanity check
				return false, fmt.Errorf("something went wrong in parsing list items in line %d", hp.lineNum)
			}
			list.addItem(li)
			nextLine = strings.TrimSpace(hp.nextLine)
		}
		block.addChild(list)
	case strings.HasPrefix(nextLine, "<<"): //next line starts with '<<' or '<<<', parse free text
		ls, err := hp.parseFreeText(openKVP)
		if err != nil {
			return false, err
		}
		hp.log(ls)
		block.addChild(ls)
	default: //this was just an empty kvp, add
		block.addChild(openKVP)
	}
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

	lineBreak := ""
	first := true
loop:
	for { //accumulate all text until cTag or an error
		ok := hp.readNextLine()
		// hp.log("PFT()", hp.lineNum, "RNL() returned:", ok, hp.errorState)
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
				//TODO: move to core
				lineBreak = "\n"
				if runtime.GOOS == "windows" {
					lineBreak = "\r\n"
				}
			}
			contents += strings.TrimPrefix(hp.currentLine, oTag) + lineBreak //store rest of first line
			//handle the situation where there is only one text line  ie << one line here >>
			if strings.HasSuffix(hp.currentLine, cTag) {
				contents = strings.TrimSuffix(contents, cTag) + lineBreak
				break loop
			}
			first = false
		case strings.HasSuffix(hp.currentLine, cTag): //the end of the text block
			contents += strings.TrimSuffix(hp.currentLine, cTag) + lineBreak
			break loop
		default: //other lines in between
			contents += hp.currentLine + lineBreak
		}
	}
	return newLiteralString(openKVP.key, contents), nil
}

func (hp *Parser) parseNextLine() Node {
	if !hp.readNextLine() {
		return hp.getErrorStatus()
	}
	line := hp.currentLine
	hp.log("PNL()", hp.lineNum, ":", line)
	trimmed := strings.TrimSpace(line)
	switch {
	case strings.HasPrefix(trimmed, "-"):
		item := line[strings.IndexByte(line, '-'):]
		item = strings.Trim(item, " -")
		return newListItem(item)
	case trimmed == "":
		return &sEmpty
	case strings.HasPrefix(trimmed, "//"):
		return &sComment
	default:
		hd := getHeading(line)
		if hd.level > 0 { //we found a heading
			/*DEBUG*/ hp.log(":", line)

			return newTokenBlock(hd.name).setLevel(hd.level)
		}
		parts := strings.SplitN(line, ":", 2) //split on the first colon
		if len(parts) != 2 {
			//fmt.Println("split", len(parts), parts)
			return newSyntaxError("likely missing ':'")
		}
		//found key:value pair
		key := trimLower(parts[0])
		value := trimLower(parts[1])
		return newKVPair(key, value) //some other key:value pair
	}
}

// readNextLine advances the scanner to the next line and return false
// if EOF encountered or error occured. Parser.Err() reports the specific error
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
		return true
	}
	//there was an error, set parser.errorState
	if hp.scanner.Err() == nil { //eof reached
		hp.errorState = &errEOF{}
	} else { //read error
		hp.errorState = hp.scanner.Err()
	}
	hp.currentLine = hp.nextLine
	return true //stil ok for this iteration
}

//FIXME: replace with one interface arg
func throw(value interface{}) (*ttBlock, error) {
	switch unboxed := value.(type) {
	case string:
		return nil, fmt.Errorf("%s", unboxed)
	case error:
		return nil, fmt.Errorf("%s", unboxed)
	default:
		panic("unsupported argument type in throw()")
	}
}

//move all these into a class in core that can be enclosed in any class
func (hp *Parser) log(a ...interface{}) {
	if hp.Debug >= DebugAll {
		//		fmt.Printf(strings.Repeat("  ", dec.depth))
		fmt.Println(a...)
	}
}

func (hp *Parser) warn(a ...interface{}) {
	if hp.Debug >= DebugWarning {
		fmt.Printf("warning: ")
		fmt.Println(a...)
	}
}

// func (hp *Parser) ParseContents(src string) (tokens []FieldToken, err error) {
// 	return parseContents(src)
// }

//DecodeBlock takes a map of key-value pairs and decodes it into the struct
//passed as an interface
// func (hp *Parser) DecodeBlock(in interface{}, block TokenBlock) error {
// 	return decode(in, block)
// }

//DecodeAttribs takes a string of key-value pairs and decodes it into the struct
//passed as an interface
// func (hp *Parser) DecodeAttribs(in interface{}, attribs string) error {
// 	block, err := parseAttribs(attribs)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse attribs '%s'", attribs)
// 	}
// 	return decode(in, block)
// }
