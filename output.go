package pixmatch

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

// Output represents an image instance and its destination, where it will be
// written.
type Output struct {
	// Dest is the destination, where the output will be written.
	Dest io.Writer

	// Image holds the image instance.
	Image *image.RGBA
}

// NewOutput creates an output image for a given file system path.
func NewOutput(dest io.Writer, w, h int) (*Output, error) {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	return &Output{
		Dest:  dest,
		Image: img,
	}, nil
}

// Save encodes and writes image data to the output destination.
func (out *Output) Save(format string) (err error) {
	switch format {
	case "gif":
		err = gif.Encode(out.Dest, out.Image, nil)
	case "jpeg":
		err = jpeg.Encode(out.Dest, out.Image, nil)
	case "png":
		err = png.Encode(out.Dest, out.Image)
	default:
		err = ErrUnknownFormat
	}
	return
}
