package pixmatch

import "errors"

var (
	//ErrDimensionsDoNotMatch represents an error when the dimensions of two
	//images do not match.
	ErrDimensionsDoNotMatch = errors.New("Images dimensions do not match")

	// ErrImageIsEmpty occurs when one of the images, or both of them,
	// are empty.
	ErrImageIsEmpty = errors.New("One or both images are empty")

	// ErrCorruptedImage occurs when the data for the image is corrupted and
	// cannot be read or decoded.
	ErrCorruptedImage = errors.New("Image data is corrupted")

	// ErrCannotWriteOutputDiff occurs when output cannot be written.
	ErrCannotWriteOutputDiff = errors.New("Cannot write diff output")

	// ErrUnknownFormat occurs when the image format is not supported or
	// unknown.
	ErrUnknownFormat = errors.New("Unknown image format")

	// ErrInvalidColorFormat occurs when user enter invalid color format.
	ErrInvalidColorFormat = errors.New("Invalid color format")

	// ErrMissingImage occurs when one or both images are missing.
	ErrMissingImage = errors.New("One or both images are missing")
)

// Exit codes that are not defined in the [BSD and Linux specifications].
//
// [BSD and Linux specifications]: https://freedesktop.org/software/systemd/man/systemd.exec.html#Process%20Exit%20Codes
const (
	// ExitOk when the program exited successfully.
	ExitOk = 0

	// ExitFSFail occurs when there is a problem with the file system.
	ExitFSFail = 100

	// ExitEmptyImage occurs when the image (or both) are empty.
	ExitEmptyImage = 101

	// ExitDimensionsNotEqual occurs when the images dimensions are not equal.
	ExitDimensionsNotEqual = 102

	// ExitInvalidInput input parameters and/or flags are invalid.
	// Check usage.
	ExitInvalidInput = 103

	// ErrMissingImage one or both images are missing.
	ExitMissingImage = 104

	// ExitUnknownFormat if format of the image is not supported.
	ExitUnknownFormat = 105

	// ExitUnknown all other failings.
	ExitUnknown = 199
)
