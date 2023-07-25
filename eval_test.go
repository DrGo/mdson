package mdson

import (
	"strings"
	"testing"
)


func TestEval(t *testing.T) {
	SetDebug(DebugWarning)
	n, err := Parse(strings.NewReader(specs))
	Equal(t, err, nil)
	if err != nil {
		return
	}
	n, err = Eval(n, DefaultParserOptions().SetDebug(DebugAll))
	Equal(t, err, nil)
	if err != nil {
		return
	}
	
	t.Logf("\n\n%s\n", newPrinter().print(n)) 
}
