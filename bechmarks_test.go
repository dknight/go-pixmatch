package pixmatch

import "testing"

var benchOpts = NewOptions()

func BenchmarkCompare_Empty(b *testing.B) {
	images := []*Image{
		NewImage(),
		NewImage(),
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		images[0].Compare(images[1], benchOpts)
	}
}

func BenchmarkCompare_Dimensions(b *testing.B) {
	paths := []string{
		"./res/kitten-c-small.png",
		"./res/kitten-b.png",
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
		"./res/kitten-b.png",
		"./res/kitten-b.png",
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

func BenchmarkCompare_Different(b *testing.B) {
	paths := []string{
		"./res/kitten-a.png",
		"./res/kitten-b.png",
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
		"./res/aa-a.png",
		"./res/aa-b.png",
	}
	images := make([]*Image, 0, len(paths))
	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}
	b.ResetTimer()

	benchOpts.DetectAA = true
	for i := 0; i < b.N; i++ {
		images[0].Compare(images[1], benchOpts)
	}
}
