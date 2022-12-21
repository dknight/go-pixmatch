package pixmatch

import (
	"image"
	"reflect"
	"testing"
)

func TestNewImage(t *testing.T) {
	img := NewImage()
	exp := "*pixmatch.Image"
	typ := reflect.TypeOf(img).String()

	if typ != exp {
		t.Error("Expected type", exp, "got", typ)
	}
}

func TestNewFromPath(t *testing.T) {
	path := "./res/kitten1.png"
	img, err := NewImageFromPath(path)
	if err != nil {
		t.Error(err)
	}
	if path != img.Path {
		t.Error("Expected", path, "got", img.Path)
	}
}

func TestNewFromPath_NotExists(t *testing.T) {
	path := "./res/kitten3.png"
	_, err := NewImageFromPath(path)
	if err == nil {
		t.Error(err)
	}
}

func TestImageSetPath(t *testing.T) {
	path := "./res/kitten1.png"
	img := NewImage()
	img.SetPath(path)
	if path != img.Path {
		t.Error("Expected", path, "got", img.Path)
	}
}

func TestImageLoad(t *testing.T) {
	path := "./res/kitten1.png"
	_, err := NewImageFromPath(path)
	if err != nil {
		t.Error("File", path, "decoded incorrectly")
	}

	path = "./res/nonexists"
	_, err = NewImageFromPath(path)
	if err == nil {
		t.Error("File", path, "should not exist")
	}

	path = "./res/corrupted.png"
	_, err = NewImageFromPath(path)
	if err == nil {
		t.Error("File", path, "cannot be decoded correctly")
	}

	path = "./res/nonimg.txt"
	_, err = NewImageFromPath(path)
	if err != image.ErrFormat {
		t.Error("File", path, "is not an image")
	}
}

func TestImageEmpty(t *testing.T) {
	exp := true
	img := NewImage()
	if !img.Empty() {
		t.Error("Expected", exp, "got", img.Empty())
	}
}

func TestDimensionsEqual(t *testing.T) {
	paths := []string{
		"./res/kitten1.png",
		"./res/kitten2.png",
		"./res/kitten-small.png",
	}
	images := make([]*Image, 0, len(paths))

	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}

	exp, _ := images[0].DimensionsEqual(images[1])
	if !exp {
		t.Error("Expected", exp, "got", false)
	}

	exp, err := images[0].DimensionsEqual(images[2])
	if exp && err == ErrDimensionsDoNotMatch {
		t.Error(ErrDimensionsDoNotMatch)
	}
}

func TestIdentical(t *testing.T) {
	paths := []string{
		"./res/kitten1.png",
		"./res/kitten2.png",
		"./res/kitten1.png",
	}
	images := make([]*Image, 0, len(paths))
	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}

	exp := false
	res := images[0].Identical(images[1])
	if res {
		t.Error("Expected", exp, "got", res)
	}

	exp = true
	res = images[0].Identical(images[2])
	if !res {
		t.Error("Expected", exp, "got", res)
	}
}

// TODO check Alpha, JPEG
func TestBytes(t *testing.T) {
	pairs := map[string]int{
		"./res/models/nrgba.png":   16,
		"./res/models/nrgba32.png": 32,
		"./res/models/rgb.png":     16,
		"./res/models/rgb32.png":   32,
		"./res/models/gray.png":    4,
		"./res/models/gray32.png":  8,
		"./res/models/graya.png":   16,
		"./res/models/graya32.png": 32,
		"./res/models/palette.png": 4,
		"./res/models/alpha.png":   16, //FIXME
		"./res/models/alpha32.png": 16, //FIXME
		"./res/models/tt.jpg":      0,  //FIXME
	}
	for path, bits := range pairs {
		img, _ := NewImageFromPath(path)
		bs := img.Bytes()
		if len(bs) != bits {
			t.Error("Expected", path, bits, "got", len(bs))
		}
	}
}

func TestPosition(t *testing.T) {
	img, err := NewImageFromPath("./res/kitten1.png")
	if err != nil {
		t.Error(nil)
	}
	res := img.Position(50, 50)
	exp := 20200
	if res != exp {
		t.Errorf("Expecte %d got %d", exp, res)
	}
}

func TestCompare_Empty(t *testing.T) {
	imgEmpty1 := NewImage()
	imgEmpty2 := NewImage()
	px, err := imgEmpty1.Compare(imgEmpty2, nil)
	if px != ExitEmptyImage && err != nil {
		t.Error("Images should be empty")
	}
}

func TestCompare_Dimensions(t *testing.T) {
	paths := []string{
		"./res/kitten1.png",
		"./res/kitten-small.png",
	}
	images := make([]*Image, len(paths))
	for i, p := range paths {
		images[i], _ = NewImageFromPath(p)
	}
	s, err := images[0].Compare(images[1], nil)
	if err != nil && s != ExitDimensionsNotEqual {
		t.Error(err)
	}
}

func TestCompare_Identical(t *testing.T) {
	paths := []string{
		"./res/kitten1.png",
		"./res/kitten1.png",
	}
	images := make([]*Image, len(paths))
	for i, p := range paths {
		images[i], _ = NewImageFromPath(p)
	}
	px, err := images[0].Compare(images[1], nil)
	if px != 0 || err != nil {
		t.Error("Images should be identical")
	}
}
