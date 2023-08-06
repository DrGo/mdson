package mdson

import (
	"fmt"
	"strings"
	"unicode"
)

func trimLower(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func trimLeftSpace(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}


// getBlockInfo returns header name and level eg "###Document" returns "document", 3}
// returned name is always lower-case (see parser_test.go for more details)
//assumes that is called for a string that starts with a space followed by #*
func getBlockInfo(line string) (string, int) {
	// hot path
	lgth := len(line)
	if lgth < 2 || line[0] != '#'{
		return "", -1 //not a block
	} 
	// find the first non-# char
	i := 1
	for ; i < lgth && line[i] == '#'; i++ { }
	//next char should be space 
	if i < lgth && line[i+1] != ' ' {
		return "", -1
	}	
	name := strings.TrimSpace(line[i:])
	if name == "" { //no name, heading but invalid
		return "", -1
	}
	return  trimLower(name),  i
}

func throw(value interface{}) (*Document, error) {
	switch unboxed := value.(type) {
	case string:
		return nil, fmt.Errorf("%s", unboxed)
	case error:
		return nil, fmt.Errorf("%s", unboxed)
	default:
		panic("unsupported argument type in throw()")
	}
}

func isArray(key string) bool {
	return strings.HasSuffix(key, " list")
}
