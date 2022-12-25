package pixmatch

import "testing"

func BenchmarkIdentical(b *testing.B) {
	paths := []string{
		"./res/kitten-b.png",
		"./res/kitten-b.png",
	}
	images := make([]*Image, 0, len(paths))
	for _, p := range paths {
		img, _ := NewImageFromPath(p)
		images = append(images, img)
	}

	for i := 0; i < b.N; i++ {
		images[0].Identical(images[1])
	}
}