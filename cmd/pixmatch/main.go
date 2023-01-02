// Command Line Interface (CLI) for pixmatch
//
//	go install github.com/dknight/go-pixmatch/cmd/pixmatch@latest
package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/dknight/go-pixmatch"
)

// Usages strings
var outputUsage = "Output file path."
var thresholdUsage = "Threshold of the maximum color delta." +
	" Values range [0..1] (default 0.1)."
var alphaUsage = "Alpha channel factor. Values range [0..1]. (default 0.1)"
var aaUsage = "Count anti-aliasing pixels as difference (default false)."
var aaColorUsage = "Color to mark anti-aliasing pixels. Works only without" +
	" -aa flag (default ffff00ff)"
var diffColorUsage = "Color to highlight the differences (default ff0000ff)"
var diffColorAltUsage = "Alternative difference color. Used to detect dark" +
	" and light differences between two images and set an alternative color" +
	" if required (default nil)."
var maskUsage = "mask renders the differences without the" +
	" original image (default false)."
var versionUsage = "Display the version of pixmatch."
var percentUsage = "Display the difference in percent, instead of pixels" +
	" (default false)."

// Initialize flags
var output string
var threshold float64
var alpha float64
var aa bool
var aaColor string
var diffColor string
var diffColorAlt string
var mask bool
var version bool
var percent bool

func init() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(out)
		fmt.Fprintln(out, "pixelmatch [flags] image1.png image2.png")
		fmt.Fprintln(out)
		fmt.Fprintln(out, "Colors are in hexadecimal format: (0x)RRGGBBAA.")
		fmt.Fprintln(out, "Examples:")
		fmt.Fprintln(out, "\t- ff00ffff")
		fmt.Fprintln(out, "\t- FF00CCFF")
		fmt.Fprintln(out, "\t- 0x0033ffff")
		fmt.Fprintln(out, "\t- 0XAAFF00FF")
		fmt.Fprintln(out, "Flags:")
		flag.PrintDefaults()
	}
	flag.StringVar(&output, "o", "", outputUsage)
	flag.Float64Var(&threshold, "t", 0, thresholdUsage)
	flag.Float64Var(&alpha, "a", 0, alphaUsage)
	flag.BoolVar(&aa, "aa", false, aaUsage)
	flag.StringVar(&aaColor, "aacolor", "", aaColorUsage)
	flag.StringVar(&diffColor, "diffcolor", "", diffColorUsage)
	flag.StringVar(&diffColorAlt, "diffcoloralt", "", diffColorAltUsage)
	flag.BoolVar(&mask, "mask", false, maskUsage)
	flag.BoolVar(&version, "v", false, versionUsage)
	flag.BoolVar(&percent, "percent", false, percentUsage)
	flag.Parse()

	// Just display version.
	if version {
		fmt.Println(pixmatch.GetVersion())
		os.Exit(pixmatch.ExitOk)
	}
	argsCount := flag.NArg()
	if argsCount == 0 {
		flag.Usage()
		os.Exit(pixmatch.ExitOk)
	}
	if argsCount < 2 {
		exitErr(pixmatch.ExitMissingImage, pixmatch.ErrMissingImage)
	}
}

func main() {
	opts := pixmatch.NewOptions()
	paths := make([]string, 2)
	images := make([]*pixmatch.Image, 2)

	// Set options -------------------------------------------
	if output != "" {
		fp, err := os.Create(output)
		if err != nil {
			exitErr(pixmatch.ExitFSFail, err)
		}
		opts.SetOutput(fp)
	}
	if threshold != 0 {
		opts.SetThreshold(threshold)
	}
	if alpha != 0 {
		opts.SetAlpha(alpha)
	}
	if aa {
		opts.SetIncludeAA(true)
	}
	if aaColor != "" {
		color, err := pixmatch.HexStringToColor(aaColor)
		if err != nil {
			exitErr(pixmatch.ExitInvalidInput, err)
		}
		opts.SetAAColor(color)
	}
	if diffColor != "" {
		color, err := pixmatch.HexStringToColor(diffColor)
		if err != nil {
			exitErr(pixmatch.ExitInvalidInput, err)
		}
		opts.SetDiffColor(color)
	}
	if diffColorAlt != "" {
		color, err := pixmatch.HexStringToColor(diffColorAlt)
		if err != nil {
			exitErr(pixmatch.ExitInvalidInput, err)
		}
		opts.SetDiffColor(color)
	}
	if mask {
		opts.SetDiffMask(true)
	}
	// -------------------------------------------------------

	args := flag.Args()
	for i, arg := range args {
		paths[i] = arg
	}

	var wg sync.WaitGroup
	for i := range paths {
		wg.Add(1)
		go func(i int) {
			defer func() {
				if r := recover(); r != nil {
					exitErr(pixmatch.ExitFSFail, r.(error))
				}
			}()
			img, err := pixmatch.NewImageFromPath(paths[i])
			if err != nil {
				panic(err)
			}
			images[i] = img
			wg.Done()
		}(i)
	}
	wg.Wait()

	px, err := images[0].Compare(images[1], opts)
	if err != nil {
		exitErr(pixmatch.ExitEmptyImage, err)
	}

	if percent {
		pct := float64(px) / float64(images[0].Size()) * 100
		fmt.Fprintf(os.Stdout, "%.2f%%", pct)
	} else {
		fmt.Fprintf(os.Stdout, "%d", px)
	}
}

func exitErr(status int, errs ...error) {
	for _, e := range errs {
		if e.Error() != "" {
			fmt.Fprintln(os.Stderr, e.Error())
		}
	}
	os.Exit(status)
}
