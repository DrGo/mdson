package wfs

import (
	"fmt"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

const test_perm = 0664

func TestWFile(t *testing.T) {
	m := NewFS()
	m.WriteFile(  "hello", []byte("hello, world\n"), test_perm)
	m.WriteFile(	"fortune/k/ken.txt", []byte("If a program is too slow, it must have a loop.\n"),test_perm)

	if err := fstest.TestFS(m, "hello", "fortune/k/ken.txt"); err != nil {
		t.Fatal(err)
	}
}

func TestWFileChmodDot(t *testing.T) {
	m := NewFS()
	m.WriteFile("a/b.txt", []byte{}, 0666)
	m.WriteFile(".", []byte{}, 0777 | fs.ModeDir)

	buf := new(strings.Builder)
	fs.WalkDir(m, ".", func(path string, d fs.DirEntry, err error) error {
		fi, err := d.Info()
		if err != nil {
			return err
		}
		fmt.Fprintf(buf, "%s: %v\n", path, fi.Mode())
		return nil
	})
	want := `
.: drwxrwxrwx
a: d---------
a/b.txt: -rw-rw-rw-
`[1:]
	got := buf.String()
	if want != got {
		t.Errorf("FS modes want:\n%s\ngot:\n%s\n", want, got)
	}
}
