package pixmatch

import (
	"image/color"
)

// Options represents the structure where common comparison's settings
// are stored.
type Options struct {
	// Threshold is threshold of the maximum color delta.
	Threshold float64

	// Alpha is alpha channel factor (multiplier). Allowed values 0..1..
	Alpha float64

	// IncludeAA sets anti-aliasing pixels as difference count.
	IncludeAA bool

	// AAColor is the color to mark anti-aliasing pixels.
	AAColor color.Color

	// DiffColor is the color to highlight the differences.
	DiffColor color.Color

	// DiffColorAlt is the alternative difference color. Used whether
	// to detect dark on light differences between two images and set
	// an alternative color if required.
	DiffColorAlt color.Color

	// DiffMask set to use mask, renders the differences without original
	// image.
	DiffMask bool

	// Output is the final output of the image.
	Output *Output
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

// NewOptions creates new Options instance. Here possible to use
// https://github.com/imdario/mergo
// Me perosnally always try to avoid dependencies where possible.
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

// SetThreshold sets threshold to options.
func (opts *Options) SetThreshold(v float64) *Options {
	opts.Threshold = v
	return opts
}

// SetAlpha sets alpha to options.
func (opts *Options) SetAlpha(v float64) *Options {
	opts.Alpha = v
	return opts
}

// SetIncludeAA sets anti-aliasing to options to counts anti-aliased
// pixels as difference.
func (opts *Options) SetIncludeAA(v bool) *Options {
	opts.IncludeAA = v
	return opts
}

// SetAAColor sets anti-aliased color to options.
func (opts *Options) SetAAColor(v color.Color) *Options {
	opts.AAColor = v
	return opts
}

// SetDiffColor sets color of difference.
func (opts *Options) SetDiffColor(v color.Color) *Options {
	opts.DiffColor = v
	return opts
}

// SetDiffColorAlt sets color of alternative difference.
func (opts *Options) SetDiffColorAlt(v color.Color) *Options {
	opts.DiffColorAlt = v
	return opts
}

// SetDiffMask sets difference mask to options.
func (opts *Options) SetDiffMask(v bool) *Options {
	opts.DiffMask = v
	return opts
}

// SetOutput sets the output as pointer to options.
func (opts *Options) SetOutput(v *Output) *Options {
	opts.Output = v
	return opts
}
