package pixmatch

import "errors"

var (
	// ErrDimensionsDoNotMatch represents an error when images dimentions
	// aren't match.
	ErrDimensionsDoNotMatch = errors.New("Images dimensions do not match")

	// ErrImageIsEmpty means that one of the image, or both or them are
	// empty.
	ErrImageIsEmpty = errors.New("One or both images are empty")

	// ErrCorruptedImage means that data fot the image is corrupted and
	// cannot be decoded.
	ErrCorruptedImage = errors.New("Image data is corrupted")

	// ErrCannotWriteOutputDiff means that output cannot be written.
	ErrCannotWriteOutputDiff = errors.New("Cannot write diff output")
)

// Exit codes not defined in BSD and Linux
// https://freedesktop.org/software/systemd/man/systemd.exec.html#id-1.20.8
const (
	// ExitOk then programm exited successfully.
	ExitOk = 0

	// ExitFSFail exit status when something is wrong with file system.
	ExitFSFail = 100

	// ExitEmptyImage exit status when the images is empty.
	ExitEmptyImage = 101

	// ExitDimensionsNotEqual status when the images dimentions are not equal.
	ExitDimensionsNotEqual = 102
)
