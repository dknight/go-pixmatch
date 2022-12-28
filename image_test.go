package pixmatch

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"reflect"
	"testing"
)

var removeDiffImages = true
var opts = NewOptions()

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
		"./res/models/tt.jpg":      64,
	}
	for path, bits := range pairs {
		img, _ := NewImageFromPath(path)
		bs := img.Bytes()
		if len(bs) != bits {
			t.Errorf("Expected %v (%v bits) got %v", path, bits, len(bs))
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

func TestImageColorDelta(t *testing.T) {
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
		image.Point{17, 76}: true,
		image.Point{3, 10}:  true, // alpha
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
		// image.Point{7, 88}: true, // FIXME
	}

	for pt, want := range pairs {
		res := img.Antialiased(img, pt)
		if res != want {
			t.Errorf("Expected %v got %v", want, res)
		}
	}
}

func TestFullCompare_PNG(t *testing.T) {
	// t.SkipNow()
	diffFileName := "diff-cat.png"
	t.Cleanup(func() {
		if removeDiffImages {
			os.Remove(diffFileName)
		}
	})
	paths := []string{
		"./res/kitten-a.png",
		"./res/kitten-b.png",
	}
	images := make([]*Image, len(paths))
	for i, p := range paths {
		images[i], _ = NewImageFromPath(p)
	}
	opts.DetectAA = false
	output, err := NewOutput(diffFileName,
		images[0].Width(), images[0].Height())
	if err != nil {
		t.Error(err)
	}

	// fmt.Printf("%+v\n", images[0].Image)

	opts.Output = output
	opts.DetectAA = false
	// opts.DiffMask = false
	diff, err := images[0].Compare(images[1], opts)
	if err != nil {
		t.Error("Compare", err.Error())
	}
	want := 2245
	if want != diff {
		t.Errorf("Expected %v got %v", want, diff)
	}
}

func TestFullCompare_PNGAA(t *testing.T) {
	// t.SkipNow()
	diffFileName := "diff-aa.png"
	t.Cleanup(func() {
		if removeDiffImages {
			os.Remove(diffFileName)
		}
	})
	paths := []string{
		"./res/aa-a.png",
		"./res/aa-b.png",
	}
	images := make([]*Image, len(paths))
	for i, p := range paths {
		images[i], _ = NewImageFromPath(p)
	}
	output, err := NewOutput(diffFileName,
		images[0].Width(), images[0].Height())
	if err != nil {
		t.Error(err)
	}

	opts.Output = output
	opts.DetectAA = true
	// opts.DiffMask = false
	diff, err := images[0].Compare(images[1], opts)
	if err != nil {
		t.Error("Compare", err.Error())
	}
	want := 51
	if want != diff {
		t.Errorf("Expected %v got %v", want, diff)
	}
}

func TestFullCompare_GIF(t *testing.T) {
	// t.SkipNow()
	diffFileName := "diff-abi.gif"
	t.Cleanup(func() {
		if removeDiffImages {
			os.Remove(diffFileName)
		}
	})
	paths := []string{
		"./res/abigail-a.gif",
		"./res/abigail-b.gif",
	}
	images := make([]*Image, len(paths))
	for i, p := range paths {
		images[i], _ = NewImageFromPath(p)
	}

	opts.DetectAA = false
	output, err := NewOutput(diffFileName,
		images[0].Width(), images[0].Height())
	if err != nil {
		t.Error(err)
	}

	opts.Output = output
	// opts.DetectAA = false
	// opts.DiffMask = true
	// opts.Threshold = .5
	diff, err := images[0].Compare(images[1], opts)
	if err != nil {
		t.Error("Compare", err.Error())
	}
	want := 294
	if want != diff {
		t.Errorf("Expected %v got %v", want, diff)
	}
}

func TestFullCompare_JPEG(t *testing.T) {
	// t.SkipNow()
	diffFileName := "diff-forest.jpg"
	t.Cleanup(func() {
		if removeDiffImages {
			os.Remove(diffFileName)
		}
	})
	paths := []string{
		"./res/forest-a.jpg",
		"./res/forest-b.jpg",
	}
	images := make([]*Image, len(paths))
	for i, p := range paths {
		images[i], _ = NewImageFromPath(p)
	}
	opts.DetectAA = false
	opts.Alpha = 1.0
	// opts.Threshold = .5
	output, err := NewOutput(diffFileName,
		images[0].Width(), images[0].Height())
	if err != nil {
		t.Error(err)
	}

	opts.Output = output
	opts.DetectAA = false
	diff, err := images[0].Compare(images[1], opts)
	if err != nil {
		t.Error("Compare", err.Error())
	}
	want := 1782
	if want != diff {
		t.Errorf("Expected %v got %v", want, diff)
	}
}
