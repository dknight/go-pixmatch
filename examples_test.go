package pixmatch

import (
	"fmt"
	"image/color"
	"log"
)

func Example() {
	img1, err := NewImageFromPath("./samples/form-a.png")
	if err != nil {
		log.Fatalln(err)
	}
	img2, err := NewImageFromPath("./samples/form-b.png")
	if err != nil {
		log.Fatalln(err)
	}
	options := NewOptions()
	options.SetThreshold(0.05)
	options.SetAlpha(0.5)
	options.SetDiffColor(color.RGBA{0, 255, 127, 255})
	// etc...

	diff, err := img1.Compare(img2, options)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(diff)
	// Output: 4933
}

func ExampleImage_Compare() {
	options := NewOptions()
	options.SetThreshold(0.25)

	img1, _ := NewImageFromPath("./samples/form-a.png")
	img2, _ := NewImageFromPath("./samples/form-b.png")
	diff, _ := img1.Compare(img2, options)

	fmt.Println(diff)
	// Output: 1626
}

func ExampleImage_DimensionsEqual() {
	img1, _ := NewImageFromPath("./samples/bird-a.jpg")
	img2, _ := NewImageFromPath("./samples/bird-c-small.jpg")
	samplesult := img1.DimensionsEqual(img2)
	fmt.Println(samplesult)
	// Output: false
}

func ExampleImage_Empty() {
	img, _ := NewImageFromPath("./samples/form-a.png")
	samplesult := img.Empty()
	fmt.Println(samplesult)
	// Output: false
}

func ExampleImage_Identical() {
	img1, _ := NewImageFromPath("./samples/form-a.png")
	img2, _ := NewImageFromPath("./samples/form-a.png")
	samplesult := img1.Identical(img2)
	fmt.Println(samplesult)
	// Output: true
}

func ExampleImage_Size() {
	img, _ := NewImageFromPath("./samples/form-a.png")
	samplesult := img.Size()
	fmt.Println(samplesult)
	// Output: 51200
}
