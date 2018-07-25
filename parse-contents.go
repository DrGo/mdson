package mds

import (
	"fmt"
	"strings"
)

//FieldTokenType just an int
type FieldTokenType int

const (
	//TtError read or parsing error
	TtError FieldTokenType = iota
	// TtText ordinary text
	TtText
	TtHtmldocx
	TtWord
	TtXML
	TtInclude
	TtScript
)

const emptySourceError = "nothing to parse"

var fieldTokenTypeMap = map[string]FieldTokenType{
	"text":     TtText,
	"htmldocx": TtHtmldocx,
	"word":     TtWord,
	"xml":      TtXML,
	"include":  TtInclude,
	"script":   TtScript,
}

type FieldToken struct {
	kind  FieldTokenType
	value string
}

func newfieldToken(kind FieldTokenType, value string) FieldToken {
	return FieldToken{
		kind:  kind,
		value: value,
	}
}

//parseContents breaks down a string into normal text and field tokens
//eg parseContents("a${b}c")--> []fieldToken{newfieldToken(ttText, "a"), newfieldToken(field of some kind, "b"), newfieldToken(ttText, "c")}
func parseContents(src string) (tokens []FieldToken, err error) {
	src = strings.TrimSpace(src)
	if src == "" {
		return nil, fmt.Errorf(emptySourceError)
	}
	const null = '\x00'
	buf := []byte(src) //convert to byte array and append null to allow the algorithm below to work
	buf = append(buf, null)
	bufLen := len(buf)
	i := 0 // is the source byte index
	j := 0 //is the byte that started the current token (the anchor)
	next := func() byte {
		if i+1 < bufLen {
			return buf[i+1]
		}
		return null
	}
	open := false               //tracks if we are in a field ie after ${
	for ; buf[i] != null; i++ { //parsing as bytes because all the characters we are looking for are also bytes in utf8
		//fmt.Printf("i=%d j=%d, md[i]=%c\n", i, j, buf[i]) //DEBUG
		switch buf[i] {
		case '$': //is this the beginning of a field
			switch next() {
			case '{': //we are in a field
				if open {
					return nil, fmt.Errorf("nested fields not allowed: error around char %d", i)
				}
				open = true
				if i > j { // if we have text before the ${, create a token for it
					tokens = append(tokens, FieldToken{kind: TtText, value: string(buf[j:i])})
				}
				//fmt.Printf("i=%d j=%d, md[i]=%c :in {\n", i, j, buf[i]) //DEBUG
				i++       //skip {
				j = i + 1 //move the anchor to the position after {
			}
		case '}':
			if !open { //we are not a field
				return nil, fmt.Errorf("} not allowed outside a field: error around char %d", i)
			}
			open = false
			if i <= j { //empty {} not allowed
				return nil, fmt.Errorf("empty ${} not allowed: error around char %d", i)
			}
			fld := string(buf[j:i])
			parts := strings.SplitN(fld, ":", 2) //split on the first colon
			if len(parts) < 2 {
				return nil, fmt.Errorf("syntax error in field %s, likely missing ':'", fld)
			}
			//fmt.Printf("key:%s, value:%s\n", parts[0], parts[1]) //DEBUG
			key := fieldTokenTypeMap[strings.TrimSpace(parts[0])]
			value := parts[1] //there must be exactly two parts
			switch key {
			case TtInclude, TtScript:
				trimmed := strings.TrimSpace(value)
				if trimmed == "" {
					return nil, fmt.Errorf("syntax error in field %s: missing value", fld)
				}
				value = trimmed
				fallthrough
			case TtHtmldocx, TtWord, TtXML:
				tokens = append(tokens, FieldToken{kind: key, value: value})
			default:
				return nil, fmt.Errorf("syntax error in field '%s': unknown key", fld)
			}
			//fmt.Printf("i=%d j=%d, md[i]=%c :in }\n ", i, j, buf[i]) //DEBUG
			j = i + 1 //move the anchor to the next position
		case '\\': //handle escaping of ${ or }
			switch c := next(); c {
			case '\\', '$', '}':
				buf = append(buf[:i], buf[i+1:]...) //excise \ and ignore the escaped character
				//TODO: handle line breaks \n ??
				//TODO: handle line breaks \t --> create tab nodes?
			}
		}
	}
	if open {
		return nil, fmt.Errorf("missing '}'")
	}
	if i > j { //capture any remaining text other than the null byte we added above
		tokens = append(tokens, FieldToken{kind: TtText, value: string(buf[j:i])})
	}
	//fmt.Printf("parseContents(): %+v\n", tokens[len(tokens)-1])
	return tokens, err
}

func isEmptySourceError(err error) bool {
	return err.Error() == emptySourceError
}
