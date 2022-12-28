package pixmatch

import (
	"image/color"
)

// Options represents options structure where common comparison settings
// are stored.
type Options struct {
	// Threshold is alpha treshold of max. color delta.
	Threshold float64
	// Alpha is alpha channel factor.
	Alpha float64
	// DetectAA determinues if comparision should take in to account
	// anti-aliasing.
	DetectAA bool
	// AAColor is the color to mark antialiasing.
	AAColor color.Color
	// DiffColor is the color to mark the difference.
	DiffColor color.Color
	// DiffColorAlt is the alternative difference color.
	DiffColorAlt color.Color
	// DiffMask set to use mask.
	DiffMask bool
	// Output is the the final output image.
	Output *Output
}

var defaultOptions = Options{
	Threshold:    0.1,
	Alpha:        0.1,
	DetectAA:     false,
	AAColor:      color.RGBA{255, 255, 0, 255},
	DiffColor:    color.RGBA{255, 0, 0, 255},
	DiffColorAlt: nil,
	DiffMask:     false,
	Output:       nil,
}

// NewOptions creates new Options instance.
func NewOptions() *Options {
	return &defaultOptions
}

// ResolveDiffColor resolves the difference color or alternate difference
// color.
func (opts *Options) ResolveDiffColor() color.Color {
	if c := opts.DiffColorAlt; c != nil {
		return c
	}
	return opts.DiffColor
}
