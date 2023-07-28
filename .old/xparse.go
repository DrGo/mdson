//go:build exclude
package mdson

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode"
)

const comment_chars="//"
type parseState int 
const (
	inHeader = iota

)

// Doc represents an entire document.
type Doc struct {
	state	   parseState	 
	Title      string
	Subtitle   string
	Summary    string
	Time       time.Time
	Authors    []Author
	Root	   *ttBlock
	Text 	   []string
	Sections   []Section
	Tags       []string
	OldURL     []string
}

func newDoc()*Doc{
	return &Doc{Root: newTokenBlock("Root")}
}

type Author struct{}
type Section struct{}

// Lines is a helper for parsing line-based input.
type Lines struct {
	rowNum    int // 0 indexed, so has 1-indexed number of last line returned
	rows    []string
}

func readLines(r io.Reader) ([]string, error) {
	var lines []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func newLinesFromFile(filename string) (*Lines, error) {
	f, err:= os.Open(filename)
	if err != nil {
		return nil, errors.Join(errors.New("error reading file "+filename), err)
	}
	return newLines(f)	
}


func newLines(r io.Reader) (*Lines, error) {
	lines, err:= readLines(r)
	if err != nil {
		return nil, err
	}
	return &Lines{0, lines}, nil
}

func (l *Lines) next() (text string, ok bool) {
	for {
		current := l.rowNum
		l.rowNum++
		if current >= len(l.rows) {
			return "", false
		}
		text = l.rows[current]
		// fmt.Println(l.rows[current])
		// found uncommented line, return it 
		if !strings.HasPrefix(text, comment_chars) {
			ok = true
			fmt.Println(current, ": ", text, )
			break
		}
	}
	return
}

func (l *Lines) back() {
	l.rowNum--
}

func (l *Lines) nextNonEmpty() (text string, ok bool) {
	for {
		text, ok = l.next()
		if !ok {
			return
		}
		if len(text) > 0 {
			break
		}
	}
	return
}

func (l *Lines) String() string {
var b strings.Builder
	for _,s := range l.rows {
		b.WriteString(s)
		b.WriteByte('\n')
	}
	return b.String()	
}

//Node implement parser's AST node
type Node interface {
	String() string
	Kind() string
}

//baseToken implements the basic token interface root of all of other tokens
type ttBase struct {
	kind string
	key  string
}
//singlton used to flag a comment 
var sComment= ttBase{kind:"Comment"}

const nodeDescLine = "type=%s, lineNum=%d, key=%s"

func (bt ttBase) String() string {
	return fmt.Sprintf(nodeDescLine, bt.Kind(), 0, bt.key)
}

func (bt ttBase) Kind() string {
	return bt.kind
}

// type ttListItem struct {
// 	ttBase
// }

func newListItem(item string) *ttBase {
	return &ttBase{kind: "ListItem", key: item}
}

func newText(item string) *ttBase {
	return &ttBase{kind: "Text", key: item}
}

type ttKV struct {
	ttBase 
	value string
}

func newKV(k,v string) *ttKV {
	return &ttKV{ttBase: ttBase{kind: "KV", key: k},value: v}
}

type ttBlock struct {
	ttBase
	// ordered list
	children []Node
	// stores block attributes that are not rendered
	attribs  map[string]string 
}

func newTokenBlock(name string) *ttBlock {
	return &ttBlock{ttBase: ttBase{kind: "Block", key: name}}
}

func (b *ttBlock) addChild(n Node) {
	b.children=append(b.children, n)
}

func trimLeftSpace(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}

func (d *Doc) parseRow(row string) Node {
	//scenario 1: empty line
	if len(row)==0 {
		return nil
	}
	//scenario 2: the first 2 chars are // -> commented line
	if strings.HasPrefix(row, "//") {
		return (&sComment)
	}
	//scenario 3: the first char is dot -> key-value pair
	if row[0] == '.' { //guaranteed to have >=1 char b/c of the empty check above
		colon := strings.Index(row, ":")
		// if no colon treat as text
		if colon == -1 {
			return newText(row)
		}
		// treat as key-value pair
		parts := strings.SplitN(row[1:], ":", 2) //split on the first colon skipping the dot
		return newKV(parts[0], parts[1])
		//TODO: handle escaped colon
		// return newListItem(s[1:])
	}
	// //scenario 4: block
	// if hd := getHeading(s); hd.level > 0 {
	// 	hp.Log(":", line, hd)
	// 	return newTokenBlock(hd.name).setLevel(hd.level)
	// }
	// //scenario 6: invalid key:value pair
	// parts := strings.SplitN(line, ":", 2) //split on the first colon
	// if len(parts) != 2 {
	// 	return newSyntaxError("likely missing ':'")
	// }
	// //scenario 7: an array of non-block types
	// key := trimLower(parts[0])
	// value := parts[1]
	// if isArray(key) && strings.TrimSpace(value) == "" {
	// 	return newList(key)
	// }
	//scenario 8: everything else is text 
	return newText(row)
}

func (d *Doc) parse(lines *Lines) error {
	// var ok bool
	// var row string
	// var current Node 
//ignore comments and any initial empty lines
		// row, ok = lines.nextNonEmpty() 
		// if !ok {
		// 	return nil //empty file or no data
		// }

	curBlock:=d.Root //start at the root 
	for _, row:= range lines.rows{
		n := d.parseRow(row)
		switch n.Kind(){
		case "Comment":
			continue
		case "Text":
			curBlock.addChild(n)
		case "KV":
			curBlock.attribs[n.(ttKV).key]=n.(ttKV).value 
		}
		//if section header, we finished parsing the header
		// if row[0] == '#'{ 
		// break 
		// }

		// a header row
	}
	// look for key-value pairs preceeding any section

	return nil
}
