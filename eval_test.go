package mdson

import (
	"testing"

	"github.com/drgo/booker/tu"
)




func TestEval(t *testing.T) {
	doc, err := ctx.ParseFile("", specs)
	tu.Equal(t, err, nil)
	if err != nil {
		return
	}
	tu.Equal(t, doc.attribs["date"], "12July2023")
	tu.Equal(t, doc.attribs["today"], "Today is 12July2023")
		
	// t.Logf("\n\n%s\n", doc.root.Children()[0]) 
}
