package pixmatch

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
)

const (
	// YIQDeltaMax is the value of 35215. This is the maximum possible value
	// for the YIQ difference metric.
	// Read about YIQ NTSC https://en.wikipedia.org/wiki/YIQ
	YIQDeltaMax = 35215
)

// Image represents the images structure.
type Image struct {
	// Path to the image in filesystem.
	Path string
	// Format format as string like (png, jpg, gif).
	Format string
	// Image embedded image from Go's standard library.
	image.Image
}

// NewImage create a new image instance.
func NewImage() *Image {
	return &Image{}
}

// NewImageFromPath creates new images instance from the file system path.
func NewImageFromPath(path string) (*Image, error) {
	img := NewImage()
	img.SetPath(path)

	fp, err := os.Open(img.Path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	if err := img.Load(fp); err != nil {
		return nil, err
	}
	return img, nil
}

// SetPath sets the path to the image in filesystem.
func (img *Image) SetPath(path string) {
	img.Path = path
}

// Load reads data from the reader.
func (img *Image) Load(rd io.Reader) (err error) {
	img.Image, img.Format, err = image.Decode(rd)
	if err != nil {
		return
	}
	return
}

// TODO: move to correct oreder of code.
// DimensionsEqual checks that dimensions of the image is equal to dimension
// of other image.
func (img *Image) DimensionsEqual(img2 *Image) (bool, error) {
	if img.Bounds().Eq(img2.Bounds()) {
		return true, nil
	}
	return false, ErrDimensionsDoNotMatch
}

// Empty checks that images is empty of has theoretical size 0x0 pixels.
func (img *Image) Empty() bool {
	if img.Image == nil {
		return true
	}
	return img.Bounds().Empty()
}

// Compare returns the number of different pixels.
func (img *Image) Compare(img2 *Image, opts *Options) (int, error) {
	diff := 0
	if opts == nil {
		opts = NewOptions()
	}

	// If empty images return error.
	if img.Empty() || img2.Empty() {
		return ExitEmptyImage, ErrImageIsEmpty
	}

	// If dimensions do not match.
	if _, err := img.DimensionsEqual(img2); err != nil {
		return ExitDimensionsNotEqual, err
	}

	// If bytes are the same just return nothing to compare more.
	if img.Identical(img2) {
		// NOTE draw output ???
		return diff, nil
	}
	maxDelta := YIQDeltaMax * math.Pow(opts.Threshold, 2.0)

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			// TODO WORK
			// pos := img.Position(x, y)
			// fmt.Printf("(%d, %d),", x, y)
		}
	}

	fmt.Println("\nmaxDelta", maxDelta)

	return diff, nil
}

// Identical checks images that these are identical on bytes level.
// Looks like bytes.Equal() is the fastest way to compare 2 bytes arrays.
// Tried reflect.DeepEqual() and loop solutions. In the most cases
// bytes.Equal() is the best choice.
func (img *Image) Identical(img2 *Image) bool {
	return bytes.Equal(img.Bytes(), img2.Bytes())
}

// Bytes get the raw bytes of the pixel data.
func (img *Image) Bytes() []byte {
	switch img.ColorModel() {
	case color.RGBAModel:
		return img.Image.(*image.RGBA).Pix
	case color.RGBA64Model:
		return img.Image.(*image.RGBA64).Pix
	case color.NRGBAModel:
		return img.Image.(*image.NRGBA).Pix
	case color.NRGBA64Model:
		return img.Image.(*image.NRGBA64).Pix
	case color.AlphaModel:
		return img.Image.(*image.Alpha).Pix
	case color.Alpha16Model:
		return img.Image.(*image.Alpha16).Pix
	case color.GrayModel:
		return img.Image.(*image.Gray).Pix
	case color.Gray16Model:
		return img.Image.(*image.Gray16).Pix
		// TODO add JPEG support
		// case color.NYCbCrAModel:
		// 	return img.Image.(*image.NYCbCrA).Y
		// case color.YCbCrModel:
		// 	return img.Image.(*image.YCbCr).Y
	}

	switch img.ColorModel().(type) {
	case color.Palette:
		return img.Image.(*image.Paletted).Pix
	}

	// retrurn empty byte slice if something is really wrong with the image.
	return []byte{}
}

// Position is the positions of the pixel in bytes data.
// Eg. x1 is (r + 0), (g + 1), (b + 2), (a + 3), x2 = ...
func (img *Image) Position(x, y int) int {
	return (y*img.Bounds().Dx() + x) * 4
}
