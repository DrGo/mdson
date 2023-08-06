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


type Context struct {
	Options
	ui.UI
}

func NewContext(options *Options) *Context {
	ctx := &Context{
		Options: *options,
		UI:            ui.NewUI(options.Debug),
	}
	return ctx
}

// Parser type for parsing MDSon files and text into a a token tree
type Parser struct {
	ctx *Context
	lineNum int
	line    string
	// holds the entire document
	doc *Document
	node        Node
	nextNode    Node
	err         error
	scanner     *bufio.Scanner
}

var errEOF = errors.New("end of file")

// NewParser returns an initialized MDsonParser
// FIXME: no need to expose since parser's funcs are not exposed
func NewParser(r io.Reader, ctx  *Context) *Parser {
	p := &Parser{
		ctx: ctx,
		scanner:       bufio.NewScanner(r),
		doc:          newDocument(ctx),
	}
	buf := make([]byte, ctx.BufferCap)
	p.scanner.Buffer(buf, ctx.BufferCap)
	return p
}

// ParseFile parses an MDSon source file into an a
//TODO: merge with parse args filename, io.reader could be nil 
//TODO: add filename field 
func (ctx *Context) ParseFile(fileName string, r io.Reader) (*Document, error) {
	if r == nil {
		f, err := os.Open(fileName)
		if err != nil {
			return throw(err)
		}
		defer f.Close()
		r= f
	}
	p := NewParser(r, ctx)
	err:= p.parse()
	if err != nil {
		return throw(fmt.Errorf("error parsing file '%s': %s", fileName, err))
	}
	err = p.doc.eval()
	if err != nil {
		return throw(fmt.Errorf("error parsing file '%s': %s", fileName, err))
	}
	// p.doc.root = root 
	return p.doc, nil
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
func (p *Parser) parse() error {
	for p.parseBlock(p.doc.root) {
		p.doc.root.AddChild(p.node)
	}
	if p.Err() != nil {
		return p.Err()
	}
	return nil
}

// read the next line and only returns non-comment lines or nil if EOF
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
		n.SetLineNum(p.lineNum)
		return n
	}
	return nil
}

func (p *Parser) advance() bool {
	p.ctx.Log("in advance(): node=", p.node, "nextNode=", p.nextNode) 
	if p.nextNode != nil { //if we already peeked, use that node
		p.node = p.nextNode
		p.nextNode = nil
		return true 
	}
	if n := p.getNextNode(); n != nil {
		p.node = n
		return true
	}
	return false
}

// func (p *Parser) peek() (bool, Node) {
// 	if n := p.getNextNode(); n != nil {
// 		p.nextNode = n //save it for future advance
// 		return true, n
// 	}
// 	return false, nil
// }

// retreat moves the parser to the previous line by
// putting current node in nextNode to be be picked
// next call to advance()
func (p *Parser) retreat() bool {
	p.ctx.Log("putting back node", p.node)
	p.nextNode = p.node
	return true
}

func (p *Parser) setError(err error) bool {
	p.err = err
	return false
}

var ln =0 
// parseBlock return values: false+nil=EOF, false+!nil=error,true+ nil continue
func (p *Parser) parseBlock(parent BlockNode) bool {
	for p.advance() {
		if ln == p.lineNum {
		break }
		ln =p.lineNum
		p.ctx.Log("after parseblock.advance()=>", p.lineNum, p.node)
		//we must have a valid non-comment node
		switch n := p.node.(type) {
		case *ttComment:
		// continue
		case *ttListItem:
			parent.AddChild(newTextLine(n.Key()))
		case *ttTextLine, *ttEmpty:
			p.ctx.Log("inside *ttTextLinei", n.Value())
			parent.AddChild(n)
		case *ttBlock:
			// p.retreat()
			return true 
		// if p.Debug >= DebugAll {
		// 	fmt.Printf(`**parseblock: parent= "%s" [%d]; current %s="%s"[%d]
		// `, parent.Key(), parent.level, "block", n.Key(), n.level)
		// }
		// ok := p.parseBlock(n) //parse it passing this token as a parent
		// if !ok {
		// 	break
		// }
		case *ttList:
			ok := p.parseList(n)
			p.ctx.Log("in *ttlist case after returning from parseList")
			parent.AddChild(n)
			if !ok {
				break
			}
			p.retreat()
		case *ttAttrib:
			// p.Log("inside parseBlock.ttkvpair:", n.key, n.value)
			p.doc.Attribs()[n.key] = n.value
		default:
			panic(fmt.Sprintf("unhandled token type in parseBlock():line %d: %v reflect.type=%s", p.lineNum, n, reflect.TypeOf(n).String()))
		} //switch
	} //for
	return false // advanced returned false
}

func (p *Parser) parseList(list *ttList) bool {
	//we just parsed a list header
	for p.advance(){
		p.ctx.Log("after parseList.advance()=>", p.lineNum, p.node)
		switch p.node.(type) {
		case *ttComment: //ignore
		case *ttListItem:
			list.AddChild(p.node)
		//TODO: allow nested lists
		default:
			// p.retreat()
			return true 
		}
	}
	return false //advance() failed
}

// }
func (p *Parser) parseLine(line string) Node {
	//scenario 1 : empty line
	if line == "" {
		return &sEmpty
	}
	//scenario 2: commented line
	if strings.HasPrefix(line, "//") {
		return &sComment
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
		// p.Log("******************** Block:", line, name, level)
		return newBlock(name, level) //FIXME: change to take name and level
	//keep tilde as another possible marker
	case '~':
		//scenario 7: a list 
		return newList(line[1:], 0)
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
		p.ctx.Log("readLine()", p.lineNum, ":", p.line)
		return true
	}
	//there was an error, set parser.errorState
	if p.scanner.Err() == nil { //eof reached
		p.err = errEOF
	} else { //read error
		p.err = p.scanner.Err()
	}
	return false
}
