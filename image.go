package pixmatch

import (
	"bytes"
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"reflect"
	"sync"
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

	// Pixel data contains data for colors.
	PixData []uint32

	// Rect is as cache for rectangle to avoid extra calculations.
	Rect image.Rectangle

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
	// Cache some stuff
	img.PixData = img.Uint32()
	img.Rect = img.Bounds()

	return
}

// Compare returns the number of different pixels.
func (img *Image) Compare(img2 *Image, opts *Options) (int, error) {
	if opts == nil {
		opts = NewOptions()
	}

	// If empty images return error.
	if img.Empty() || img2.Empty() {
		return ExitEmptyImage, ErrImageIsEmpty
	}

	// If dimensions do not match, return error.
	if !img.DimensionsEqual(img2) {
		return ExitDimensionsNotEqual, ErrDimensionsDoNotMatch
	}

	// If bytes are the same just return nothing to compare more.
	if img.Identical(img2) {
		// NOTE
		// We don't work to generate output image if it has no differences.
		// but original pixelmatch.js has it, maybe add later extra
		// option for this.
		return 0, nil
	}

	maxDelta := YIQDeltaMax * opts.Threshold * opts.Threshold
	bpc := img.BytesPerColor()
	diffColor := opts.ResolveDiffColor()
	diff := 0

	var wg sync.WaitGroup
	var mu sync.Mutex
	for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
		wg.Add(1)
		go func(y int) {
			for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
				point := image.Pt(x, y)
				pos := img.Position(point)
				delta := ColorDelta(img.PixData, img2.PixData, pos, pos, bpc, false)

				if math.Abs(delta) > maxDelta {
					if opts.DetectAA && (img.Antialiased(img2, point) || img2.Antialiased(img, point)) {
						if opts.Output != nil && !opts.DiffMask {
							opts.Output.Image.Set(x, y, opts.AAColor)
						}
					} else {
						if opts.Output != nil {
							opts.Output.Image.Set(x, y, diffColor)
						}
						mu.Lock()
						diff++
						mu.Unlock()
					}
				} else if opts.Output != nil && !opts.DiffMask {
					r, g, b, a := img.At(x, y).RGBA()
					gray := NewColor(r, g, b, a).BlendToGray(opts.Alpha)
					opts.Output.Image.Set(x, y, gray)
				}
			}
			wg.Done()
		}(y)
	}

	wg.Wait()

	if opts.Output != nil {
		err := opts.Output.Save(img.Format)
		if err != nil {
			return -1, err
		}
	}
	return diff, nil
}

// Empty checks that images is empty of has theoretical size 0x0 pixels.
func (img *Image) Empty() bool {
	if img.Image == nil {
		return true
	}
	return img.Rect.Empty()
}

// DimensionsEqual checks that dimensions of the image is equal to dimension
// of other image.
func (img *Image) DimensionsEqual(img2 *Image) bool {
	return img.Rect.Eq(img2.Rect)
}

// Identical checks images that these are identical on bytes level.
// Looks like bytes.Equal() is the fastest way to compare 2 bytes arrays.
// Tried reflect.DeepEqual() and loop solutions. In the most cases
// bytes.Equal() is the best choice.
//
//	loops - the slowest
//	reflect.DeepEqual() - slower
//	bytes.Compare() - better
//	bytes.Equal() - even better
func (img *Image) Identical(img2 *Image) bool {
	return bytes.Equal(img.Bytes(), img2.Bytes())
}

// Bytes get the raw bytes of the pixel data.
// Reflection is never clear (https://go-proverbs.github.io/)
func (img *Image) Bytes() []byte {
	val := reflect.ValueOf(img.Image)
	ptr := reflect.Indirect(val)
	pixs := ptr.FieldByName("Pix")
	// return empty byte slice if something is really wrong with the image.
	if pixs.IsValid() {
		return pixs.Bytes()
	}
	// For JPEG
	y := ptr.FieldByName("Y")
	if y.IsValid() {
		return y.Bytes()
	}
	return nil
}

// Stride get generic stride. Default return value is zero.
// Reflection is never clear (https://go-proverbs.github.io/)
func (img *Image) Stride() int {
	val := reflect.ValueOf(img.Image)
	ptr := reflect.Indirect(val)
	stride := ptr.FieldByName("Stride")
	if stride.IsValid() {
		return int(stride.Int()) // int64
	}

	// for JPEG (NOTE it is very dirty)
	strideY := ptr.FieldByName("YStride")
	if strideY.IsValid() {
		return int(strideY.Int())
	}
	return 0
}

// Position is the positions of the pixel in bytes data.
// Eg. x1 is (r + 0), (g + 1), (b + 2), (a + 3), x2 = ...
func (img *Image) Position(p image.Point) int {
	return (p.Y-img.Rect.Min.Y)*img.Stride() +
		(p.X-img.Rect.Min.X)*img.BytesPerColor()
}

// BytesPerColor resolves count of bytes per color.
func (img *Image) BytesPerColor() int {
	switch img.ColorModel() {
	case color.AlphaModel, color.GrayModel:
		return 1
	case color.Alpha16Model, color.Gray16Model:
		return 2
	case color.CMYKModel, color.NRGBAModel, color.RGBAModel:
		return 4
	case color.NRGBA64Model, color.RGBA64Model:
		return 8
	}

	// NOTE need?
	switch img.Image.(type) {
	case *image.YCbCr:
	case *image.Paletted:
		return 1
	}

	return 1 // default for any other possible case if exists...
}

// ColorDelta is the squared YUV distance between colors at this pixel
// position, returns negative if the img2 pixel is darker.
// If argument onlyY is true, the only brightness level will be returned
// (Y component of YIQ model).
func ColorDelta(pix1, pix2 []uint32, m, n int, bpc int, onlyY bool) float64 {
	var r1, g1, b1, a1 uint32
	var r2, g2, b2, a2 uint32
	switch bpc {
	case 1:
		r1, g1, b1, a1 = pix1[m], pix1[m], pix1[m], pix1[m]
		r2, g2, b2, a2 = pix2[n], pix2[n], pix2[n], pix2[n]
	case 2:
		r1, g1, b1, a1 = pix1[m], pix1[m], pix1[m], pix1[m+1]
		r2, g2, b2, a2 = pix2[n], pix2[n], pix2[n], pix2[n+1]
	default:
	case 4:
		r1, g1, b1, a1 = pix1[m], pix1[m+1], pix1[m+2], pix1[m+3]
		r2, g2, b2, a2 = pix2[n], pix2[n+1], pix2[n+2], pix2[n+3]
	case 8:
		r1, r2 = pix1[0]<<8|pix1[1], pix2[0]<<8|pix2[1]
		g1, g2 = pix1[2]<<8|pix1[3], pix2[2]<<8|pix2[3]
		b1, b2 = pix1[4]<<8|pix1[5], pix2[4]<<8|pix2[5]
		a1, a2 = pix1[6]<<8|pix1[7], pix2[6]<<8|pix2[7]
	}

	color1 := NewColor(r1, g1, b1, a1)
	color2 := NewColor(r2, g2, b2, a2)

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
	delta := 0.5053*y*y + 0.299*i*i + 0.1957*q*q

	if y1 > y2 {
		return -delta
	}
	return delta
}

// Uint32 converts bytes array into []uint32 slice.
func (img *Image) Uint32() []uint32 {
	bs := img.Bytes()
	ui32 := make([]uint32, len(bs))
	for i, b := range bs {
		ui32[i] = uint32(b)
	}
	return ui32
}

// Antialiased checks that point is anti-aliased.
func (img *Image) Antialiased(img2 *Image, pt image.Point) bool {
	neibrs := 0
	x1 := intMax(pt.X-1, img.Rect.Min.X)
	y1 := intMax(pt.Y-1, img.Rect.Min.Y)
	x2 := intMin(pt.X+1, img.Rect.Max.X-1)
	y2 := intMin(pt.Y+1, img.Rect.Max.Y-1)
	pos := img.Position(pt)

	if pt.X == x1 || pt.X == x2 || pt.Y == y1 || pt.Y == y2 {
		neibrs++
	}

	min, max := 0.0, 0.0
	minX, minY, maxX, maxY := 0, 0, 0, 0
	bpc := img.BytesPerColor()

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			if x == pt.X && y == pt.Y {
				continue
			}

			pos2 := img.Position(image.Pt(x, y))
			delta := ColorDelta(img.PixData, img.PixData, pos, pos2, bpc, true)
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

	return (img.SameNeighbors(image.Pt(minX, minY), 3) &&
		img2.SameNeighbors(image.Pt(minX, minY), 3)) ||
		(img.SameNeighbors(image.Pt(maxX, maxY), 3) &&
			img2.SameNeighbors(image.Pt(maxX, maxY), 3))
}

// SameNeighbors checks if a pixel has 3+ adjacent pixels of the
// same color.
func (img *Image) SameNeighbors(pt image.Point, n int) bool {
	neibrs := 0
	x1 := intMax(pt.X-1, img.Rect.Min.X)
	y1 := intMax(pt.Y-1, img.Rect.Min.Y)
	x2 := intMin(pt.X+1, img.Rect.Max.X-1)
	y2 := intMin(pt.Y+1, img.Rect.Max.Y-1)
	pos1 := img.Position(pt)

	if pt.X == x1 || pt.X == x2 || pt.Y == y1 || pt.Y == y2 {
		neibrs++
	}

	bs := img.Bytes()
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			if x == pt.X && y == pt.Y {
				continue
			}

			pos2 := img.Position(image.Pt(x, y))
			ok := true
			for i := 0; i < img.BytesPerColor(); i++ {
				if bs[pos1+i] != bs[pos2+i] {
					ok = false
					break
				}
			}
			if ok {
				neibrs++
			}
			if neibrs >= n {
				return true
			}
		}
	}
	return false
}
