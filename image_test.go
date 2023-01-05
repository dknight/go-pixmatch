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

// NOTE
// Environemnt variable UPDATEDIFFS updates test instances.
// Use it only if you want to update diff or to be sure that everything
// is fine. Be careful about when updated. Test results may fail in future
// and cause a mess.

// Set environment variable UPDATEDIFFS=1 to update test samples.
var updateDiffImages bool

func init() {
	if _, ok := os.LookupEnv("UPDATEDIFFS"); ok {
		updateDiffImages = true
	}
}

func TestNewImage(t *testing.T) {
	img := NewImage(10, 10, DefaultFormat)
	want := "*pixmatch.Image"
	samples := reflect.TypeOf(img).String()

	if samples != want {
		t.Errorf("Expected %v got %v", want, samples)
	}
}

func TestNewImage_DefaultFormat(t *testing.T) {
	img := NewImage(10, 10, "")
	want := DefaultFormat
	samples := img.Format
	if samples != want {
		t.Errorf("Expected %v got %v", want, samples)
	}
}

func TestNewFromPath(t *testing.T) {
	want := "./samples/form-a.png"
	img, err := NewImageFromPath(want)
	if err != nil {
		t.Error(err)
	}
	if want != img.Path {
		t.Errorf("Expected %v got %v", want, img.Path)
	}
}

func TestNewFromPath_NotExists(t *testing.T) {
	path := "./samples/form-xxx.png"
	_, err := NewImageFromPath(path)
	if err == nil {
		t.Error(err)
	}
}

func TestImageSize(t *testing.T) {
	img, err := NewImageFromPath("./samples/gray8-a.png")
	if err != nil {
		t.Error(err)
	}
	want := 256
	samples := img.Size()
	if want != samples {
		t.Errorf("Expected %v got %v", want, samples)
	}
}

func TestImageLoad(t *testing.T) {
	path := "./samples/form-a.png"
	_, err := NewImageFromPath(path)
	if err != nil {
		t.Errorf("File %v decoded incorrectly", path)
	}

	path = "./samples/nonexists"
	_, err = NewImageFromPath(path)
	if err == nil {
		t.Errorf("File %v should not exist", path)
	}

	path = "./samples/corrupted.png"
	_, err = NewImageFromPath(path)
	if err == nil {
		t.Errorf("File %v cannot be decoded correctly", path)
	}

	path = "./samples/not-image"
	_, err = NewImageFromPath(path)
	if err != image.ErrFormat {
		t.Errorf("File %v is not an image", path)
	}
}

func TestImageEmpty(t *testing.T) {
	exp := true
	img := NewImage(0, 0, DefaultFormat)
	if !img.Empty() {
		t.Errorf("Expected %v got %v", exp, img.Empty())
	}

	exp = false
	img = NewImage(10, 10, DefaultFormat)
	if img.Empty() {
		t.Errorf("Expected %v got %v", exp, img.Empty())
	}
}

func TestDimensionsEqual(t *testing.T) {
	paths := []string{
		"./samples/bird-a.jpg",
		"./samples/bird-b.jpg",
		"./samples/bird-c-small.jpg",
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
		"./samples/form-a.png",
		"./samples/form-b.png",
		"./samples/form-a.png",
	}
	images := make([]*Image, 0, len(paths))
	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}

	want := false
	samples := images[0].Identical(images[1])
	if samples {
		t.Errorf("Expected %v got %v", want, samples)
	}

	want = true
	samples = images[0].Identical(images[2])
	if !samples {
		t.Errorf("Expected %v got %v", want, samples)
	}
}

func TestBytes(t *testing.T) {
	pairs := map[string]int{
		"./samples/models/nrgba.png":   16,
		"./samples/models/nrgba32.png": 32,
		"./samples/models/rgb.png":     16,
		"./samples/models/rgb32.png":   32,
		"./samples/models/gray.png":    4,
		"./samples/models/gray32.png":  8,
		"./samples/models/graya.png":   16,
		"./samples/models/graya32.png": 32,
		"./samples/models/palette.png": 4,
		"./samples/models/alpha.png":   16,
		"./samples/models/alpha32.png": 16,
		"./samples/models/tt.jpg":      64,
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
	img, err := NewImageFromPath("./samples/form-a.png")
	if err != nil {
		t.Error(nil)
	}
	point := image.Point{50, 50}
	samples := img.Position(point)
	want := 40200
	if samples != want {
		t.Errorf("Expected %v got %v", want, samples)
	}
}

func TestCompare_Empty(t *testing.T) {
	imgEmpty1 := NewImage(0, 0, DefaultFormat)
	imgEmpty2 := NewImage(0, 0, DefaultFormat)
	px, err := imgEmpty1.Compare(imgEmpty2, nil)
	if px > 0 && err != nil {
		t.Error(ErrImageIsEmpty.Error())
	}
}

func TestCompare_Dimensions(t *testing.T) {
	paths := []string{
		"./samples/bird-a.jpg",
		"./samples/bird-c-small.jpg",
	}
	images := make([]*Image, len(paths))
	for i, p := range paths {
		images[i], _ = NewImageFromPath(p)
	}
	px, err := images[0].Compare(images[1], nil)
	if px > 0 && err != nil {
		t.Error(ErrDimensionsDoNotMatch.Error())
	}
}

func TestCompare_Identical(t *testing.T) {
	paths := []string{
		"./samples/form-a.png",
		"./samples/form-a.png",
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
		"./samples/form-a.png",
		"./samples/form-b.png",
	}
	images := make([]*Image, len(paths))
	for i, p := range paths {
		images[i], _ = NewImageFromPath(p)
	}

	pairs := map[image.Point]float64{
		image.Point{35, 140}:  8097.090700072483,
		image.Point{50, 15}:   0,
		image.Point{105, 105}: -30.08278788879798,
	}

	for pt, want := range pairs {
		pos := images[0].Position(pt)
		samples := images[0].ColorDelta(images[1], pos, pos, false)
		if samples != want {
			t.Errorf("Expected %v got %v", want, samples)
		}
	}

	// Only Y (brigthness) component.
	pos := images[0].Position(image.Point{40, 220})
	samples := images[0].ColorDelta(images[1], pos, pos, true)
	want := 7.784791640000009
	if samples != want {
		t.Errorf("Expected %v got %v", want, samples)
	}
}

func TestSameNeighbors(t *testing.T) {
	img, _ := NewImageFromPath("./samples/form-a.png")
	pairs := map[image.Point]bool{
		image.Point{33, 218}: true,
		image.Point{107, 68}: false,
		image.Point{13, 72}:  true,
		image.Point{17, 72}:  true,
		image.Point{6, 70}:   true,
		image.Point{68, 119}: false,
	}

	for pt, want := range pairs {
		samples := img.SameNeighbors(pt, 3)
		if samples != want {
			t.Errorf("Expected %v got %v", want, samples)
		}
	}
}

func TestAntialiased(t *testing.T) {
	img, _ := NewImageFromPath("./samples/form-a.png")
	pairs := map[image.Point]bool{
		image.Point{0, 0}:    false,
		image.Point{50, 50}:  false,
		image.Point{42, 127}: true,
	}

	for pt, want := range pairs {
		samples := img.Antialiased(img, pt)
		if samples != want {
			t.Errorf("Expected %v got %v", want, samples)
		}
	}
}

// ---------------------------- PAIRS ------------------------------------

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
		pathA:        "./samples/gray8-a.png",
		pathB:        "./samples/gray8-b.png",
		pathDiff:     "./samples/gray8-diff.png",
		expectedDiff: 4,
		skip:         false,
		options:      NewOptions().SetIncludeAA(true),
	},
	testPair{
		name:         "PNG+GRAY16",
		pathA:        "./samples/gray16-a.png",
		pathB:        "./samples/gray16-b.png",
		pathDiff:     "./samples/gray16-diff.png",
		expectedDiff: 15,
		skip:         false,
		options:      NewOptions().SetIncludeAA(true),
	},
	testPair{
		name:         "PNG+Alpha",
		pathA:        "./samples/form-a.png",
		pathB:        "./samples/form-b.png",
		pathDiff:     "./samples/form-diff.png",
		expectedDiff: 2909,
		skip:         false,
		options:      NewOptions(),
	},
	testPair{
		name:         "PNG+Alpha+Anti-aliasing",
		pathA:        "./samples/form-a.png",
		pathB:        "./samples/form-b.png",
		pathDiff:     "./samples/form-aa-diff.png",
		expectedDiff: 3864,
		skip:         false,
		options:      NewOptions().SetIncludeAA(true),
	},
	testPair{
		name:         "GIF",
		pathA:        "./samples/landscape-a.gif",
		pathB:        "./samples/landscape-b.gif",
		pathDiff:     "./samples/landscape-diff.gif",
		expectedDiff: 9225,
		skip:         false,
		options:      NewOptions().SetAlpha(.5).SetThreshold(0.05).SetIncludeAA(true),
	},
	testPair{
		name:         "JPEG",
		pathA:        "./samples/bird-a.jpg",
		pathB:        "./samples/bird-b.jpg",
		pathDiff:     "./samples/bird-diff.jpg",
		expectedDiff: 1102,
		skip:         false,
		options:      NewOptions().SetAlpha(.5),
	},
	testPair{
		name:         "pixelmatch.js_100",
		pathA:        "./samples/original/1a.png",
		pathB:        "./samples/original/1b.png",
		pathDiff:     "./samples/original/1diff.png",
		expectedDiff: 143,
		skip:         false,
		options:      NewOptions().SetThreshold(.05),
	},
	testPair{
		name:         "pixelmatch.js_101",
		pathA:        "./samples/original/1a.png",
		pathB:        "./samples/original/1b.png",
		pathDiff:     "./samples/original/1diffmask.png",
		expectedDiff: 143,
		skip:         false,
		options:      NewOptions().SetThreshold(.05).SetDiffMask(true),
	},
	testPair{
		name:         "pixelmatch.js_102",
		pathA:        "./samples/original/1a.png",
		pathB:        "./samples/original/1b.png",
		pathDiff:     "./samples/original/1emptydiffmask.png",
		expectedDiff: 0,
		skip:         false,
		options:      NewOptions().SetThreshold(1).SetDiffMask(true),
	},
	testPair{
		name:         "pixelmatch.js_200",
		pathA:        "./samples/original/2a.png",
		pathB:        "./samples/original/2b.png",
		pathDiff:     "./samples/original/2diff.png",
		expectedDiff: 12437,
		skip:         false,
		options:      NewOptions().SetThreshold(.05).SetAlpha(.5).SetAAColor(color.RGBA{0, 0xc0, 0, 0xff}).SetDiffColor(color.RGBA{0xff, 0, 0xff, 0xff}),
	},
	testPair{
		name:         "pixelmatch.js_300",
		pathA:        "./samples/original/3a.png",
		pathB:        "./samples/original/3b.png",
		pathDiff:     "./samples/original/3diff.png",
		expectedDiff: 212,
		skip:         false,
		options:      NewOptions().SetThreshold(.05),
	},
	testPair{
		name:         "pixelmatch.js_400",
		pathA:        "./samples/original/4a.png",
		pathB:        "./samples/original/4b.png",
		pathDiff:     "./samples/original/4diff.png",
		expectedDiff: 36049,
		skip:         false,
		options:      NewOptions().SetThreshold(.05),
	},
	testPair{
		name:         "pixelmatch.js_500",
		pathA:        "./samples/original/5a.png",
		pathB:        "./samples/original/5b.png",
		pathDiff:     "./samples/original/5diff.png",
		expectedDiff: 0, // 0 because no AA included
		skip:         false,
		options:      NewOptions().SetThreshold(.05),
	},
	testPair{
		name:         "pixelmatch.js_600",
		pathA:        "./samples/original/6a.png",
		pathB:        "./samples/original/6b.png",
		pathDiff:     "./samples/original/6diff.png",
		expectedDiff: 51,
		skip:         false,
		options:      NewOptions().SetThreshold(.05),
	},
	testPair{
		name:         "pixelmatch.js_601",
		pathA:        "./samples/original/6a.png",
		pathB:        "./samples/original/6b.png",
		pathDiff:     "./samples/original/6empty.png",
		expectedDiff: 0,
		skip:         false,
		options:      NewOptions().SetThreshold(1.0).SetKeepEmptyDiff(true),
	},
	testPair{
		name:         "pixelmatch.js_700",
		pathA:        "./samples/original/7a.png",
		pathB:        "./samples/original/7b.png",
		pathDiff:     "./samples/original/7diff.png",
		expectedDiff: 2448,
		skip:         false,
		options:      NewOptions().SetDiffColorAlt(color.RGBA{0, 0xff, 0, 0xff}),
	},
}

func TestImageCompare(t *testing.T) {
	for _, pair := range testPairs {
		t.Run(pair.name, func(t *testing.T) {
			if pair.skip {
				t.SkipNow()
			}
			imageA, err := NewImageFromPath(pair.pathA)
			if err != nil {
				t.Error(err)
			}
			imageB, err := NewImageFromPath(pair.pathB)
			if err != nil {
				t.Error(err)
			}
			if updateDiffImages {
				fp, err := os.Create(pair.pathDiff)
				if err != nil {
					t.Error(err)
				}
				pair.options.SetOutput(fp)
			}
			diff, err := imageA.Compare(imageB, pair.options)
			if err != nil {
				t.Error(err.Error())
			}
			if diff != pair.expectedDiff {
				t.Errorf("Expected %v got %v", pair.expectedDiff, diff)
			}
		})
	}
}
