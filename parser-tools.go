package mdson

import (
	"fmt"
	"strings"
)

type heading struct {
	name  string
	level int
}

func getHeading(line string) heading {
	i := 0
	for ; i < len(line) && line[i] == '#'; i++ { //no utf8 needed, we are only looking for a byte #
	}
	if i == 0 { //# not found, not a heading
		return heading{name: "", level: 0}
	}
	name := trimLower(line[i:])
	if name == "" { //no name, heading but invalid
		return heading{name: "", level: -1}
	}
	return heading{name: name, level: i}
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
