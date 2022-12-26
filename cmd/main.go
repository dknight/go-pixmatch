package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/dknight/go-pixmatch"
)

func main() {
	paths := []string{
		"../res/kitten1.png",
		"../res/kitten2.png",
	}

	opts := pixmatch.NewOptions()
	images := make([]*pixmatch.Image, 0, len(paths))

	var wg sync.WaitGroup
	wg.Add(len(paths))

	for _, path := range paths {
		go func(p string) {
			defer func() {
				if r := recover(); r != nil {
					exitErr(pixmatch.ExitFSFail, r.(error))
				}
			}()
			img, err := pixmatch.NewImageFromPath(p)
			if err != nil {
				panic(err)
			}
			images = append(images, img)
			wg.Done()
		}(path)
	}

	wg.Wait()

	px, err := images[0].Compare(images[1], opts)
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
