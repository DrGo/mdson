package mdson

import "runtime"

var debug int

// SetDebug sets the default debug level for all package routines
func SetDebug(level int) {
	debug = level
}

const (
	//DebugSilent print errors only
	DebugSilent int = iota
	//DebugWarning print warnings and errors
	DebugWarning
	//DebugUpdates print execution updates, warnings and errors
	DebugUpdates
	//DebugAll print internal debug messages, execution updates, warnings and errors
	DebugAll
)

//ParserOptions holds parsing options
type ParserOptions struct {
	Debug int
}

//DefaultParserOptions returns reasonable default for parsing
func DefaultParserOptions() *ParserOptions {
	return &ParserOptions{
		Debug: DebugUpdates,
	}
}

//SetDebug sets verbosity level
func (po *ParserOptions) SetDebug(d int) *ParserOptions {
	po.Debug = d
	return po
}

var lineBreak = "\n"

func init() {
	if runtime.GOOS == "windows" {
		lineBreak = "\r\n"
	}
}
