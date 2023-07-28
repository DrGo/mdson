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

//Options holds parsing options
type Options struct {
	Debug ui.Debug
	BufferCap int
	// style of md list generated when list style is not specified
	// ol=ordered, ul=unordered 
	DefaultListStyle string 
}

//DefaultOptions returns reasonable default for parsing
func DefaultOptions() *Options {
	return &Options{
		Debug: ui.DebugUpdates,
		BufferCap: 1024 * 10,
		DefaultListStyle: "ol",
	}
}


func (po Options ) String() string {
	return fmt.Sprintf(
	"Settings: Debug: %s | Buffer Capacity %d\n", po.Debug, po.BufferCap)
}
//SetDebug sets verbosity level
func (po *Options) SetDebug(d ui.Debug) *Options {
	po.Debug = d
	return po
}

var lineBreak = "\n"

func init() {
	if runtime.GOOS == "windows" {
		lineBreak = "\r\n"
	}
}
