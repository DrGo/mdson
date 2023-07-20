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

// Parser type for parsing MDSon files and text into a a token tree
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

// NewParser returns an initialized MDsonParser
// FIXME: no need to expose since parser's funcs are not exposed
func NewParser(r io.Reader, options *ParserOptions) (*Parser, error) {
	p := &Parser{
		ParserOptions: *options,
		scanner:       bufio.NewScanner(r),
		UI:            ui.NewUI(options.Debug),
	}
	// fmt.Println("hp.Debug:", hp.Debug)
	bufCap := 1024 * 1024 //1 megabyte buffer
	buf := make([]byte, bufCap)
	p.scanner.Buffer(buf, bufCap)
	//prime the scanner
	if p.scanner.Scan() {
		p.nextLine = p.scanner.Text()
		return p, nil
	}
	//Scan() failed from the start; either EOF or an error
	if p.scanner.Err() == nil { //eof reached
		return nil, &errEOF{}
	}
	return nil, p.scanner.Err()
}

// ParseFile parses an MDSon source file into an a
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

// Parse parses an MDSon source into an an AST
func Parse(r io.Reader) (root Node, err error) {
	hp, err := NewParser(r, DefaultParserOptions().SetDebug(debug))
	if err != nil {
		return nil, err
	}
	return hp.parse()
}

// Err return parser error state after last advance() call
func (hp *Parser) Err() error {
	return hp.errorState
}

// Parse parses an MDson source
// FIXME: validate block name uniqueness
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
	// if len(root.children) != 1 {
	// 	return throw("there must be only exactly one first-level (#) heading")
	// }
	//discard the root we added above
	// if root, ok := root.children[0].(*ttBlock); ok {
	return root, nil
	// }
	// return throw("no valid first-level (#) heading")
}

// peek parses the nextline without advancing the cursor or changing lineNum
func (p *Parser) peek() Node {
	if p.Debug >= DebugAll {
		p.Log("peek()", p.lineNum, ":", p.nextLine)
	}
	return p.parseLine(p.nextLine)
}

func (p *Parser) advance() Node {
	if !p.readNextLine() {
		return p.getErrorStatus()
	}
	p.Log("advance()", p.lineNum, ":", p.currentLine)
	n := p.parseLine(p.currentLine)
	n.setLineNum(p.lineNum)
	if p.Debug >= DebugAll {
		switch t := n.(type) {
		case *ttKVPair:
		default:
			/*DEBUG*/ p.Log(p.lineNum, t)
		}
	}
	return n
}

// parseBlock return values: false+nil=EOF, false+!nil=error,true+ nil continue
func (p *Parser) parseBlock(parent *ttBlock) (bool, error) {
	for {
		n := p.advance()
		p.Log("after advance()-->line", p.lineNum, n)
		switch n := n.(type) {
		case *ttReadError:
			return false, fmt.Errorf("line %d: %s", p.lineNum, n.err)
		case *ttSyntaxError:
			//TODO: convert to EsynatxError
			return false, fmt.Errorf("line %d: %s", p.lineNum, n.err)
		case *ttComment, *ttEmpty:
			continue
		case *ttEOF:
			return false, nil
		case *ttTextLine:
		case *ttListItem:
			return false, fmt.Errorf("line %d: '-' outside a list", p.lineNum)
		case *ttBlock:
			var ok bool
			var err error
			if p.Debug >= DebugAll {
				fmt.Println("parseblock: parent", parent.Name(), parent.level, "current", n.Name(), n.level)
			}
		loop:
			for { //loop to process all children
				switch {
				// case n.level-parent.level > 1: ////handle nesting errors
				// return false, fmt.Errorf("line %d: block level %d is too deeply nested for its level", hp.lineNum)
				case n.level-parent.level == 1: // this is a child
					ok, err = p.parseBlock(n) //parse it passing this token as a parent
					if err != nil {
						return false, err
					}
					parent.addChild(n)
					if !ok { //FIXME: verify this is correct
						break loop
					}
					if p.pendingLeaf != nil { //do we have unhandled blocken?
						n, _ = p.pendingLeaf.(*ttBlock)
						p.pendingLeaf = nil
						continue
					}
				default: //sibling or higher
					p.pendingLeaf = n //park it in pendingLeaf for ancestors to handle
					return true, nil
				}
			}
		case *ttList:
			_, err := p.parseList(parent, n)
			if err != nil { //error
				return false, fmt.Errorf("line %d: %s", p.lineNum, err)
			}
			// if !ok { //end of file
			// 	return false, nil
			// }
		case *ttKVPair:
			// if parent.isArray() {
			// 	return false, ESyntaxError{hp.lineNum, "key-value pairs not permitted within a list"}
			// }
			p.Log("inside parseBlock.ttkvpair:", n.key, n.value)
			// if n.value != "" {
			parent.addChild(n)
			// continue
			// }
			//value is empty, is this the start of free text entry
			// _, err := hp.parseFreeTextOrKVP(parent, n)
			// if err != nil { //error
			// 	return false, fmt.Errorf("line %d: %s", hp.lineNum, err)
			// }
			// if !ok { //end of file
			// 	return false, nil
			// }
		default:
			panic(fmt.Sprintf("unhandled token type in parseBlock():line %d: %v reflect.type=%s", p.lineNum, n, reflect.TypeOf(n).String()))
		}
	}
}

func (p *Parser) parseList(block *ttBlock, list *ttList) (bool, error) {
	for {
		n := p.advance()
		switch n := n.(type) {
		case *ttListItem:
			list.addItem(n)
		default:
			break
		}
		block.addChild(list)
		return true, nil
	}
}

// func (hp *Parser) parseFreeTextOrKVP(block *ttBlock, openKVP *ttKVPair) (bool, error) {
// 	nextLine := trimLeftSpace(hp.nextLine)
// 	if strings.HasPrefix(nextLine, "<<") { //next line starts with '<<' or '<<<', parse free text
// 		ls, err := hp.parseFreeText(openKVP)
// 		if err != nil {
// 			return false, err
// 		}
// 		hp.Log(ls)
// 		block.addChild(ls)
// 		return true, nil
// 	}
// 	//this was just an empty kvp, add
// 	block.addChild(openKVP)
// 	return true, nil
//
// }

func (p *Parser) getErrorStatus() Node {
	switch p.Err().(type) {
	case *errEOF: //eof reached
		return &sEOF
	default: //read error, return it
		return newReadError(p.scanner.Err())
	}
}

func (p *Parser) isEOF() bool {
	if p.Err() == nil {
		return false
	}
	switch p.Err().(type) {
	case *errEOF: //eof reached
		return true
	default: //read error, return it
		return false
	}
}

// parseFreeText parses free text fields
//
//	func (p *Parser) parseFreeText(openKVP *ttKVPair) (*ttLiteralString, error) {
//		contents := ""
//		oTag, cTag := "<<", ">>" // we are in a block that starts with << or <<<
//		eol := ""
//		first := true
//
// loop:
//
//		for { //accumulate all text until cTag or an error
//			ok := p.readNextLine()
//			// hp.Log("PFT()", hp.lineNum, "RNL() returned:", ok, hp.errorState)
//			switch {
//			case !ok: //error or nothing left to parse
//				if !p.isEOF() {
//					return nil, p.errorState
//				}
//				// end of input, break out of loop
//				break loop
//			case first: //first line in free text block
//				if strings.HasPrefix(p.currentLine, "<<<") { // is this a literal block
//					oTag = "<<<"
//					cTag = ">>>"
//					eol = lineBreak
//				}
//				contents += strings.TrimPrefix(p.currentLine, oTag) + eol //store rest of first line
//				//handle the situation where there is only one text line  ie << one line here >>
//				if strings.HasSuffix(p.currentLine, cTag) {
//					contents = strings.TrimSuffix(contents, cTag) + eol
//					break loop
//				}
//				first = false
//			case strings.HasSuffix(p.currentLine, cTag): //the end of the text block
//				contents += strings.TrimSuffix(p.currentLine, cTag) + eol
//				break loop
//			default: //other lines in between
//				contents += p.currentLine + eol
//			}
//		}
//		return newLiteralString(openKVP.key, contents), nil
//	}
func (p *Parser) parseLine(line string) Node {
	//scenario 1 : empty line
	if line == "" {
		return (&sEmpty)
	}
	//scenario 2: commented line
	if strings.HasPrefix(line, "//") {
		return (&sComment)
	}
	//get the first unicode coding point guaranteed to have >=1 char b/c of the empty check above
	switch ch := []rune(line)[0]; ch {
	//scenario 3: list item
	case '-':
		item := line[1:] //skip the minus
		return newListItem(item)
	//scenario 4: block
	case '#':
		if name, level := getBlockInfo(line); level > 0 {
			p.Log(":", line, name, level)
			return newTokenBlock(name).setLevel(level) //FIXME: change to take name and level
		}
	//keep tilde as another possible marker
	case '.':
		colon := strings.Index(line, ":")
		// scenario 5, regular text starting with dot ( no colon
		if colon == -1 {
			return newTextLine(line)
		}
		// treat as attribute
		parts := strings.SplitN(line[1:], ":", 2) //split on the first colon skipping the first char
		// scenario 6, regular text starting with a dot, 1+ spaces and :
		if strings.TrimSpace(parts[0]) == "" {
			return newTextLine(line)
		}
		//scenario 7: an array since colon is not followed by a value on the same line
		if strings.TrimSpace(parts[1]) == "" {
			return newList(parts[0])
		}
		//scenario 8: attribute; key:value
		return newKVPair(parts[0], parts[1])
	default:
		return newTextLine(line)
	}
	return newTextLine(line)
}

// readNextLine advances the scanner to the next line and return false
// if EOF encountered or error occurred. Parser.Err() reports the specific error
// otherwise it return true
// first time called there is always something to read
func (p *Parser) readNextLine() bool {
	if p.errorState != nil { //we have reached eof or encountered an error in previous call
		p.currentLine = ""
		p.nextLine = ""
		return false
	}
	if p.scanner.Scan() {
		p.currentLine = p.nextLine
		p.nextLine = p.scanner.Text()
		p.lineNum++
		return true
	}
	//there was an error, set parser.errorState
	if p.scanner.Err() == nil { //eof reached
		p.errorState = &errEOF{}
	} else { //read error
		p.errorState = p.scanner.Err()
	}
	p.currentLine = p.nextLine
	p.lineNum++
	return true //still ok for this iteration
}
