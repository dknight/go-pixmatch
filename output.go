package pixmatch

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

// Output represents the structure of output and its parameters.
type Output struct {
	// Dest is the writer where output will be written.
	Dest io.Writer
	// Image holds the image information.
	Image *image.RGBA
}

// NewOutput creates output image for given path in filesystem.
func NewOutput(path string, w, h int) (*Output, error) {
	dest, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	return &Output{
		Dest:  dest,
		Image: img,
	}, nil
}

// Save writes data to output.
func (out *Output) Save(format string) (err error) {
	switch format {
	case "gif":
		err = gif.Encode(out.Dest, out.Image, nil)
	case "jpeg", "jpg":
		err = jpeg.Encode(out.Dest, out.Image, nil)
	case "png":
		err = png.Encode(out.Dest, out.Image)
	default:
		err = errors.New("Unknown image format")
	}
	return err
}
