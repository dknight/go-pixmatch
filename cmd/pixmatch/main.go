// Command Line Interface (CLI) for pixmatch
//
//	go install github.com/dknight/go-pixmatch/cmd/pixmatch@latest
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/dknight/go-pixmatch"
)

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
var keepUsage = "Keep empty output files. Valid only with -o flag."
var nUsage = "Do not output the trailing newline."
var watchUsage = "Experimental: Watch for input pair of images"

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
var keep bool
var n bool
var watch bool

var fpOutout *os.File

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
	flag.BoolVar(&keep, "keep", false, keepUsage)
	flag.BoolVar(&n, "n", false, nUsage)
	flag.BoolVar(&watch, "w", false, watchUsage)
	flag.Parse()

	// Just display version.
	if version {
		fmt.Println(pixmatch.GetVersion())
		os.Exit(pixmatch.ExitOk)
	}
	argsCount := flag.NArg()
	if argsCount == 0 && !watch {
		flag.Usage()
		os.Exit(pixmatch.ExitOk)
	}
	if argsCount < 2 && !watch {
		exitErr(pixmatch.ExitMissingImage, pixmatch.ErrMissingImage)
	}
}

func runComparison(paths []string) (int, int) {
	opts := pixmatch.NewOptions()
	images := make([]*pixmatch.Image, 2)

	setupOptions(opts)

	// Load images
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

	// Compare images
	px, err := images[0].Compare(images[1], opts)
	if err != nil {
		switch err {
		case pixmatch.ErrDimensionsDoNotMatch:
			exitErr(pixmatch.ExitDimensionsNotEqual, err)
		case pixmatch.ErrImageIsEmpty:
			exitErr(pixmatch.ExitEmptyImage, err)
		case pixmatch.ErrUnknownFormat:
			exitErr(pixmatch.ExitUnknownFormat, err)
		default:
			exitErr(pixmatch.ExitUnknown, err)
		}
	}

	return px, images[0].Size()
}

func main() {
	paths := make([]string, 2)
	args := flag.Args()
	var px int
	var size int
	// TODO deal with watch()
	if watch {
		homeDir, _ := os.UserHomeDir()
		fname := fmt.Sprintf("%s/%s", homeDir, "pixmatch.stream")
		file, err := os.OpenFile(fname, os.O_CREATE|os.O_RDONLY, 0644)
		if err != nil {
			log.Fatal("Open named file error:", err)
		}
		rd := bufio.NewReader(file)
		for {
			line, _ := rd.ReadString('\n')
			// if err != nil {
			// 	log.Fatal(err)
			// }
			paths = strings.Fields(line)
			if len(paths) > 0 {
				px, size = runComparison(paths)
			}
		}
	}

	for i, arg := range args {
		paths[i] = arg
	}
	px, size = runComparison(paths)

	// If no diference remove file.
	if output != "" && px <= 0 && !keep {
		os.Remove(output)
	}

	output := format(px, percent, size)
	fmt.Fprint(os.Stdout, output)

	fpOutout.Close()
}

func format(d int, isPct bool, size int) string {
	if isPct {
		format := "%.2f%%"
		if !n {
			format += "\n"
		}
		pct := float64(d) / float64(size) * 100
		return fmt.Sprintf(format, pct)
	}
	format := "%d"
	if !n {
		format += "\n"
	}
	return fmt.Sprintf(format, d)
}

func setupOptions(opts *pixmatch.Options) {
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
}

func exitErr(status int, errs ...error) {
	for _, e := range errs {
		if e.Error() != "" {
			fmt.Fprintln(os.Stderr, e.Error())
		}
	}
	os.Exit(status)
}
