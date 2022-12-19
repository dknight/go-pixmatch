package pixmatch

import "image/color"

// Options represents options structure where common comparison settings
// are stored.
type Options struct {
	Threshold    float64
	IncludeAA    bool
	Alpha        float64
	AAColor      color.Color
	DiffColor    color.Color
	DiffColorAlt color.Color
	DiffMask     bool
}

// NewOptions creates new options.
func NewOptions() *Options {
	return &defaultOptions
}

var defaultOptions = Options{
	Threshold:    0.1,
	IncludeAA:    false,
	Alpha:        0.1,
	AAColor:      color.RGBA{255, 255, 0, 0},
	DiffColor:    color.RGBA{255, 0, 0, 0},
	DiffColorAlt: nil,
	DiffMask:     false,
}
