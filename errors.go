package pixmatch

import "errors"

// ErrDimensionsDoNotMatch represents an error when images dimentions
// aren't match.
var ErrDimensionsDoNotMatch = errors.New("Images dimensions do not match.")

const (
	// ExitOk then programm exited successfully.
	ExitOk = 0
	// ExitFSFail exit status when something is wrong with file system.
	ExitFSFail = 1
)
