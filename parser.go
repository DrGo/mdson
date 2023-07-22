// Package mdson a package to parse and process the contents of an MDson file
package mdson

import (
	"bufio"
	"errors"
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
	line string
	// nextLine    string
	node  Node  
	nextNode 	Node 
	err  error
	pendingLeaf Node
	scanner     *bufio.Scanner
}

var errEOF = errors.New("end of file")

// NewParser returns an initialized MDsonParser
// FIXME: no need to expose since parser's funcs are not exposed
func NewParser(r io.Reader, options *ParserOptions) (*Parser) {
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
	// if p.scanner.Scan() {
	// 	p.nextLine = p.scanner.Text()
	// 	return p, nil
	// }
	// //we are here so Scan() failed from the start; either EOF or an error
	// if p.scanner.Err() == nil { //eof reached
	// 	return nil, errEOF
	// }
	return p  //nil, p.scanner.Err()
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
	p := NewParser(r, DefaultParserOptions().SetDebug(debug))
	// if err != nil {
	// 	return nil, err
	// }
	return p.parse()
}

// Err return parser error state after last advance() call
func (p *Parser) Err() error {	
	if p.err == errEOF {
		return nil   
	}	
	return p.err 
}

// Parse parses an MDson source
// FIXME: validate block name uniqueness
func (p *Parser) parse() (root *ttBlock, err error) {
	root = newTokenBlock("root", 0)
	for p.parseBlock(root){}
	if p.Err()!= nil {
		return nil, p.Err()
	}	
	return root, nil
}


// read the next line and only returns non-comment lines  
func (p *Parser) getNextNode() Node {
	//TODO: verify error propagation is working
	for p.readLine() {
			n := p.parseLine(p.line)
			switch n.Kind() {
			case LtComment:
				continue
			case LtEOF:
				return nil
			}
			n.setLineNum(p.lineNum)
			return n 
		}
	return nil
}

func (p *Parser) advance() bool {
	if p.nextNode != nil { //if we already peeked, use that node 
		p.node = p.nextNode
		p.nextNode = nil 
	}
	if n:= p.getNextNode(); n != nil {
				p.node = n
		return true 
	}
	return false 
}

func (p *Parser) peek() (bool, Node) {
	if n:= p.getNextNode(); n != nil {
		p.nextNode = n //save it for future advance
		return true, n 
	}
	return false, nil 
}

 func (p *Parser) setError(err error) bool {
 	p.err= err
	return false 
 }

// parseBlock return values: false+nil=EOF, false+!nil=error,true+ nil continue
func (p *Parser) parseBlock(parent *ttBlock) bool {
	for p.advance() {
		// p.Log("after p.next(naPeek): ", p.lineNum, n)
		//we must have a valid non-comment node
		switch n := p.node.(type) {
		case *ttComment:
		// continue
		case *ttListItem:
			parent.addChild(newTextLine(n.Key()))
		case *ttTextLine, *ttEmpty:
			parent.addChild(n)
		case *ttBlock:
			if p.Debug >= DebugAll {
				fmt.Printf("**parseblock: parent= %s [%d]; current %s=%s[%d]\n", parent.Key(), parent.level, "block", n.Key(), n.level)
			}
			if n.level > parent.level { // this is a child block parse it				
				ok := p.parseBlock(n) //parse it passing this token as a parent
				parent.addChild(n)
				if !ok {
					return false 
				}	
			}
			// if this is a sibling, we let caller handle it
			return true
		case *ttList:
			ok:= p.parseList(n)
			n.level = parent.level + 1
			parent.addChild(n)
			fmt.Printf("**parseblock: parent= %s [%d]; current %s=%s[%d]\n", parent.Key(), parent.level, "list", n.Key(), n.level)
			if !ok { 
				return ok
			}
		case *ttAttrib:
			p.Log("inside parseBlock.ttkvpair:", n.key, n.value)
			parent.attribs[n.key] = n.value
		default:
			panic(fmt.Sprintf("unhandled token type in parseBlock():line %d: %v reflect.type=%s", p.lineNum, n, reflect.TypeOf(n).String()))
		} //switch
	} //for
	return false // advanced returned false 
}

func (p *Parser) parseList(list *ttList) bool{
loop:
	for {
		ok, n := p.peek() 
		if !ok {
			return ok
		}
		p.Log("inside parseList=>", n)
		switch n.(type) {
		case *ttComment: //ignore
		case *ttListItem:
			p.advance()
			list.addChild(p.node)
			continue
		//TODO: allow nested lists
		default:
			break loop
		}
	}
	return true
}

// }
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
		name, level := getBlockInfo(line)
		p.Log("******************** Block:", line, name, level)
		return newTokenBlock(name, level) //FIXME: change to take name and level
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
			return newList(parts[0], 0)
		}
		//scenario 8: attribute; key:value
		return newAttrib(parts[0], parts[1])
	default:
		return newTextLine(line)
	}
}

// readNextLine advances the scanner to the next line and return false
// if EOF encountered or error occurred. Parser.Err() reports the specific error
// otherwise it return true
// first time called there is always something to read
func (p *Parser) readLine() bool {
	if p.err != nil { //we have reached eof or encountered an error in previous call
		return false
	}
	if p.scanner.Scan() {
		p.line = p.scanner.Text()
		p.lineNum++
		p.Log("readLine()", p.lineNum, ":", p.line)
		return true
	}
	//there was an error, set parser.errorState
	if p.scanner.Err() == nil { //eof reached
		p.err= errEOF
	} else { //read error
		p.err= p.scanner.Err()
	}
	return false 
}
