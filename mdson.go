package mdson

import (
	"fmt"
	"runtime"

	"github.com/drgo/core/ui"
)

// var debug  ui.Debug
// // SetDebug sets the default debug level for all package routines
// func SetDebug(level ui.Debug) {
// 	debug = level
// }
//

//ParserOptions holds parsing options
type ParserOptions struct {
	Debug ui.Debug
	BufferCap int 
}

//DefaultParserOptions returns reasonable default for parsing
func DefaultParserOptions() *ParserOptions {
	return &ParserOptions{
		Debug: ui.DebugUpdates,
		BufferCap: 1024 * 10,
	}
}


func (po ParserOptions ) String() string {
	return fmt.Sprintf(
	"Settings: Debug: %s | Buffer Capacity %d\n", po.Debug, po.BufferCap)
}
//SetDebug sets verbosity level
func (po *ParserOptions) SetDebug(d ui.Debug) *ParserOptions {
	po.Debug = d
	return po
}

var lineBreak = "\n"

func init() {
	if runtime.GOOS == "windows" {
		lineBreak = "\r\n"
	}
}
