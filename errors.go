package pixmatch

import "errors"

var (
	// ErrDimensionsDoNotMatch represents an error when images dimentions
	// aren't match.
	ErrDimensionsDoNotMatch = errors.New("Images dimensions do not match")

	// ErrImageIsEmpty means that one of the image, or both or them are
	// empty.
	ErrImageIsEmpty = errors.New("One or both images are empty")
)

const (
	// ExitOk then programm exited successfully.
	ExitOk = 0

	// ExitFSFail exit status when something is wrong with file system.
	ExitFSFail = 1

	// ExitEmptyImage exit status when the images is empty.
	ExitEmptyImage = 2
)
