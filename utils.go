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
//assumes that is called for a string that starts with #
func getBlockInfo(line string) (string, int) {
	// hot path
	if line[0] != '#'{
		return "", -1 //not a block
	} 
	// find the first non-# char
	i := 1
	for ; i < len(line) && line[i] == '#'; i++ { //FIXME: for utf8 use
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
