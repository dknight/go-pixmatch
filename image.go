package pixmatch

import (
	"bytes"
	"image"

	// _ "image/gif"
	// _ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"os"
	"reflect"
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

// TODO: move to correct order of code.

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
		opts = DefaultOptions()
	}
	diffColor := opts.ResolveDiffColor()

	// If empty images return error.
	if img.Empty() || img2.Empty() {
		return ExitEmptyImage, ErrImageIsEmpty
	}

	// If dimensions do not match, return error.
	if _, err := img.DimensionsEqual(img2); err != nil {
		return ExitDimensionsNotEqual, err
	}

	// If bytes are the same just return nothing to compare more.
	if img.Identical(img2) {
		// TODO draw output gray
		return diff, nil
	}
	maxDelta := YIQDeltaMax * math.Pow(opts.Threshold, 2.0)

	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			point := image.Point{x, y}
			pos := img.Position(point)
			delta := img.ColorDelta(img2, pos, pos, false)

			if math.Abs(delta) > maxDelta {
				if opts.DetectAA &&
					(img.Antialiased(img2, point) ||
						img2.Antialiased(img, point)) {
					if opts.Output != nil {
						opts.Output.Image.Set(x, y, opts.AAColor)
					}
				} else {
					if opts.Output != nil {
						opts.Output.Image.Set(x, y, diffColor)
					}
					diff++
				}
			} else if opts.Output != nil {
				// TODO draw output gray
			}
		}
	}
	if opts.Output != nil {
		err := opts.Output.Save()
		if err != nil {
			return -1, err
		}
	}
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
// Reflection is never clear (https://go-proverbs.github.io/)
// TODO add JPEG support
func (img *Image) Bytes() (bs []byte) {
	if img == nil {
		return bs
	}
	val := reflect.ValueOf(img.Image)
	ptr := reflect.Indirect(val)
	pixs := ptr.FieldByName("Pix")
	// retrurn empty byte slice if something is really wrong with the image.
	if pixs.IsValid() {
		bs = pixs.Bytes()
	}
	return bs

	// -----
	// If reflection is not enough clear use something like this or generics!
	//------

	// switch img.ColorModel() {
	// case color.RGBAModel:
	// 	return img.Image.(*image.RGBA).Pix
	// case color.RGBA64Model:
	// 	return img.Image.(*image.RGBA64).Pix
	// case color.NRGBAModel:
	// 	return img.Image.(*image.NRGBA).Pix
	// case color.NRGBA64Model:
	// 	return img.Image.(*image.NRGBA64).Pix
	// case color.AlphaModel:
	// 	return img.Image.(*image.Alpha).Pix
	// case color.Alpha16Model:
	// 	return img.Image.(*image.Alpha16).Pix
	// case color.GrayModel:
	// 	return img.Image.(*image.Gray).Pix
	// case color.Gray16Model:
	// -------------- JPEG ------------------
	// 	return img.Image.(*image.Gray16).Pix
	// 	// case color.NYCbCrAModel:
	// 	// 	return img.Image.(*image.NYCbCrA).Y
	// 	// case color.YCbCrModel:
	// 	// 	return img.Image.(*image.YCbCr).Y
	// }

	// switch img.ColorModel().(type) {
	// case color.Palette:
	// 	return img.Image.(*image.Paletted).Pix
	// }
}

// Uint32 converts bytes array to []uint32 slice.
func (img *Image) Uint32() []uint32 {
	bs := img.Bytes()
	ui32 := make([]uint32, len(bs))
	for i, b := range bs {
		ui32[i] = uint32(b)
	}
	return ui32
}

// Position is the positions of the pixel in bytes data.
// Eg. x1 is (r + 0), (g + 1), (b + 2), (a + 3), x2 = ...
func (img *Image) Position(p image.Point) int {
	return (p.Y*img.Bounds().Dx() + p.X) * 4
}

// ColorDelta is the squared YUV distance between colors at this pixel
// position, returns negative if the img2 pixel is darker.
// If argument onlyY is true, the only brightness level will be returned
// (Y component of YIQ model).
func (img *Image) ColorDelta(img2 *Image, m, n int, onlyY bool) float64 {
	bs1, bs2 := img.Uint32(), img2.Uint32()
	color1 := NewColor(bs1[m+0], bs1[m+1], bs1[m+2], bs1[m+3])
	color2 := NewColor(bs2[m+0], bs2[m+1], bs2[m+2], bs2[m+3])

	// If all colors are the same then no delta.
	if color1.Equals(color2) {
		return 0
	}

	if color1.A < 255 {
		color1 = color1.Blend(color1.A)
	}

	if color2.A < 255 {
		color2 = color2.Blend(color2.A)
	}

	y1, y2 := color1.Y(), color2.Y()
	y := y1 - y2
	if onlyY {
		return y
	}

	i := color1.I() - color2.I()
	q := color1.Q() - color2.Q()
	delta := 0.5053*y*y + 0.299*i*i + 0.1957*q*q // math.Pow(x, 2.0)

	if y1 > y2 {
		delta *= -1.0
	}

	return delta
}

// Antialiased checks that point is anti-aliased.
// TODO use vector points? same as SameColorNeighbors
// TODO not correctly?
func (img *Image) Antialiased(img2 *Image, pt image.Point) bool {
	neibrs := 0
	x1 := intMax(pt.X-1, 0)
	y1 := intMax(pt.Y-1, 0)
	x2 := intMin(pt.X+1, img.Bounds().Dx()-1)
	y2 := intMin(pt.Y+1, img.Bounds().Dy()-1)
	pos := img.Position(pt)

	if pt.X == x1 || pt.X == x2 || pt.Y == y1 || pt.Y == y2 {
		neibrs++
	}

	min, max := 0.0, 0.0
	minX, minY, maxX, maxY := 0, 0, 0, 0

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			if x == pt.X && y == pt.Y {
				continue
			}

			pos2 := img.Position(image.Point{x, y})
			delta := img.ColorDelta(img, pos, pos2, true)

			if delta == 0 {
				neibrs++
				if neibrs > 2 {
					return false
				}
			} else if delta < min {
				min = delta
				minX = x
				minY = y
			} else if delta > max {
				max = delta
				maxX = x
				maxY = y
			}
		}
	}

	if min == 0 || max == 0 {
		return false
	}

	return (img.HasLeast3Neighbors(image.Point{minX, minY}) &&
		img2.HasLeast3Neighbors(image.Point{minX, minY})) ||
		(img.HasLeast3Neighbors(image.Point{maxX, maxY}) &&
			img2.HasLeast3Neighbors(image.Point{maxX, maxY}))
}

// HasLeast3Neighbors returns true if pixel has at least 3 neighbors.
func (img *Image) HasLeast3Neighbors(pt image.Point) bool {
	return img.SameColorNeighbors(pt) > 2
}

// SameColorNeighbors checks if a pixel has 3+ adjacent pixels of the
// same color.
// TODO use vector points?
func (img *Image) SameColorNeighbors(pt image.Point) int {
	neibrs := 0
	x1 := intMax(pt.X-1, 0)
	y1 := intMax(pt.Y-1, 0)
	x2 := intMin(pt.X+1, img.Bounds().Dx()-1)
	y2 := intMin(pt.Y+1, img.Bounds().Dy()-1)
	pos1 := img.Position(pt)

	if pt.X == x1 || pt.X == x2 || pt.Y == y1 || pt.Y == y2 {
		neibrs++
	}

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			if x == pt.X && y == pt.Y {
				continue
			}

			pos2 := img.Position(image.Point{x, y})
			bs := img.Bytes()
			if bs[pos1+0] == bs[pos2+0] &&
				bs[pos1+1] == bs[pos2+1] &&
				bs[pos1+2] == bs[pos2+2] &&
				bs[pos1+3] == bs[pos2+3] {
				neibrs++
			}
		}
	}
	return neibrs
}

// CreateOutput creates output image for given path in filesystem.
func CreateOutput(path string, w, h int) (*image.RGBA, error) {
	_, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	return img, nil
}
