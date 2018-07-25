package mds

import (
	"fmt"
	"strings"
	"testing"
)

func Test_parseHDScript(t *testing.T) {
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
		{"wrong block type", wrongBlockType, "general", "", "", true},
	}
	hp := NewParser(DefaultParserOptions().SetDebug(DebugAll))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTokens, err := hp.Parse(strings.NewReader(tt.src))
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHDScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(gotTokens)
			if gotv := gotTokens[tt.block][tt.wantk]; gotv != tt.wantv {
				t.Errorf("parseHDScript() = %v, want %v", gotv, tt.wantv)
			}
		})
	}
}

const full = `+++
Settings: general
Version: 1.0.0
Command: run
OverWriteOutputFile: true
OutputFileName: fromconfig.docx
Debug: 2
WorkDirName: 
//comments: empty lines ignored
+++
Settings: rosewood
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
+++

Document: document1
TemplateFileName: 
InputDir: 
+++
Section: section1
InputDir: /Users/salah/Dropbox/code/go/src/github.com/drgo/htmldocx/cmd
AddPageBreakAfterEachInputFile: true
Size: Width: 12240, Height: 15840, Orientation: portrait
Margins: Top: 10, Right: 1440, Bottom: 1440, Left: 1440,  Header: 360, Footer: 360,  Gutter: 0
Headers: default: header1
Footers: default: footer1
Contents: 
${include: tab1old.html}${include: test.html}
+++   

Header: header1
InputDir: 
Contents:
${htmldocx: timestamp} text in next line is created using raw xml ${xml:<w:p><w:rPr><w:tab 
w:val="right" w:leader="dot" w:pos="2160"/></w:rPr><w:r><w:rPr><w:b/><w:i/><w:color 
w:val="8300ff"/><w:rFonts w:ascii="Courier New" w:hAnsi="Times New Roman" w:cs="Times New 
Roman"/></w:rPr><w:tab/><w:tab/><w:t>tabbed-bold-italic-magenta-Courier New</w:t></w:r></w:p>}
+++

Footer: footer1
InputDir: "empty"
Contents:
${word: PAGE}
some other text
+++`

const wrongBlockType = `+++
xSettings: general
Version: 1.0.0`
