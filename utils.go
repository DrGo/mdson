package mdson

import (
	"strings"
	"unicode"
)

func trimLower(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func trimLeftSpace(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}
