package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func inlineCSS(srcFileName, cssFileName string) ([]byte, error) {
	srcBuf, err := ioutil.ReadFile(srcFileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %s", err)
	}
	headPos := bytes.Index(srcBuf, []byte("</head>"))
	if headPos == -1 {
		return nil, fmt.Errorf("file [%s] missing </head>", srcFileName)
	}
	css, err := os.Open(cssFileName)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %s", err)
	}
	var out bytes.Buffer
	out.Write(srcBuf[:headPos])
	out.WriteString("<style>\n")
	if _, err = io.Copy(&out, css); err != nil {
		return nil, fmt.Errorf("cannot read file: %s", err)
	}
	out.WriteString("</style>\n")
	out.Write(srcBuf[headPos:])
	return out.Bytes(), nil
}
