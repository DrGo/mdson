package mdson

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
