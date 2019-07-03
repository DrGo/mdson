package mdson

import (
	"fmt"
	"strings"
)

type heading struct {
	name  string
	level int
}

// getHeading returns header name and level eg "###Document" returns heading{name: "document", level: 3}
// returned name is always lower-case (see parser_test.go for more details)
func getHeading(line string) heading {
	i := 0
	for ; i < len(line) && line[i] == '#'; i++ { //no utf8 needed, we are only looking for a byte #
	}
	if i == 0 { //# not found, not a heading
		return heading{name: "", level: -1}
	}
	name := strings.TrimSpace(line[i:])
	if name == "" { //no name, heading but invalid
		return heading{name: "", level: -1}
	}
	return heading{name: trimLower(name), level: i}
}

func throw(value interface{}) (*ttBlock, error) {
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
