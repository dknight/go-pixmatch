package pixmatch

import "errors"

var (
	//ErrDimensionsDoNotMatch represents an error when the dimensions of two
	//images do not match.
	ErrDimensionsDoNotMatch = errors.New("Images dimensions do not match")

	// ErrImageIsEmpty means that one of the images, or both of them,
	// are empty.
	ErrImageIsEmpty = errors.New("One or both images are empty")

	// ErrCorruptedImage means that the data for the image is corrupted and
	// cannot be read or decoded.
	ErrCorruptedImage = errors.New("Image data is corrupted")

	// ErrCannotWriteOutputDiff means that output cannot be written.
	ErrCannotWriteOutputDiff = errors.New("Cannot write diff output")

	// ErrUnknownFormat means that the image format is not supported or
	// unknown.
	ErrUnknownFormat = errors.New("Unknown image format")
)

// Exit codes that are not defined in the BSD and Linux specifications
// https://freedesktop.org/software/systemd/man/systemd.exec.html#id-1.20.8
const (
	// ExitOk when the program exited successfully.
	ExitOk = 0

	// ExitFSFail occurs when there is a problem with the file system.
	ExitFSFail = 100

	// ExitEmptyImage occurs when the image (or both) are empty.
	ExitEmptyImage = 101

	// ExitDimensionsNotEqual occurs when the images dimensions are not equal.
	ExitDimensionsNotEqual = 102
)
