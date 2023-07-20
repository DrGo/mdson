package main

import (
	"io"
	"os"
	"testing"
)

func TestCollection_ForEach(t *testing.T) {
	tests := []struct {
		pattern string
		ffn     FileFunc
		want    []io.Reader
		wantErr bool
	}{
		{"test-files/src/contents/datasets/*.mdson", cat, nil, false},
	}
	var err error
	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			c, _ := NewCollection(tt.pattern)
			tt.want, err = c.ForEach(tt.ffn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Collection.ForEach() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_, err = io.Copy(os.Stdout, tt.want[0])
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Collection.ForEach() = %v, want %v", got, tt.want)
			// }
		})
	}
}
