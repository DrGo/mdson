package mdson

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestTransformMD(t *testing.T) {
	SetDebug(DebugWarning)
	n, err := Parse(strings.NewReader(data))
	Equal(t, err, nil)
	if err != nil {
		return
	}
	w := bufio.NewWriter(os.Stdout)
	md := MDTransformer{}
	md.Transform(w, n)
	fmt.Println("********************************output************")
	w.Flush()
}
