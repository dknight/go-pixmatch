package main

import (
	"fmt"
	"os"

	"github.com/dknight/go-pixmatch"
)

func main() {
	opts := pixmatch.NewOptions()

	f1 := "../res/kitten1.png"
	f2 := "../res/kitten2.png"

	img1 := pixmatch.NewImage()
	img1.SetPath(f1)
	img2 := pixmatch.NewImage()
	img2.SetPath(f2)

	images := [pixmatch.ImagesCount]*pixmatch.Image{img1, img2}
	err := pixmatch.LoadImages(images)
	if err != nil {
		exitErr(pixmatch.ExitFSFail, err)
	}

	px, err := img1.Compare(img2, opts)
	if err != nil {
		exitErr(pixmatch.ExitEmptyImage, err)
	}
	fmt.Println(px)
}

func exitErr(status int, errs ...error) {
	for _, e := range errs {
		if e.Error() != "" {
			fmt.Fprintln(os.Stderr, e.Error())
		}
	}
	os.Exit(status)
}
