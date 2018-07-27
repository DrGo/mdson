package mds

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
		{"correct", full, "general", "version", "1.0.0", false},

		//the following should produce an error, wantErr = true
		// {"wrong block type", wrongBlockType, "general", "", "", true},
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
				file, _ := gotTokens.children[0].(*ttBlock)
				for _, b := range file.children {
					fmt.Printf("******************\n%s\n", b)
				}
			}
			// if gotv := gotTokens[tt.block][tt.wantk]; gotv != tt.wantv {
			// 	t.Errorf("parseHDScript() = %v, want %v", gotv, tt.wantv)
			// }
		})
	}
}

const full = `# file: file1
Version: 1.0.0
Command: run
OverWriteOutputFile: true
OutputFileName: fromconfig.docx
Debug: 2
WorkDirName: 
//comments: empty lines ignored

## Settings: rosewood
ConvertOldVersions: false
ConvertFromVersion: v01
DoNotInlineCSS: false
MandatoryCol: false
MaxConcurrentWorkers: 30
PreserveWorkFiles: false
ReportAllError: false
SaveConvertedFile: false
StyleSheetName: 
TrimCellContents: false

## Document: document1
TemplateFileName: 
InputDir: 

### Section: section1
ID: section1
contents: "hello"
InputFiles:

### Section: section2
ID: section2
contents: "goodbye"
InputFiles:
- file1
- file2
- file3

## Documen: document2
Title: new document
`

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
	const fileName = "/Users/salah/Dropbox/code/go/src/github.com/drgo/mds/carpenter.mdon"
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

// - file1
// - file2
// - file3
