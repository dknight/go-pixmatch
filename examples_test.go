package pixmatch

import (
	"fmt"
	"image/color"
	"log"
)

func Example() {
	img1, err := NewImageFromPath("./res/kitten-a.png")
	if err != nil {
		log.Fatalln(err)
	}
	img2, err := NewImageFromPath("./res/kitten-b.png")
	if err != nil {
		log.Fatalln(err)
	}
	options := NewOptions()
	options.SetThreshold(0.05)
	options.SetAlpha(0.5)
	options.SetDiffColor(color.RGBA{0, 255, 128, 255})
	// etc...

	diff, err := img1.Compare(img2, options)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(diff)
	// Output: 1723
}

func ExampleImage_Compare() {
	options := NewOptions()
	options.SetThreshold(0.25)

	img1, _ := NewImageFromPath("./res/kitten-a.png")
	img2, _ := NewImageFromPath("./res/kitten-b.png")
	diff, _ := img1.Compare(img2, options)

	fmt.Println(diff)
	// Output: 1563
}

func ExampleImage_DimensionsEqual() {
	img1, _ := NewImageFromPath("./res/kitten-a.png")
	img2, _ := NewImageFromPath("./res/kitten-c-small.png")
	result := img1.DimensionsEqual(img2)
	fmt.Println(result)
	// Output: false
}

func ExampleImage_Empty() {
	img, _ := NewImageFromPath("./res/kitten-a.png")
	result := img.Empty()
	fmt.Println(result)
	// Output: false
}

func ExampleImage_Identical() {
	img1, _ := NewImageFromPath("./res/kitten-a.png")
	img2, _ := NewImageFromPath("./res/kitten-a.png")
	result := img1.Identical(img2)
	fmt.Println(result)
	// Output: true
}

func ExampleImage_Size() {
	img, _ := NewImageFromPath("./res/kitten-a.png")
	result := img.Size()
	fmt.Println(result)
	// Output: 10000
}
