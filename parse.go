package mdson

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const comment_chars="//"

// Doc represents an entire document.
type Doc struct {
	Title      string
	Subtitle   string
	Summary    string
	Time       time.Time
	Authors    []Author
	TitleNotes []string
	Sections   []Section
	Tags       []string
	OldURL     []string
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
		fmt.Println(l.rows[current])
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
