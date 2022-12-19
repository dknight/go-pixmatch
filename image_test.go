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
		t.Error(err)
	}
}

func TestLoadImages(t *testing.T) {
	f1 := "./res/kitten1.png"
	f2 := "./res/kitten-small.png"
	f3 := "./res/corrupted.png"

	img1 := NewImage()
	img1.SetPath(f1)
	img2 := NewImage()
	img2.SetPath(f2)
	img3 := NewImage()
	img3.SetPath(f3)

	images := [ImagesCount]*Image{img1, img2}
	err := LoadImages(images)
	if err != nil {
		t.Error(err)
	}

	images = [ImagesCount]*Image{img3, img2}
	err = LoadImages(images)
	if err != nil {
		if err != io.ErrUnexpectedEOF && err != io.EOF {
			t.Error("File", f3, "cannot be decoded, file is corrupted.")
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
	img1 := NewImage()
	img1.SetPath("./res/kitten1.png")

	img2 := NewImage()
	img2.SetPath("./res/kitten2.png")

	img3 := NewImage()
	img3.SetPath("./res/kitten-small.png")

	images := [ImagesCount]*Image{img1, img2}
	err := LoadImages(images)
	if err != nil {
		t.Error(err)
	}
	exp, _ := images[0].DimensionsEqual(images[1])
	if !exp {
		t.Error("Expected", exp, "got", false)
	}

	images = [ImagesCount]*Image{img2, img3}
	err = LoadImages(images)
	if err != nil {
		t.Error(err)
	}

	exp, err = images[0].DimensionsEqual(images[1])
	if exp && err == ErrDimensionsDoNotMatch {
		t.Error(ErrDimensionsDoNotMatch)
	}
}
