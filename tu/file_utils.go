package tu

import (
	"bytes"
	"io"
	"os"
)


type TFile struct {
}

var File TFile

func (f TFile) MustRead(filename string) io.Reader {	
	content, err := os.ReadFile(filename)
	if err != nil {
		panic("test setup failed" + err.Error())
	}
	return bytes.NewReader(content) 
}


