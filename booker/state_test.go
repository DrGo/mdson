package booker

import (
	"os"
	"testing"

	"github.com/drgo/core/ui"
	"github.com/drgo/mdson"
	"github.com/drgo/mdson/tu"
)

const (
	TestDir= "../test/gen"
	FileCount= 100
	TestDebug = ui.DebugSilent
)

var cfg = newConfig()

func Test(t *testing.T) {
	
	s := newState(os.DirFS(TestDir), cfg)	
	files, err := s.Glob("")

	tu.Equal(t, err, nil)
	tu.Equal(t, len(files), FileCount)
	ctx:= mdson.NewContext(mdson.DefaultOptions().SetDebug(ui.Debug(TestDebug)))
	for _, fn := range files[:5] {
		s.Log("**************************Filename: ", fn)
		r, err := s.sfs.Open(fn)
		tu.Equal(t, err, nil)
		doc, err:= ctx.ParseFile(fn, r)
		tu.Equal(t, err, nil)
		if err== nil {
			continue
		}
		tu.Equal(t, doc.Attribs()["name"], fn)
		t.Log(doc)
			
	}	
}
