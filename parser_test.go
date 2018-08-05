package mdson

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_Parser_parse(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		block   string
		wantk   string
		wantv   string
		wantErr bool
	}{
		//TODO: revize this test or delete
		{"correct", full, "general", "version", "1.0.0", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hp, _ := NewParser(strings.NewReader(tt.src), DefaultParserOptions().SetDebug(DebugAll))
			gotTokens, err := hp.parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHDScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				node := gotTokens.children[0].(*ttKVPair)
				fmt.Printf("******************\n%s\n", node)
			}
			// if gotv := gotTokens[tt.block][tt.wantk]; gotv != tt.wantv {
			// 	t.Errorf("parseHDScript() = %v, want %v", gotv, tt.wantv)
			// }
		})
	}
}

const wrongBlockType = `+++
xSettings: general
Version: 1.0.0`

func Test_getHeading(t *testing.T) {
	tests := []struct {
		name string
		line string
		want heading
	}{
		{"0H", "Document", heading{name: "", level: -1}},
		{"1 H", " #Document", heading{name: "", level: -1}},
		{"1H empty", "# ", heading{name: "", level: -1}},
		{"1H ", "#  Document", heading{name: "document", level: 1}},
		{"1H", "#Document", heading{name: "document", level: 1}},
		{"3H", "###Document", heading{name: "document", level: 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getHeading(tt.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getHeading() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	const fileName = "/Users/salah/Dropbox/code/go/src/github.com/drgo/mdson/carpenter.mdson"
	tests := []struct {
		name     string
		fileName string
		wantRoot *ttBlock
		wantErr  bool
	}{
		{"", fileName, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRoot, err := ParseFile(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRoot, tt.wantRoot) {
				t.Errorf("ParseFile() = %v, want %v", gotRoot, tt.wantRoot)
			}
		})
	}
}
