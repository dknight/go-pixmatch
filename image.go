package pixmatch

import (
	"bytes"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"sync"
)

// ImagesCount total images to compare probably won't ever change.
const ImagesCount = 2

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

// SetPath sets the path to the image in filesystem.
func (img *Image) SetPath(path string) {
	img.Path = path
}

// DimensionsEqual checks that dimensions of the image is equal to dimension
// of other image.
func (img1 *Image) DimensionsEqual(img2 *Image) (bool, error) {
	if img1.Bounds().Eq(img2.Bounds()) {
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

// LoadImages loads multiple (2) images asynchronously.
func LoadImages(images [ImagesCount]*Image) error {
	ch := make(chan error)
	var wg sync.WaitGroup
	wg.Add(ImagesCount)

	for i := 0; i < ImagesCount; i++ {
		go func(i int) {
			ch <- images[i].Load()
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for e := range ch {
		if e != nil {
			return e
		}
	}

	return nil
}

// Load loads image from the path of filesystem.
func (img *Image) Load() error {
	fp, err := os.Open(img.Path)
	if err != nil {
		return err
	}
	defer fp.Close()

	img.Image, img.Format, err = image.Decode(fp)
	if err != nil {
		return err
	}
	return nil
}

// Identical checks images that these are identical on bytes level.
// Looks like bytes.Equal() is the fastest way to compare 2 bytes arrays.
// Tried reflect.DeepEqual() and loop solutions. In the most cases
// bytes.Equal() is the best choice.
func (img1 *Image) Identical(img2 *Image) bool {
	return bytes.Equal(img1.Bytes(), img2.Bytes())
}

// Bytes get the raw bytes of the image.
// NOTE Is there any better way to make this in better way?
// TODO make jpeg
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
		// case color.NYCbCrAModel:
		// 	return img.Image.(*image.NYCbCrA).Y
		// case color.YCbCrModel:
		// 	return img.Image.(*image.YCbCr).Y
	}

	switch img.ColorModel().(type) {
	case color.Palette:
		return img.Image.(*image.Paletted).Pix
	}

	return []byte{}
}
