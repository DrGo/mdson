package mdson

import (
	"bytes"
	"fmt"
	"strings"
)

//TODO: test struct tags "mdson"
func ExampleOldPa() {
	SetDebug(DebugSilent)
	root, err := Parse(strings.NewReader(full))
	if err != nil {
		fmt.Printf("error parsing %v \n", err)
		return
	}
	rootBlk, ok := root.(*ttBlock)
	// fmt.Println("rootBlk", rootBlk)
	if !ok {
		fmt.Println("parser returned unexpected type: root is not *ttBlock")
	}
	doc, ok := rootBlk.getChildByName("document").(*ttBlock)
	if !ok {
		fmt.Println("parser returned unexpected type: doc is not *ttBlock")
	}
	sections := doc.getChildByName("sections")
	if sections == nil {
		fmt.Println("nil sections")
		return
	}
	fmt.Println(sections.Name(), sections.Kind())
	blk := sections.(*ttBlock)
	fmt.Println(len(blk.children))
	sec1 := blk.children[0]
	fmt.Println(sec1.Name())
	blk = sec1.(*ttBlock)
	fmt.Println(blk.getChildByName("contents").Name())
	// kvp := blk.children[1].(*ttKVPair)
	// fmt.Println(kvp.key, kvp.value)
	list := blk.getChildByName("inputfiles").(*ttList)
	fmt.Println(list.Name())
	fmt.Println(list.items[0].key)
	ft := blk.getChildByName("freetext").(*ttLiteralString)
	fmt.Println(ft.value)
	// Output:
	// sections Block
	// 2
	// section1
	// contents
	// inputfiles
	// file1
	//  anything here goes <xml of some sort />some other test
}

func ExampleUnmarshal() {
	SetDebug(DebugSilent)
	var job Job
	if err := Unmarshal(strings.NewReader(correct), &job); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(job.Version)
	fmt.Println(len(job.Document.Sections))
	fmt.Println(job.Document.Sections[0].InputDir)
	fmt.Println(job.Document.Sections[0].Contents)
	fmt.Println(job.Document.Sections[0].Props.Size.Width)
	fmt.Println(job.Document.Sections[0].Props.HeadersFooters[0].ID)
	// Output:
	//
	// 1
	// /users/salah/dropbox/code/go/src/github.com/drgo/htmldocx/cmd
	// ${include: tab1old.html}${include: test.html}
	// 12240
	// header1
}

func ExampleMarshal() {
	SetDebug(DebugAll)
	var job Job
	if err := Unmarshal(strings.NewReader(correct), &job); err != nil {
		fmt.Println(err)
		return
	}
	//test byte slice
	job.ByteSlice = []byte("this is a byte slice")
	var buf []byte
	var err error
	if buf, err = Marshal(&job); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(buf))
	// Output:
	//
	// 1
	// /users/salah/dropbox/code/go/src/github.com/drgo/htmldocx/cmd
	// ${include: tab1old.html}${include: test.html}
	// 12240
	// header1
}

func ExampleMarshalSubBlock() {
	SetDebug(DebugAll)
	doc := &Document{}
	var buf bytes.Buffer
	enc := NewEncoder(&buf) //creating own encoder to control depth level
	enc.SetBlockLevel(2)    //to get ## Document
	if err := enc.Encode(doc); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(buf.String())
	// Output:
	//
	// 1
	// /users/salah/dropbox/code/go/src/github.com/drgo/htmldocx/cmd
	// ${include: tab1old.html}${include: test.html}
	// 12240
	// header1
}

type Job struct {
	Version             string
	Command             string
	OverWriteOutputFile bool
	OutputFileName      string
	Debug               int
	WorkDirName         string
	ByteSlice           []byte
	InputFileNames      []string
	Rosewood            *Settings
	Document            *Document
}

func (f *Job) String() string {
	buf, err := Marshal(f)
	if err == nil {
		return string(buf)
	}
	return "marshalling error " + err.Error()
}

type Document struct {
	ID               ID
	TemplateFileName string `mdson:",omitempty"`
	InputDir         string
	// HeadersFooters   []*HeaderFooter
	Sections []*Section
}

type Settings struct {
	ConvertOldVersions   bool
	ConvertFromVersion   string
	DoNotInlineCSS       bool
	MandatoryCol         bool `mdson:"-"`
	MaxConcurrentWorkers int
}

type Section struct {
	ID         ID
	Contents   string
	InputDir   string
	InputFiles []string
	FreeText   string
	Props      *SectionProps
}

type SectionProps struct {
	Size           *PageSize
	Margins        *PageMargins
	HeadersFooters []*SectionHeaderFooter
}

type PageSize struct {
	Width       string
	Height      string
	Orientation string
}

//PageMargins holds info on a section page margins
type PageMargins struct {
	Top    string
	Right  string
	Bottom string
	Left   string
	Header string
	Footer string
	Gutter string
}

//SectionHeaderFooter holds info on headers and footers attached to this section
type SectionHeaderFooter struct {
	ID     string //ID of the header/footer in the Document headerfooters list
	HFType string `json:"type"` //default, odd etc
}

const full = `#file1
Version: 1.0.0
Command: run
OverWriteOutputFile: true
OutputFileName: fromconfig.docx
Debug: 2
WorkDirName: 
//comments: empty lines ignored

## Document
TemplateFileName: 
InputDir: 

### Sections List
#### Section1
contents: "hello"
InputFiles List:
- file1
- file2
- file3
freeText:
<< anything here goes 
<xml of some sort />
some other test>>

#### Section2
contents: "goodbye"
InputFiles List:
freeText:
<<< anything here 'goes'
<head> other \n stuff </head>
some "other" test>>>
<<${include: tab1old.html}${include: test.html}>>
`

const correct = `# Job
Command: run
OverWriteOutputFile: true
Output File Name: fromconfig.docx
Debug: 2
WorkDirName:
InputFileNames List:
- tab1old.html
- test.html

## Rosewood
ConvertOldVersions: false
ConvertFromVersion: v01
DoNotInlineCSS: false
MandatoryCol: false
MaxConcurrentWorkers: 30
PreserveWorkFiles: false
ReportAllError: true
SaveConvertedFile: false
StyleSheetName: 
TrimCellContents: false

## Document
TemplateFileName: 
InputDir: 

### Sections List
#### Section1	  
InputDir: /Users/salah/Dropbox/code/go/src/github.com/drgo/htmldocx/cmd
AddPageBreakAfterEachInputFile: true
Contents: 
<<${include: tab1old.html}${include: test.html}>>

##### Props
###### Size 
Width       : 12240
Height      : 15840
Orientation : portrait

###### HeadersFooters List	
####### HeaderFooter1 
ID     : header1
HFType   : default

### HeaderFooters List 
#### Header1
InputDir: 
Contents:
<<${htmldocx: timestamp} text in next line is created using raw xml ${xml:<w:p><w:rPr><w:tab 
w:val="right" w:leader="dot" w:pos="2160"/></w:rPr><w:r><w:rPr><w:b/><w:i/><w:color 
w:val="8300ff"/><w:rFonts w:ascii="Courier New" w:hAnsi="Times New Roman" w:cs="Times New 
Roman"/></w:rPr><w:tab/><w:tab/><w:t>tabbed-bold-italic-magenta-Courier New</w:t></w:r></w:p>}>>
#### Footer1
InputDir:
Contents:
<<${word: PAGE}
some other text>>
`
