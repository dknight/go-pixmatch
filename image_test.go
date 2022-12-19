package pixmatch

import (
	"image"
	"io"
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

func TestImageSetPath(t *testing.T) {
	exp := "./hello.png"
	img := NewImage()
	img.SetPath(exp)
	if exp != img.Path {
		t.Error("Expected", exp, "got", img.Path)
	}
}

func TestImageLoad(t *testing.T) {
	path := "./res/kitten1.png"
	img := NewImage()
	img.SetPath(path)
	err := img.Load()
	if err != nil {
		t.Error("File", path, "decoded incorrectly")
	}

	path = "./res/nonexists"
	img = NewImage()
	img.SetPath(path)
	err = img.Load()
	if err == nil {
		t.Error("File", path, "should not exist")
	}

	path = "./res/corrupted.png"
	img = NewImage()
	img.SetPath(path)
	err = img.Load()
	if err == nil {
		t.Error("File", path, "cannot be decoded correctly.")
	}

	path = "./res/nonimg.txt"
	img = NewImage()
	img.SetPath(path)
	err = img.Load()
	if err != image.ErrFormat {
		t.Error("File", path, "is not an image.")
	}
}

func TestLoadImages(t *testing.T) {
	paths := []string{
		"./res/kitten1.png",
		"./res/kitten-small.png",
		"./res/corrupted.png",
	}
	images := make([]*Image, 3)

	for i, p := range paths {
		images[i] = NewImage()
		images[i].SetPath(p)
	}
	imgs := [ImagesCount]*Image{images[0], images[1]}
	err := LoadImages(imgs)
	if err != nil {
		t.Error(err)
	}

	imgs = [ImagesCount]*Image{images[2], images[1]}
	err = LoadImages(imgs)
	if err != nil {
		if err != io.ErrUnexpectedEOF && err != io.EOF {
			t.Error("File", images[2].Path,
				"cannot be decoded, file is corrupted.")
		}
	}
}

func TestImageEmpty(t *testing.T) {
	exp := true
	img := NewImage()
	if !img.Empty() {
		t.Error("Expected", exp, "got", img.Empty())
	}

	exp = false
	img.SetPath("./res/1x1.png")
	err := img.Load()
	if err != nil {
		t.Error(err)
	}
	if img.Empty() {
		t.Error("Expected", exp, "got", img.Empty())
	}
}

func TestDimensionsEqual(t *testing.T) {
	paths := []string{
		"./res/kitten1.png",
		"./res/kitten2.png",
		"./res/kitten-small.png",
	}
	images := make([]*Image, 3)

	for i, p := range paths {
		images[i] = NewImage()
		images[i].SetPath(p)
		images[i].Load()
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
	images := make([]*Image, len(paths))
	for i, p := range paths {
		images[i] = NewImage()
		images[i].SetPath(p)
		images[i].Load()
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
	images := make([]*Image, len(pairs))
	i := 0
	for k, bits := range pairs {
		images[i] = NewImage()
		images[i].SetPath(k)
		_ = images[i].Load()
		bs := images[i].Bytes()
		if len(bs) != bits {
			t.Error("Expected", k, bits, "got", len(bs))
		}
		i++
	}
}
