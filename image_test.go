package pixmatch

import (
	"image"
	"reflect"
	"testing"
)

func TestNewImage(t *testing.T) {
	img := NewImage()
	want := "*pixmatch.Image"
	typ := reflect.TypeOf(img).String()

	if typ != want {
		t.Error("Expected type", want, "got", typ)
	}
}

func TestNewFromPath(t *testing.T) {
	want := "./res/kitten1.png"
	img, err := NewImageFromPath(want)
	if err != nil {
		t.Error(err)
	}
	if want != img.Path {
		t.Error("Expected", want, "got", img.Path)
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
	want := "./res/kitten1.png"
	img := NewImage()
	img.SetPath(want)
	if want != img.Path {
		t.Error("Expected", want, "got", img.Path)
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

	want, _ := images[0].DimensionsEqual(images[1])
	if !want {
		t.Error("Expected", want, "got", false)
	}

	want, err := images[0].DimensionsEqual(images[2])
	if want && err == ErrDimensionsDoNotMatch {
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

	want := false
	res := images[0].Identical(images[1])
	if res {
		t.Error("Expected", want, "got", res)
	}

	want = true
	res = images[0].Identical(images[2])
	if !res {
		t.Error("Expected", want, "got", res)
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
	point := image.Point{50, 50}
	res := img.Position(point)
	want := 20200
	if res != want {
		t.Errorf("Expecte %d got %d", want, res)
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

func TestImageColorDelta(t *testing.T) {
	paths := []string{
		"./res/kitten1.png",
		"./res/kitten2.png",
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

func TestSameColorNeighbors(t *testing.T) {
	img, _ := NewImageFromPath("./res/kitten1.png")
	pairs := map[image.Point]int{
		image.Point{0, 99}:  4,
		image.Point{7, 99}:  6,
		image.Point{13, 72}: 4,
		image.Point{17, 72}: 8,
		image.Point{17, 76}: 3,
		image.Point{3, 10}:  8, // alpha
		image.Point{11, 14}: 0,
	}

	for pt, want := range pairs {
		res := img.SameColorNeighbors(pt)
		if res != want {
			t.Errorf("Expected %v got %v", want, res)
		}
	}
}

func TestHasLeast3Neighbors(t *testing.T) {
	img, _ := NewImageFromPath("./res/kitten1.png")
	pairs := map[image.Point]bool{
		image.Point{0, 99}:  true,
		image.Point{3, 10}:  true,
		image.Point{11, 14}: false,
	}

	for pt, want := range pairs {
		res := img.HasLeast3Neighbors(pt)
		if res != want {
			t.Errorf("Expected %v got %v", want, res)
		}
	}
}

func TestAntialiased(t *testing.T) {
	img, _ := NewImageFromPath("./res/kitten1.png")
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

func TestFullCompare(t *testing.T) {
	paths := []string{
		"./res/kitten1.png",
		"./res/kitten2.png",
	}
	images := make([]*Image, len(paths))
	for i, p := range paths {
		images[i], _ = NewImageFromPath(p)
	}
	opts := DefaultOptions()
	output, err := NewOutput("diff.png",
		images[0].Bounds().Dx(), images[0].Bounds().Dy())
	if err != nil {
		t.Error(err)
	}
	opts.Output = output
	opts.DetectAA = true
	diff, err := images[0].Compare(images[1], opts)
	if err != nil {
		t.Error("Compare", err.Error())
	}
	want := 2245
	if want != diff {
		t.Errorf("Expected %v got %v", want, diff)
	}
}
