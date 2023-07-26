package mdson

import (
	"flag"
	"testing"
)


var TestDebug=flag.Int("debug", 0, "")


func TestEval(t *testing.T) {
	doc, err := ctx.ParseFile("", specs)
	Equal(t, err, nil)
	if err != nil {
		return
	}
	Equal(t, doc.attribs["date"], "12July2023")
	Equal(t, doc.attribs["today"], "Today is 12July2023")
	// t.Logf("\n\n%s\n", newPrinter().print(doc.root)) 
}