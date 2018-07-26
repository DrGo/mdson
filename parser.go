// Package mds a package to parse and process the contents of an MDS file
package mds

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

//Parser type for parsing MDS script files and text
type Parser struct {
	ParserOptions
	lineNum     int
	currentLine string
	nextLine    string
	errorState  error
	scanner     *bufio.Scanner
}

//NewParser returns an initialized MDS Parser
func NewParser(r io.Reader, options *ParserOptions) (*Parser, error) {
	hp := &Parser{
		ParserOptions: *options,
		lineNum:       0,
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

//ParseFile parses an MDS script file
func ParseFile(fileName string) (root *ttBlock, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return throw(err)
	}
	hp, err := NewParser(file, DefaultParserOptions())
	if err != nil {
		return throw(fmt.Errorf("error parsing file '%s': %s", fileName, err))
	}
	root, err = hp.parse()
	if err != nil {
		return throw(fmt.Errorf("error parsing file '%s': %s", fileName, err))
	}
	return root, nil
}

//Err return parser error state after last advance() call
func (hp *Parser) Err() error {
	return hp.errorState
}

//Parse parses an MDS source
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
	return root, nil
}

//first time called there is always something to read
func (hp *Parser) readNextLine() bool {
	hp.currentLine = hp.nextLine
	if hp.errorState != nil { //we have reached eof or encountered an error in previous call
		return false
	}
	if hp.scanner.Scan() {
		hp.nextLine = hp.scanner.Text()
		return true
	}
	//there was an error, set parser.errorState
	if hp.scanner.Err() == nil { //eof reached
		hp.errorState = &errEOF{}
	} else { //read error
		hp.errorState = hp.scanner.Err()
	}
	hp.nextLine = ""
	return false
}

func (hp *Parser) advance() token {
	iFace := hp.parseNextLine()
	hp.lineNum += iFace.linesRead()
	if hp.Debug >= DebugAll {
		/*DEBUG*/ fmt.Println(hp.lineNum, iFace)
	}
	return iFace
}

//parseBlock return values: false+nil=EOF, false+!nil=error,true+ nil continue
func (hp *Parser) parseBlock(block *ttBlock) (bool, error) {
	var openKVP token
	for {
		iFace := hp.advance()
		switch token := iFace.(type) {
		case *ttReadError:
			return false, fmt.Errorf("line %d: %s", hp.lineNum, token.err)
		case *ttSyntaxError:
			return false, fmt.Errorf("line %d: %s", hp.lineNum, token.err)
		case *ttComment, *ttEmpty:
			continue
		case *ttEOF:
			if openKVP != nil {
				block.addChild(openKVP) //add the pending token
			}
			return false, nil
		case *ttBlock:
			if openKVP != nil {
				block.addChild(openKVP) //add the pending token
			}
			return hp.parseBlock(block.addChild(token))
		case *ttList:
			if openKVP == nil {
				return false, fmt.Errorf("line %d: '-' outside a list", hp.lineNum)
			}

			openKVP = nil
		case *ttKVPair:
			if token.value != "" {
				if openKVP != nil {
					block.addChild(openKVP) //add the pending token
				}
				block.addChild(token)
				continue
			}
			//value is empty, is this the start of a list or a free text entry
			//store current token to add if if it is not
			openKVP = token
		default:
			panic(fmt.Sprintf("unhandled token type in parseBlock():line %d", hp.lineNum))
		}
	}
}

// func (hp *Parser) parseList(name string, list *ttList) (bool, error) {
// 	var openKVP token
// 	list.setName(name) //first name the list token
// 	for {
// 		iFace := hp.advance()
// 		switch token := iFace.(type) {
// 		case ttList:

// 		}
// 	}
// }

//ParseContents parses contents fields
func (hp *Parser) ParseContents(src string) (tokens []FieldToken, err error) {
	return parseContents(src)
}

//DecodeBlock takes a map of key-value pairs and decodes it into the struct
//passed as an interface
func (hp *Parser) DecodeBlock(in interface{}, block TokenBlock) error {
	return decode(in, block)
}

//DecodeAttribs takes a string of key-value pairs and decodes it into the struct
//passed as an interface
func (hp *Parser) DecodeAttribs(in interface{}, attribs string) error {
	block, err := parseAttribs(attribs)
	if err != nil {
		return fmt.Errorf("failed to parse attribs '%s'", attribs)
	}
	return decode(in, block)
}

func (hp *Parser) parseNextLine() token {
	if !hp.readNextLine() {
		switch hp.Err().(type) {
		case *errEOF: //eof reached
			return &sEOF
		default: //read error, return it
			return sReadError.setError(hp.scanner.Err())
		}
	}
	line := hp.currentLine
	if hp.Debug >= DebugAll {
		fmt.Println(":", line)
	}
	trimmed := strings.TrimSpace(line)
	// if inContentBlock { //in Contents only '+++' has special meaning
	// 	// if strings.HasPrefix(trimmed, "+++") {
	// 	// 	return ltSeparator, 1, "", ""
	// 	// }
	// 	return ltContents, 1, "", line
	// }
	switch {
	case strings.HasPrefix(trimmed, "-"):
		item := line[strings.IndexByte(line, '-'):]
		return newList(item)
	case trimmed == "":
		return &sEmpty
	case strings.HasPrefix(trimmed, "//"):
		return &sComment
	default:
		hd := getHeading(line)
		if hd.level > 0 { //we found a heading
			return newTokenBlock(hd.name).setLevel(hd.level)
		}
		parts := strings.SplitN(line, ":", 2) //split on the first colon
		if len(parts) != 2 {
			//fmt.Println("split", len(parts), parts)
			return sSyntaxError.setError("likely missing ':'")
		}
		//found key:value pair
		key := trimLower(parts[0])
		value := trimLower(parts[1])
		return newKVPair(key, value) //some other key:value pair

		// switch key {
		// case "file", "settings", "document", "section", "header", "footer":
		// 	return ltBlockType, 1, key, v
		// case "contents": //in content the only valid values are ltSeparator or text
		// 	if strings.TrimSpace(value) != "" { // text after Contents is not allowed
		// 		return ltSyntaxError, 1, "", "text after tag 'Contents' is not allowed"
		// 	}
		// 	contents := ""
		// 	linesRead = 1 //we just read one line
		// 	linesInBlock := 0
		// 	for { //accumulate all Contents content until +++ or an error
		// 		lt, linesInBlock, key, value = hp.parseNextLine(scanner, true)
		// 		switch lt {
		// 		case ltSeparator:
		// 			return ltContents, linesRead + 1 /*previously read lines plus this line*/, key, contents
		// 		case ltEOF, ltReadError:
		// 			return ltSyntaxError, linesRead, "", "Contents must end with +++"
		// 		default:
		// 			linesRead += linesInBlock
		// 			contents += value
		// 		}
		// 	}
		// default:

		// }
	}
}

// func (hp *Parser) readFirstLine(scanner *bufio.Scanner) error {
// 	lt, _, _, value := hp.parseNextLine(scanner, false /*outside a content block*/)
// 	switch lt {
// 	case ltSeparator:
// 		return nil
// 	case ltEOF:
// 		return fmt.Errorf("file is empty")
// 	case ltReadError, ltSyntaxError:
// 		return fmt.Errorf("error reading first line: " + value)
// 	default: //anything else is error
// 		return fmt.Errorf("file must start with '+++'")
// 	}
// }

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

type heading struct {
	name  string
	level int
}

func getHeading(line string) heading {
	i := 0
	for ; i < len(line) && line[i] == '#'; i++ { //no utf8 needed, we are only looking for a byte #
	}
	if i == 0 { //# not found, not a heading
		return heading{name: "", level: 0}
	}
	name := trimLower(line[i:])
	if name == "" { //no name, heading but invalid
		return heading{name: "", level: -1}
	}
	return heading{name: name, level: i}
}
