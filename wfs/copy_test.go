package wfs

import (
	"os"
	"testing"
	"testing/fstest"

	"github.com/drgo/mdson/tu"
)

const copyTestDir= "test/copy"

func TestFlushFS(t *testing.T) {
	sfs := newTestFS(t)
	err:= FlushFS(copyTestDir, sfs)
	tu.Equal(t, err, nil)
	
	if err := fstest.TestFS(os.DirFS(copyTestDir), "hello", "fortune/k/ken.txt"); err != nil {
		t.Fatal(err)
	}
}
