package pixmatch

import "testing"

func BenchmarkIdentical(b *testing.B) {
	paths := []string{"./res/kitten2.png", "./res/kitten2.png"}
	images := make([]*Image, len(paths))
	for i, p := range paths {
		images[i] = NewImage()
		images[i].SetPath(p)
		_ = images[i].Load()
	}

	for i := 0; i < b.N; i++ {
		images[0].Identical(images[1])
	}
}
