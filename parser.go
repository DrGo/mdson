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
	currentLine string
	nextLine    string
	errorState  error
	pendingLeaf Node
	scanner     *bufio.Scanner
}

var errEOF= errors.New("end of file")
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
	//we are here so Scan() failed from the start; either EOF or an error
	if p.scanner.Err() == nil { //eof reached
		return nil, errEOF
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
	p, err := NewParser(r, DefaultParserOptions().SetDebug(debug))
	if err != nil {
		return nil, err
	}
	return p.parse()
}

// Err return parser error state after last advance() call
func (p *Parser) Err() error {
	return p.errorState
}

// Parse parses an MDson source
// FIXME: validate block name uniqueness
func (p *Parser) parse() (root *ttBlock, err error) {
	root = newTokenBlock("root", 0)
	for {
		do, err := p.parseBlock(root)
		if err != nil && err!=errEOF {
			return throw(err)
		}
		if !do {
			break
		}
	}
	return root, nil
}

type nextAct int

const (
	// parses the next line but does not move the cursor
	naPeek nextAct = iota
	// only move the cursor to next line; no parsing
	naNext
)

// peek parses the nextline without advancing the cursor or changing lineNum
func (p *Parser) next(act nextAct) (ok bool, n Node) {
	if p.Debug >= DebugAll {
		p.Log("next()", p.lineNum, ":", p.nextLine)
	}
	for p.readNextLine() {
		switch act {
		case naPeek:
			n := p.parseLine(p.nextLine)
			switch n.Kind() {
			case LtComment:
				continue
			case LtEOF:
				return false, nil   
			}	
			n.setLineNum(p.lineNum)
			return true, n
		case naNext:
			return true, nil
		default:
			//TODO: panic
			return false, nil
		}
	}
	//readNextLine is false
	return false, nil
}

// parseBlock return values: false+nil=EOF, false+!nil=error,true+ nil continue
func (p *Parser) parseBlock(parent *ttBlock) (ok bool, err error) {
	for {
		// peek at next line; cannot fail
		_, n := p.next(naPeek)
		// p.Log("after advance()-->line", p.lineNum, n)
		switch n := n.(type) {
		case *ttComment:
		// continue
		case *ttListItem:
			return false, fmt.Errorf("line %d: '-' outside a list", p.lineNum)
		case *ttTextLine, *ttEmpty:
			parent.addChild(n)
		case *ttBlock:
			if p.Debug >= DebugAll {
				fmt.Println("************parseblock: parent", parent.Key(), parent.level, "current", n.Key(), n.level)
			}
			if n.level > parent.level { // this is a child block parse it
				ok, err = p.parseBlock(n) //parse it passing this token as a parent
				parent.addChild(n)
				return ok, err
			}
			// if this is a sibling, we let caller handle it
			return true, nil
		case *ttList:
			ok, err = p.parseList(n)
			parent.addChild(n)
			if !ok || err != nil { //error
				return ok, fmt.Errorf("line %d: %s", p.lineNum, err)
			}
		case *ttAttrib:
			p.Log("inside parseBlock.ttkvpair:", n.key, n.value)
			parent.attribs[n.key] = n.value
		default:
			panic(fmt.Sprintf("unhandled token type in parseBlock():line %d: %v reflect.type=%s", p.lineNum, n, reflect.TypeOf(n).String()))
		} //switch
		//TODO: move error handling to p.next
		ok,_ = p.next(naNext) //handled next line so move to next
		//FIXME: return more error info
		if !ok {
			return ok, p.errorState
		}
	} //for
}

func (p *Parser) parseList(list *ttList) (bool, error) {
	loop:
	for {
		_, n := p.next(naPeek) //cannot fail		
		switch n := n.(type) {
		case *ttComment:			
		case *ttListItem:
			list.addItem(n)
		//TODO: allow nested lists
		default:
			break loop
		}
		ok, _ := p.next(naNext)
		if !ok {
			return ok, p.errorState
		}
	}
	return true, nil
}

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
			return newList(parts[0])
		}
		//scenario 8: attribute; key:value
		return newAttrib(parts[0], parts[1])
	default:
		return newTextLine(line)
	}
	// return newTextLine(line)
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
		p.errorState = errEOF
	} else { //read error
		p.errorState = p.scanner.Err()
	}
	p.currentLine = p.nextLine
	p.lineNum++
	return true //still ok for this iteration
}
