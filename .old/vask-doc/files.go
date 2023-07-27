package main

import (
	"io"
	"os"
	"path/filepath"
)

const (
	errorModeStop int = iota
	errorModeSkip
)

// FileFunc receives a reader based on a file contained in the collection
// and should return an active Reader and error status.
type FileFunc func(*Collection, io.Reader, int) (io.Reader, error)

// Collection holds info on a file collection
type Collection struct {
	errorMode int
	files     []string
	Values    map[string]interface{}
}

// NewCollection creates a new file collection
func NewCollection(pattern string) (*Collection, error) {
	c := Collection{Values: make(map[string]interface{})}
	var err error
	c.files, err = filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

//ForEach calls ffn for each file in the collection (as determined by the glob pattern
// passed to the NewCollection func) and returns an array of io.Readers (potentially with nil values) and error status
func (c *Collection) ForEach(ffn FileFunc) ([]io.Reader, error) {
	rws := make([]io.Reader, len(c.files))
	for i := 0; i < len(c.files); i++ {
		r, err := os.Open(c.files[i])
		if err != nil {
			if c.errorMode == errorModeStop {
				return nil, err
			}
			continue
		}
		rw, err := ffn(c, r, i)
		r.Close()
		if err != nil {
			if c.errorMode == errorModeStop {
				return nil, err
			}
			continue
		}
		rws[i] = rw
	}
	return rws, nil
}

// GetFileName returns filename at index
func (c *Collection) GetFileName(index int) string {
	return c.files[index]
}
