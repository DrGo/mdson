// Package mds a package to parse and process the contents of an MDS file
package mds

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type lineType int

const (
	ltReadError lineType = iota
	ltSyntaxError
	ltEOF
	ltEmpty
	ltComment
	ltSeparator
	ltBlockType
	ltContents
	ltKVPair
)

//FIXME:  check for out of range
func t2s(lt lineType) string {
	desc := []string{"read error", "syntax error", "eof", "empty", "comment",
		"separator", "block", "contents", "kv pair"}
	return desc[lt]
}

//Parser type for parsing MDS script files and text
type Parser struct {
	ParserOptions
}

//NewParser returns an initialized MDS Parser
func NewParser(options *ParserOptions) *Parser {
	return &Parser{
		ParserOptions: *options,
	}
}

//ParseFile parses an MDS script file
func (hp *Parser) ParseFile(fileName string) (tokens TokenMap, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return throw(err, "")
	}
	tokens, err = hp.Parse(file)
	if err != nil {
		return throw(fmt.Errorf("error parsing file '%s': %s", fileName, err), "")
	}
	return tokens, nil
}

//Parse parses an MDS source
//FIXME: validate block name uniqueness
func (hp *Parser) Parse(r io.Reader) (tokens TokenMap, err error) {
	lineNum := 0
	blockNum := 0
	scanner := bufio.NewScanner(r)
	if err := hp.readFirstLine(scanner); err != nil {
		return throw(err, "")
	}
	//readFirstLine successed, so we read first line and started first block
	lineNum++
	blockNum++
	blockName := ""
	tokens = newTokenMap()
	parseBlock := func() (bool, error) { //false + nil=EOF, false +!nil=error, true+ nil continue
		firstLine := true
	loop:
		for {
			lt, linesRead, key, value := hp.readNextLine(scanner, false /*outside a content block*/)
			lineNum += linesRead
			if hp.Debug >= DebugAll {
				fmt.Println(lineNum, t2s(lt), "lines-read", linesRead, key, value) //DEBUG
			}
			if firstLine {
				switch lt {
				case ltEOF:
					return false, nil //no more data
				case ltSyntaxError, ltReadError:
					return false, fmt.Errorf("line %d: %s", lineNum, value)
				case ltEmpty, ltComment:
					continue //skip
				case ltBlockType:
					blockName = value
					//TODO: check for blockname uniqueness
					tokens.addEntry(blockName, "name", blockName)
					tokens.addEntry(blockName, "type", key)
					tokens.addEntry(blockName, "order", strconv.Itoa(blockNum))
					firstLine = false
					continue
				default:
					return false, fmt.Errorf("line %d: expected block type declaration, found '%s'", lineNum, key)
				}
			}
			switch lt {
			case ltBlockType:
				return false, fmt.Errorf("line %d: block type declaration '%s' not allowed", lineNum, key)
			case ltEOF:
				if blockName != "" {
					return false, fmt.Errorf("line %d:+++ not found at end of section %s", lineNum, blockName)
				}
				return false, nil
			case ltReadError, ltSyntaxError:
				return false, fmt.Errorf("line %d: %s", lineNum, value)
			case ltKVPair:
				tokens.addEntry(blockName, key, value)
			case ltContents:
				tokens.addEntry(blockName, "contents", value) //ltContents also indicate end of block
				fallthrough
			case ltSeparator:
				blockName = "" //no longer in this section
				blockNum++
				break loop
			case ltComment, ltEmpty:
				continue
			default:
				panic(fmt.Sprintf("unhandled token type in parseBlock():line %d-tt %d,%s,%s",
					lineNum, lt, key, value))
			}

		}
		return true, nil
	}
	for {
		do, err := parseBlock()
		if err != nil {
			return throw(err, "")
		}
		if !do {
			break
		}
	}
	return tokens, nil
}

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

func (hp *Parser) readNextLine(scanner *bufio.Scanner, inContentBlock bool) (lt lineType, linesRead int, key, value string) {
	if !scanner.Scan() {
		switch scanner.Err() {
		case nil: //eof reached
			return ltEOF, 0, "", ""
		default: //read error, return its description
			return ltReadError, 0, "", scanner.Err().Error()
		}
	}
	line := scanner.Text()
	if hp.Debug >= DebugAll {
		fmt.Println(":", line)
	}
	trimmed := strings.TrimSpace(line)
	if inContentBlock { //in Contents only '+++' has special meaning
		if strings.HasPrefix(trimmed, "+++") {
			return ltSeparator, 1, "", ""
		}
		return ltContents, 1, "", line
	}

	switch {
	case trimmed == "":
		return ltEmpty, 1, "", ""
	case strings.HasPrefix(trimmed, "//"):
		return ltComment, 1, "", line
	case strings.HasPrefix(trimmed, "+++"):
		return ltSeparator, 1, "", ""
	default:
		parts := strings.SplitN(line, ":", 2) //split on the first colon
		if len(parts) != 2 {
			fmt.Println("split", len(parts), parts)
			return ltSyntaxError, 1, "", "likely missing ':'"
		}
		//found key:value pair
		key = trimLower(parts[0])
		value = trimLower(parts[1])
		switch key {
		case "file", "settings", "document", "section", "header", "footer":
			return ltBlockType, 1, key, value
		case "contents": //in content the only valid values are ltSeparator or text
			if strings.TrimSpace(value) != "" { // text after Contents is not allowed
				return ltSyntaxError, 1, "", "text after tag 'Contents' is not allowed"
			}
			contents := ""
			linesRead = 1 //we just read one line
			linesInBlock := 0
			for { //accumulate all Contents content until +++ or an error
				lt, linesInBlock, key, value = hp.readNextLine(scanner, true)
				switch lt {
				case ltSeparator:
					return ltContents, linesRead + 1 /*previously read lines plus this line*/, key, contents
				case ltEOF, ltReadError:
					return ltSyntaxError, linesRead, "", "Contents must end with +++"
				default:
					linesRead += linesInBlock
					contents += value
				}
			}
		default:
			return ltKVPair, 1, key, value //some other key:value pair
		}
	}
}

func (hp *Parser) readFirstLine(scanner *bufio.Scanner) error {
	lt, _, _, value := hp.readNextLine(scanner, false /*outside a content block*/)
	switch lt {
	case ltSeparator:
		return nil
	case ltEOF:
		return fmt.Errorf("file is empty")
	case ltReadError, ltSyntaxError:
		return fmt.Errorf("error reading first line: " + value)
	default: //anything else is error
		return fmt.Errorf("file must start with '+++'")
	}
}

func throw(err error, msg string) (TokenMap, error) {
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return nil, fmt.Errorf("%s", msg)
}
