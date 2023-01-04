package pixmatch

import (
	"image/color"
	"io"
)

// Options is the structure that stores the settings for common comparisons.
type Options struct {
	// Output is stucture where final image will be written.
	Output io.Writer

	// Threshold is the threshold of the maximum color delta.
	// Values range [0, 1.0].
	Threshold float64

	// Alpha is the alpha channel factor (multiplier). Values range [0, 1.0].
	// NOTE it is interesting to experiment with overflow and underflow
	// ranges.
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

	// KeepEmptyDiff removes empty diff files.
	KeepEmptyDiff bool
}

// defaultOptions are just default options.
var defaultOptions = Options{
	Output:        nil,
	Threshold:     0.1,
	Alpha:         0.1,
	IncludeAA:     false,
	AAColor:       color.RGBA{0xff, 0xff, 0, 0xff},
	DiffColor:     color.RGBA{0xff, 0, 0, 0xff},
	DiffColorAlt:  nil,
	DiffMask:      false,
	KeepEmptyDiff: false,
}

// NewOptions creates a new Options instance. It is possible to use
// https://github.com/imdario/mergo in this case. Personally, I try to avoid
// dependencies whenever possible.
func NewOptions() *Options {
	return &Options{
		Output:        defaultOptions.Output,
		Threshold:     defaultOptions.Threshold,
		Alpha:         defaultOptions.Alpha,
		IncludeAA:     defaultOptions.IncludeAA,
		AAColor:       defaultOptions.AAColor,
		DiffColor:     defaultOptions.DiffColor,
		DiffColorAlt:  defaultOptions.DiffColorAlt,
		DiffMask:      defaultOptions.DiffMask,
		KeepEmptyDiff: defaultOptions.KeepEmptyDiff,
	}
}

// SetOutput sets the output as pointer to the options.
func (opts *Options) SetOutput(v io.Writer) *Options {
	opts.Output = v
	return opts
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

// SetKeepEmptyDiff sets difference mask to the options.
func (opts *Options) SetKeepEmptyDiff(v bool) *Options {
	opts.KeepEmptyDiff = v
	return opts
}
