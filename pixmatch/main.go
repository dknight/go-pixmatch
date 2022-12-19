package main

import (
	"fmt"
	"os"

	"github.com/dknight/go-pixmatch"
)

func main() {
	f1 := "../res/kitten1.png"
	// f2 := "./kitten2.png"
	f2 := "../res/kitten-small.png"

	img1 := pixmatch.NewImage()
	img1.SetPath(f1)
	img2 := pixmatch.NewImage()
	img2.SetPath(f2)

	images := [pixmatch.ImagesCount]*pixmatch.Image{img1, img2}
	err := pixmatch.LoadImages(images)
	if err != nil {
		exitErr(pixmatch.ExitFSFail, err)
	}

	fmt.Println(images[0].DimensionsEqual(images[1]))
	fmt.Printf("%+v\n", images[0])
	fmt.Printf("%+v\n", images[1])
}

func exitErr(status int, errs ...error) {
	for _, e := range errs {
		if e.Error() != "" {
			fmt.Fprintln(os.Stderr, e.Error())
		}
	}
	os.Exit(status)
}
