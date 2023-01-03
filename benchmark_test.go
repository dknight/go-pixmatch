package pixmatch

import "testing"

var benchOpts = NewOptions()

func BenchmarkCompare_Empty(b *testing.B) {
	images := []*Image{
		NewImage(0, 0, DefaultFormat),
		NewImage(0, 0, DefaultFormat),
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		images[0].Compare(images[1], benchOpts)
	}
}

func BenchmarkCompare_Dimensions(b *testing.B) {
	paths := []string{
		"./samples/bird-c-small.jpg",
		"./samples/bird-b.jpg",
	}
	images := make([]*Image, 0, len(paths))
	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		images[0].Compare(images[1], benchOpts)
	}
}

func BenchmarkCompare_Identical(b *testing.B) {
	paths := []string{
		"./samples/form-b.png",
		"./samples/form-b.png",
	}
	images := make([]*Image, 0, len(paths))
	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		images[0].Compare(images[1], benchOpts)
	}
}

func BenchmarkCompare_DifferentPNG(b *testing.B) {
	paths := []string{
		"./samples/form-a.png",
		"./samples/form-b.png",
	}
	images := make([]*Image, 0, len(paths))
	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		images[0].Compare(images[1], benchOpts)
	}
}

func BenchmarkCompare_DifferentGIF(b *testing.B) {
	paths := []string{
		"./samples/landscape-a.gif",
		"./samples/landscape-b.gif",
	}
	images := make([]*Image, 0, len(paths))
	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		images[0].Compare(images[1], benchOpts)
	}
}

func BenchmarkCompare_DifferentJPEG(b *testing.B) {
	paths := []string{
		"./samples/bird-a.jpg",
		"./samples/bird-b.jpg",
	}
	images := make([]*Image, 0, len(paths))
	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		images[0].Compare(images[1], benchOpts)
	}
}

func BenchmarkCompare_DifferentAA(b *testing.B) {
	paths := []string{
		"./samples/form-a.png",
		"./samples/form-b.png",
	}
	images := make([]*Image, 0, len(paths))
	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}
	b.ResetTimer()

	benchOpts.SetIncludeAA(true)
	for i := 0; i < b.N; i++ {
		images[0].Compare(images[1], benchOpts)
	}
}
