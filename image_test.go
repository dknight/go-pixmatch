package pixmatch

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"reflect"
	"testing"
)

// false is only for debugging purposes.
var removeDiffImages = true

func TestNewImage(t *testing.T) {
	img := NewImage()
	want := "*pixmatch.Image"
	typ := reflect.TypeOf(img).String()

	if typ != want {
		t.Errorf("Expected %v got %v", want, typ)
	}
}

func TestNewFromPath(t *testing.T) {
	want := "./res/kitten-a.png"
	img, err := NewImageFromPath(want)
	if err != nil {
		t.Error(err)
	}
	if want != img.Path {
		t.Errorf("Expected %v got %v", want, img.Path)
	}
}

func TestNewFromPath_NotExists(t *testing.T) {
	path := "./res/kitten-xxx.png"
	_, err := NewImageFromPath(path)
	if err == nil {
		t.Error(err)
	}
}

func TestImageSetPath(t *testing.T) {
	want := "./res/kitten-a.png"
	img := NewImage()
	img.SetPath(want)
	if want != img.Path {
		t.Errorf("Expected %v got %v", want, img.Path)
	}
}

func TestImageSize(t *testing.T) {
	img, err := NewImageFromPath("./res/gray8-a.png")
	if err != nil {
		t.Error(err)
	}
	want := 256
	res := img.Size()
	if want != res {
		t.Errorf("Expected %v got %v", want, res)
	}
}

func TestImageLoad(t *testing.T) {
	path := "./res/kitten-a.png"
	_, err := NewImageFromPath(path)
	if err != nil {
		t.Errorf("File %v decoded incorrectly", path)
	}

	path = "./res/nonexists"
	_, err = NewImageFromPath(path)
	if err == nil {
		t.Errorf("File %v should not exist", path)
	}

	path = "./res/corrupted.png"
	_, err = NewImageFromPath(path)
	if err == nil {
		t.Errorf("File %v cannot be decoded correctly", path)
	}

	path = "./res/not-image"
	_, err = NewImageFromPath(path)
	if err != image.ErrFormat {
		t.Errorf("File %v is not an image", path)
	}
}

func TestImageEmpty(t *testing.T) {
	exp := true
	img := NewImage()
	if !img.Empty() {
		t.Errorf("Expected %v got %v", exp, img.Empty())
	}
}

func TestDimensionsEqual(t *testing.T) {
	paths := []string{
		"./res/kitten-a.png",
		"./res/kitten-b.png",
		"./res/kitten-c-small.png",
	}
	images := make([]*Image, 0, len(paths))

	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}

	ok := images[0].DimensionsEqual(images[1])
	if !ok {
		t.Errorf("Expected %v got %v", ok, !ok)
	}

	notOk := images[0].DimensionsEqual(images[2])
	if notOk {
		t.Error(ErrDimensionsDoNotMatch)
	}
}

func TestIdentical(t *testing.T) {
	paths := []string{
		"./res/kitten-a.png",
		"./res/kitten-b.png",
		"./res/kitten-a.png",
	}
	images := make([]*Image, 0, len(paths))
	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}

	want := false
	res := images[0].Identical(images[1])
	if res {
		t.Errorf("Expected %v got %v", want, res)
	}

	want = true
	res = images[0].Identical(images[2])
	if !res {
		t.Errorf("Expected %v got %v", want, res)
	}
}

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
		"./res/models/alpha.png":   16,
		"./res/models/alpha32.png": 16,
		"./res/models/tt.jpg":      64,
	}
	for path, bits := range pairs {
		img, _ := NewImageFromPath(path)
		bs := img.Bytes()
		if len(bs) != bits {
			t.Errorf("Expected %v (%vb) got %v", path, bits, len(bs))
		}
	}
}

func TestPosition(t *testing.T) {
	img, err := NewImageFromPath("./res/kitten-a.png")
	if err != nil {
		t.Error(nil)
	}
	point := image.Point{50, 50}
	res := img.Position(point)
	want := 20200
	if res != want {
		t.Errorf("Expected %v got %v", want, res)
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
		"./res/kitten-a.png",
		"./res/kitten-c-small.png",
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
		"./res/kitten-a.png",
		"./res/kitten-a.png",
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

func TestColorDelta(t *testing.T) {
	paths := []string{
		"./res/kitten-a.png",
		"./res/kitten-b.png",
	}
	images := make([]*Image, len(paths))
	for i, p := range paths {
		images[i], _ = NewImageFromPath(p)
	}

	pairs := map[image.Point]float64{
		image.Point{11, 58}: -6365.96249947337,
		image.Point{50, 50}: 0,
		image.Point{13, 16}: 15476.475726033921,
	}

	for pt, want := range pairs {
		pos := images[0].Position(pt)
		res := images[0].ColorDelta(images[1], pos, pos, false)
		if res != want {
			t.Errorf("Expected %v got %v", want, res)
		}
	}

	// Only Y (brigthness) component.
	pos := images[0].Position(image.Point{11, 58})
	res := images[0].ColorDelta(images[1], pos, pos, true)
	want := 111.97145726
	if res != want {
		t.Errorf("Expected %v got %v", want, res)
	}
}

func TestSameNeighbors(t *testing.T) {
	img, _ := NewImageFromPath("./res/kitten-a.png")
	pairs := map[image.Point]bool{
		image.Point{0, 99}:  true,
		image.Point{7, 99}:  true,
		image.Point{13, 72}: true,
		image.Point{17, 72}: true,
		image.Point{6, 70}:  true,
		image.Point{3, 10}:  true,
		image.Point{11, 14}: false,
	}

	for pt, want := range pairs {
		res := img.SameNeighbors(pt, 3)
		if res != want {
			t.Errorf("Expected %v got %v", want, res)
		}
	}
}

func TestHasLeast3Neighbors(t *testing.T) {
	img, _ := NewImageFromPath("./res/kitten-a.png")
	pairs := map[image.Point]bool{
		image.Point{0, 99}:  true,
		image.Point{3, 10}:  true,
		image.Point{11, 14}: false,
	}

	for pt, want := range pairs {
		res := img.SameNeighbors(pt, 3)
		if res != want {
			t.Errorf("Expected %v got %v", want, res)
		}
	}
}

func TestAntialiased(t *testing.T) {
	img, _ := NewImageFromPath("./res/kitten-a.png")
	pairs := map[image.Point]bool{
		image.Point{0, 0}:   false,
		image.Point{17, 61}: false,
		image.Point{12, 88}: true,
	}

	for pt, want := range pairs {
		res := img.Antialiased(img, pt)
		if res != want {
			t.Errorf("Expected %v got %v", want, res)
		}
	}
}

// -------------------------- PAIRS --------------------------------------

type testPair struct {
	name         string
	pathA        string
	pathB        string
	pathDiff     string
	expectedDiff int
	skip         bool
	options      *Options
}

var testPairs = []testPair{
	testPair{
		name:         "PNG+GRAY8",
		pathA:        "./res/gray8-a.png",
		pathB:        "./res/gray8-b.png",
		pathDiff:     "diff-gray8.png",
		expectedDiff: 4,
		skip:         false,
		options:      NewOptions().SetIncludeAA(true),
	},
	testPair{
		name:         "PNG+GRAY16",
		pathA:        "./res/gray16-a.png",
		pathB:        "./res/gray16-b.png",
		pathDiff:     "diff-gray16.png",
		expectedDiff: 15,
		skip:         false,
		options:      NewOptions().SetIncludeAA(true),
	},
	testPair{
		name:         "PNG+Alpha",
		pathA:        "./res/kitten-a.png",
		pathB:        "./res/kitten-b.png",
		pathDiff:     "diff-kitten.png",
		expectedDiff: 1706,
		skip:         false,
		options:      NewOptions(),
	},
	testPair{
		name:         "PNG+Alpha+Anti-aliasing",
		pathA:        "./res/kitten-a.png",
		pathB:        "./res/kitten-b.png",
		pathDiff:     "diff-kitten-aa.png",
		expectedDiff: 2090,
		skip:         false,
		options:      NewOptions().SetIncludeAA(true),
	},
	testPair{
		name:         "GIF",
		pathA:        "./res/abigail-a.gif",
		pathB:        "./res/abigail-b.gif",
		pathDiff:     "diff-abigail.gif",
		expectedDiff: 114,
		skip:         false,
		options:      NewOptions().SetIncludeAA(true),
	},
	testPair{
		name:         "JPEG",
		pathA:        "./res/forest-a.jpg",
		pathB:        "./res/forest-b.jpg",
		pathDiff:     "diff-forest.jpg",
		expectedDiff: 656,
		skip:         false,
		options:      NewOptions().SetAlpha(.8),
	},
	testPair{
		name:         "pixelmatch.js_100",
		pathA:        "./res/original/1a.png",
		pathB:        "./res/original/1b.png",
		pathDiff:     "./res/original/1diff~.png",
		expectedDiff: 143,
		skip:         false,
		options:      NewOptions().SetThreshold(.05),
	},
	testPair{
		name:         "pixelmatch.js_101",
		pathA:        "./res/original/1a.png",
		pathB:        "./res/original/1b.png",
		pathDiff:     "./res/original/1diffmask~.png",
		expectedDiff: 143,
		skip:         false,
		options:      NewOptions().SetThreshold(.05).SetDiffMask(true),
	},
	testPair{
		name:         "pixelmatch.js_102",
		pathA:        "./res/original/1a.png",
		pathB:        "./res/original/1b.png",
		pathDiff:     "./res/original/1emptydiffmask~.png",
		expectedDiff: 0,
		skip:         false,
		options:      NewOptions().SetThreshold(1).SetDiffMask(true),
	},
	testPair{
		name:         "pixelmatch.js_200",
		pathA:        "./res/original/2a.png",
		pathB:        "./res/original/2b.png",
		pathDiff:     "./res/original/2diff~.png",
		expectedDiff: 12437,
		skip:         false,
		options:      NewOptions().SetThreshold(.05).SetAlpha(.5).SetAAColor(color.RGBA{0, 192, 0, 255}).SetDiffColor(color.RGBA{255, 0, 255, 255}),
	},
	testPair{
		name:         "pixelmatch.js_300",
		pathA:        "./res/original/3a.png",
		pathB:        "./res/original/3b.png",
		pathDiff:     "./res/original/3diff~.png",
		expectedDiff: 212,
		skip:         false,
		options:      NewOptions().SetThreshold(.05),
	},
	testPair{
		name:         "pixelmatch.js_400",
		pathA:        "./res/original/4a.png",
		pathB:        "./res/original/4b.png",
		pathDiff:     "./res/original/4diff~.png",
		expectedDiff: 36049,
		skip:         false,
		options:      NewOptions().SetThreshold(.05),
	},
	testPair{
		name:         "pixelmatch.js_500",
		pathA:        "./res/original/5a.png",
		pathB:        "./res/original/5b.png",
		pathDiff:     "./res/original/5diff~.png",
		expectedDiff: 0, // 0 because no AA included
		skip:         false,
		options:      NewOptions().SetThreshold(.05),
	},
	testPair{
		name:         "pixelmatch.js_600",
		pathA:        "./res/original/6a.png",
		pathB:        "./res/original/6b.png",
		pathDiff:     "./res/original/6diff~.png",
		expectedDiff: 51,
		skip:         false,
		options:      NewOptions().SetThreshold(.05),
	},
	testPair{
		name:         "pixelmatch.js_601",
		pathA:        "./res/original/6a.png",
		pathB:        "./res/original/6b.png",
		pathDiff:     "./res/original/6empty~.png",
		expectedDiff: 0,
		skip:         false,
		options:      NewOptions().SetThreshold(1.0),
	},
	testPair{
		name:         "pixelmatch.js_700",
		pathA:        "./res/original/7a.png",
		pathB:        "./res/original/7b.png",
		pathDiff:     "./res/original/7diff~.png",
		expectedDiff: 2448,
		skip:         false,
		options:      NewOptions().SetDiffColorAlt(color.RGBA{0, 255, 0, 255}),
	},
}

func TestImageCompare(t *testing.T) {
	for _, pair := range testPairs {
		t.Run(pair.name, func(t *testing.T) {
			if pair.skip {
				t.SkipNow()
			}
			t.Cleanup(func() {
				if removeDiffImages {
					os.Remove(pair.pathDiff)
				}
			})
			imageA, err := NewImageFromPath(pair.pathA)
			if err != nil {
				t.Error(err)
			}
			imageB, err := NewImageFromPath(pair.pathB)
			if err != nil {
				t.Error(err)
			}
			fp, err := os.Create(pair.pathDiff)
			if err != nil {
				t.Error(err)
			}
			output, err := NewOutput(fp, imageA.Rect.Dx(), imageB.Rect.Dy())
			if err != nil {
				t.Error(err)
			}
			pair.options.SetOutput(output)
			diff, err := imageA.Compare(imageB, pair.options)
			if err != nil {
				t.Error("Compare", err.Error())
			}
			if diff != pair.expectedDiff {
				t.Errorf("Expected %v got %v", pair.expectedDiff, diff)
			}
		})
	}
}
