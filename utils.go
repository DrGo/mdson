package mdson

import "strings"

func trimLower(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
