package pixmatch

import (
	"image/color"
	"io"
)

// Options is the structure that stores the settings for common comparisons.
type Options struct {
	// Threshold is the threshold of the maximum color delta.
	// Values range [0, 1.0].
	Threshold float64

	// Alpha is the alpha channel factor (multiplier). Values range [0, 1.0].
	Alpha float64

	// IncludeAA sets anti-aliasing pixels as difference counts.
	IncludeAA bool

	// AAColor is the color to mark anti-aliasing pixels.
	AAColor color.Color

	// DiffColor is the color to highlight the differences.
	DiffColor color.Color

	// DiffColorAlt is the alternative difference color. Used to detect dark
	// and light differences between two images and set an alternative color if
	// required.
	DiffColorAlt color.Color

	// DiffMask sets to use mask, renders the differences without the original
	// image.
	DiffMask bool

	// Output is stucture where final image will be written.
	Output io.Writer
}

// defaultOptions are just default options.
var defaultOptions = Options{
	Threshold:    0.1,
	Alpha:        0.1,
	IncludeAA:    false,
	AAColor:      color.RGBA{255, 255, 0, 255},
	DiffColor:    color.RGBA{255, 0, 0, 255},
	DiffColorAlt: nil,
	DiffMask:     false,
	Output:       nil,
}

// NewOptions creates a new Options instance. It is possible to use
// https://github.com/imdario/mergo in this case. Personally, I try to avoid
// dependencies whenever possible.
func NewOptions() *Options {
	return &Options{
		Threshold:    defaultOptions.Threshold,
		Alpha:        defaultOptions.Alpha,
		IncludeAA:    defaultOptions.IncludeAA,
		AAColor:      defaultOptions.AAColor,
		DiffColor:    defaultOptions.DiffColor,
		DiffColorAlt: defaultOptions.DiffColorAlt,
		DiffMask:     defaultOptions.DiffMask,
		Output:       defaultOptions.Output,
	}
}

// SetThreshold sets threshold to the options.
func (opts *Options) SetThreshold(v float64) *Options {
	opts.Threshold = v
	return opts
}

// SetAlpha sets alpha to the options.
func (opts *Options) SetAlpha(v float64) *Options {
	opts.Alpha = v
	return opts
}

// SetIncludeAA sets anti-aliasing to the options to counts anti-aliased
// pixels as differences.
func (opts *Options) SetIncludeAA(v bool) *Options {
	opts.IncludeAA = v
	return opts
}

// SetAAColor sets anti-aliased color to the options.
func (opts *Options) SetAAColor(v color.Color) *Options {
	opts.AAColor = v
	return opts
}

// SetDiffColor sets color of differences to the options.
func (opts *Options) SetDiffColor(v color.Color) *Options {
	opts.DiffColor = v
	return opts
}

// SetDiffColorAlt sets color of alternative difference to the options.
func (opts *Options) SetDiffColorAlt(v color.Color) *Options {
	opts.DiffColorAlt = v
	return opts
}

// SetDiffMask sets difference mask to the options.
func (opts *Options) SetDiffMask(v bool) *Options {
	opts.DiffMask = v
	return opts
}

// SetOutput sets the output as pointer to the options.
func (opts *Options) SetOutput(v io.Writer) *Options {
	opts.Output = v
	return opts
}
