package pixmatch

import (
	"bytes"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

// Common constants for pixmatch package.
const (
	// YIQDeltaMax is the value of 35215. This is the maximum possible value
	// for the YIQ difference metric.
	// Read more about YIQ NTSC https://en.wikipedia.org/wiki/YIQ
	YIQDeltaMax = 35215

	// DefaultFormat is used if format is not specified.
	DefaultFormat = FormatPNG

	FormatPNG  = "png"
	FormatGIF  = "gif"
	FormatJPEG = "jpeg"
)

// Image represents the image structure.
type Image struct {
	// Path to the image in file system.
	Path string

	// Format as a string like (PNG, JEPG, GIF).
	Format string

	// PixData contains data for colors as uint32 numbers.
	PixData []uint32

	// BPC is the number of bytes per color.
	BPC int

	// Image is an embedded [image.Image] from the standard library.
	image.Image
}

// NewImage creates a new image instance.
func NewImage(w, h int, format string) *Image {
	if format == "" {
		format = DefaultFormat
	}
	return &Image{
		Image:  image.NewRGBA(image.Rect(0, 0, w, h)),
		Format: format,
	}
}

// NewImageFromPath creates a new image instance from the file system path.
func NewImageFromPath(path string) (*Image, error) {
	ext := filepath.Ext(path)
	format := FormatPNG
	switch ext {
	case ".png":
		format = FormatPNG
	case ".gif":
		format = FormatGIF
	case ".jpeg", ".jpg":
		format = FormatJPEG
	}

	img := NewImage(0, 0, format)
	img.Path = path

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

// Load reads data from the reader.
func (img *Image) Load(rd io.Reader) (err error) {
	img.Image, img.Format, err = image.Decode(rd)
	if err != nil {
		return err
	}
	// Cache pixel data because Uint32() is very expensive.
	img.PixData = img.Uint32()
	img.BPC = img.BytesPerColor()
	return
}

// Size gives the total size of the image in pixels.
func (img *Image) Size() int {
	return img.Bounds().Dx() * img.Bounds().Dy()
}

// Compare returns the number of different pixels between two comparable
// images. Zero is returned if no difference found.Returns negative values
// if something went wrong but in this case error also returned.
//
// Looks like process row of the pixel in a single goroutine is the most
// performant way to do this, but I can mistake here.
func (img *Image) Compare(img2 *Image, opts *Options) (int, error) {
	if opts == nil {
		opts = NewOptions()
	}

	// If empty images return error.
	if img.Empty() || img2.Empty() {
		return -1, ErrImageIsEmpty
	}

	// If dimensions do not match return error.
	if !img.DimensionsEqual(img2) {
		return -1, ErrDimensionsDoNotMatch
	}

	// If bytes are the same just return nothing to compare more.
	if img.Identical(img2) {
		return 0, nil
	}

	maxDelta := YIQDeltaMax * opts.Threshold * opts.Threshold
	diff := 0
	output := NewImage(img.Bounds().Dx(), img.Bounds().Dy(), img.Format)

	// Looks like the mutex + WaitGroup is the fastest found solution by me.
	// sync/atomic also shows the same results.
	var wg sync.WaitGroup
	var mu sync.Mutex
	wg.Add(img.Bounds().Max.Y)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		go func(y int) {
			defer wg.Done()
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				point := image.Pt(x, y)
				pos := img.Position(point)
				delta := img.ColorDelta(img2, pos, pos, false)

				if math.Abs(delta) > maxDelta {
					if !opts.IncludeAA &&
						(img.Antialiased(img2, point) ||
							img2.Antialiased(img, point)) {
						if opts.Output != nil && !opts.DiffMask {
							output.Image.(*image.RGBA).Set(x, y, opts.AAColor)
						}
					} else {
						if opts.Output != nil {
							diffColor := getDiffColor(opts, delta)
							output.Image.(*image.RGBA).Set(x, y, diffColor)
						}
						mu.Lock()
						diff++
						mu.Unlock()
					}
				} else if opts.Output != nil && !opts.DiffMask {
					r, g, b, a := img.At(x, y).RGBA()
					gray := NewColor(r, g, b, a).BlendToGray(opts.Alpha)
					output.Image.(*image.RGBA).Set(x, y, gray)
				}
			}
		}(y)
	}

	wg.Wait()

	// If no output given or there is no difference do not create diff file.
	if opts.Output != nil {
		err := output.Save(opts.Output)
		if err != nil {
			return -1, err
		}
	}
	return diff, nil
}

// getDiffColor get diff color.
func getDiffColor(opts *Options, delta float64) color.Color {
	diffColor := opts.DiffColor
	if delta < 0 && opts.DiffColorAlt != nil {
		diffColor = opts.DiffColorAlt
	}
	return diffColor
}

// Save encodes and writes image data to the destination.
func (img *Image) Save(wr io.Writer) (err error) {
	switch img.Format {
	case FormatGIF:
		err = gif.Encode(wr, img.Image, nil)
	case FormatJPEG:
		err = jpeg.Encode(wr, img.Image, nil)
	case FormatPNG:
		err = png.Encode(wr, img.Image)
	default:
		err = ErrUnknownFormat
	}
	return
}

// Empty checks that the image is empty or has a theoretical size of 0 pixels.
func (img *Image) Empty() bool {
	// Some failing happens on benchmarking, how 'img' can be nil
	// I have no idea.
	if img == nil {
		return true
	}
	return img.Bounds().Empty()
}

// DimensionsEqual checks that the dimensions of the two images are equal.
func (img *Image) DimensionsEqual(img2 *Image) bool {
	return img.Bounds().Eq(img2.Bounds())
}

// Identical determines whether or not images are identical at the byte level.
// This means that all the bytes of both images are the same.
//
// Alternative possible ways to compare:
//
//	loops - the slowest (faster for smaller images)
//	 reflect.DeepEqual() - slower
//	 bytes.Compare() - better
//	 bytes.Equal() - even better
func (img *Image) Identical(img2 *Image) bool {
	return bytes.Equal(img.Bytes(), img2.Bytes())
}

// Bytes are the raw bytes of the pixel data.
// Reflection is never clear (https://go-proverbs.github.io/)
func (img *Image) Bytes() []byte {
	val := reflect.ValueOf(img.Image)
	ptr := reflect.Indirect(val)
	pixs := ptr.FieldByName("Pix")
	if pixs.IsValid() {
		return pixs.Bytes()
	}

	// For JPEG
	y := ptr.FieldByName("Y")
	if y.IsValid() {
		return y.Bytes()
	}
	return []byte{}
}

// Stride gets the stride from the image. The default value is 1.
// Reflection is never clear (https://go-proverbs.github.io/)
func (img *Image) Stride() int {
	val := reflect.ValueOf(img.Image)
	ptr := reflect.Indirect(val)
	stride := ptr.FieldByName("Stride")
	if stride.IsValid() {
		return int(stride.Int()) // int64
	}

	// for JPEG
	strideY := ptr.FieldByName("YStride")
	if strideY.IsValid() {
		return int(strideY.Int())
	}
	return 1
}

// Position is the position of the pixel in the array of bytes.
//
// Formula
//
//	(y2-y1)*Stride + (x2-x1)*BPC
func (img *Image) Position(p image.Point) int {
	return (p.Y-img.Bounds().Min.Y)*img.Stride() +
		(p.X-img.Bounds().Min.X)*img.BPC
}

// BytesPerColor resolves the count of the bytes per color: 1, 2, 4, or 8.
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

	switch img.Image.(type) {
	case *image.Paletted:
		return 1
	}
	// The last frontier for any other possible cases.
	return 1
}

// ColorDelta is the squared YUV distance between colors at the pixel's
// position, returns a negative value if the img2 pixel is darker, and vice
// versa. If the argument onlyY is true, the only brightness level will be
// returned (Y component of the YIQ color space).
func (img *Image) ColorDelta(img2 *Image, m, n int, onlyY bool) float64 {
	px1 := img.PixData
	px2 := img2.PixData
	var r1, g1, b1, a1 uint32
	var r2, g2, b2, a2 uint32

	switch img.BPC {
	case 1:
		r1, g1, b1, a1 = px1[m], px1[m], px1[m], px1[m]
		r2, g2, b2, a2 = px2[n], px2[n], px2[n], px2[n]
	case 2:
		r1, g1, b1, a1 = px1[m], px1[m], px1[m], px1[m+1]
		r2, g2, b2, a2 = px2[n], px2[n], px2[n], px2[n+1]
	case 4:
		r1, g1, b1, a1 = px1[m], px1[m+1], px1[m+2], px1[m+3]
		r2, g2, b2, a2 = px2[n], px2[n+1], px2[n+2], px2[n+3]
	// NOTE not sure about this
	case 8:
		r1, r2 = px1[0]*px1[1], px2[0]*px2[1]
		g1, g2 = px1[2]*px1[3], px2[2]*px2[3]
		b1, b2 = px1[4]*px1[5], px2[4]*px2[5]
		a1, a2 = px1[6]*px1[7], px2[6]*px2[7]
	}

	switch img.Image.(type) {
	case *image.Paletted:
		x := px1[m]
		y := px2[n]
		palette1 := img.Image.(*image.Paletted).Palette
		palette2 := img2.Image.(*image.Paletted).Palette
		r1, g1, b1, a1 = palette1[x].RGBA()
		r2, g2, b2, a2 = palette2[y].RGBA()
		r1 >>= 8
		g1 >>= 8
		b1 >>= 8
		a1 >>= 8
		r2 >>= 8
		g2 >>= 8
		b2 >>= 8
		a2 >>= 8
		// case *image.YCbCr, *image.NYCbCrA:
		// Maybe do something here later...
	}

	color1 := NewColor(r1, g1, b1, a1)
	color2 := NewColor(r2, g2, b2, a2)

	// If all colors are the same then zero delta.
	if color1.Equals(color2) {
		return 0
	}

	if color1.A < 0xff {
		color1 = color1.Blend(float64(color1.A) / 0xff)
	}

	if color2.A < 0xff {
		color2 = color2.Blend(float64(color2.A) / 0xff)
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

// Uint32 converts image.Bytes() into a []uint32 slice. Be careful; this might
// be an expensive operation, used once and cached in image.PixData on image
// loading
func (img *Image) Uint32() []uint32 {
	bs := img.Bytes()
	ui32 := make([]uint32, len(bs))
	for i, b := range bs {
		ui32[i] = uint32(b)
	}
	return ui32
}

// Antialiased checks that the point is anti-aliased.
//
// NOTE Probably, better algorithms are required here.
func (img *Image) Antialiased(img2 *Image, pt image.Point) bool {
	neibrs := 0
	n := 2
	x1 := intMax(pt.X-1, img.Bounds().Min.X)
	y1 := intMax(pt.Y-1, img.Bounds().Min.Y)
	x2 := intMin(pt.X+1, img.Bounds().Max.X-1)
	y2 := intMin(pt.Y+1, img.Bounds().Max.Y-1)
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

			pos2 := img.Position(image.Pt(x, y))
			delta := img.ColorDelta(img, pos, pos2, true)
			if delta == 0 {
				neibrs++
				if neibrs > n {
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

	return (img.SameNeighbors(image.Pt(minX, minY), n) &&
		img2.SameNeighbors(image.Pt(minX, minY), n)) ||
		(img.SameNeighbors(image.Pt(maxX, maxY), n) &&
			img2.SameNeighbors(image.Pt(maxX, maxY), n))
}

// SameNeighbors determines whether a pixel has n+ adjacent pixels that are
// the same color.
func (img *Image) SameNeighbors(pt image.Point, n int) bool {
	neibrs := 0
	x1 := intMax(pt.X-1, img.Bounds().Min.X)
	y1 := intMax(pt.Y-1, img.Bounds().Min.Y)
	x2 := intMin(pt.X+1, img.Bounds().Max.X-1)
	y2 := intMin(pt.Y+1, img.Bounds().Max.Y-1)
	pos1 := img.Position(pt)

	if pt.X == x1 || pt.X == x2 || pt.Y == y1 || pt.Y == y2 {
		neibrs++
	}

	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			if x == pt.X && y == pt.Y {
				continue
			}

			pos2 := img.Position(image.Pt(x, y))
			ok := true
			for i := 0; i < img.BPC; i++ {
				if img.PixData[pos1+i] != img.PixData[pos2+i] {
					ok = false
					break
				}
			}
			if ok {
				neibrs++
			}
			if neibrs > n {
				return true
			}
		}
	}
	return false
}
