package pixmatch

import (
	"fmt"
	"image/color"
	"math"
)

// Color represents colors struct it's components (R)ed, (G)reen, (B)lue,
// (Alpha).
type Color struct {
	R, G, B, A uint32
}

// NewColor create new color with certain type.
func NewColor(r, g, b, a uint32) *Color {
	return &Color{r, g, b, a}
}

// Equals checks color equality, that all color channels are equals.
// returns true if colors are equal.
func (c Color) Equals(c2 *Color) bool {
	return c.R == c2.R && c.G == c2.G && c.B == c2.B && c.A == c2.A
}

// RGBA implementation of color.Color interface but with generics.
func (c Color) RGBA() (r, g, b, a uint32) {
	return c.R, c.G, c.B, c.A
}

// YIQ converts RGB to YIQ color space. See wiki page about YIQ:
// https://en.wikipedia.org/wiki/YIQ
func (c Color) YIQ() (float64, float64, float64) {
	return c.Y(), c.I(), c.Q()
}

// Y is RBG to Y (brightness) conversion.
func (c Color) Y() float64 {
	return float64(c.R)*0.29889531 +
		float64(c.G)*0.58662247 +
		float64(c.B)*0.11448223
}

// I is RBG to I (chrominance) conversion.
func (c Color) I() float64 {
	return float64(c.R)*0.59597799 -
		float64(c.G)*0.27417610 -
		float64(c.B)*0.32180189
}

// Q is RBG to Q (chrominance) conversion.
func (c Color) Q() float64 {
	return float64(c.R)*0.21147017 -
		float64(c.G)*0.52261711 +
		float64(c.B)*0.31114694
}

// Blend is for blending colors with alpha.
func (c Color) Blend(alpha uint32) *Color {
	r := 255 + (c.R-255)*alpha
	g := 255 + (c.G-255)*alpha
	b := 255 + (c.B-255)*alpha
	return NewColor(r, g, b, alpha)
}

// BlendToGray draws greyscaled color multiplied by alpha factor.
func (c Color) BlendToGray(alpha float64) color.Color {
	gray := uint8(math.Round(255.0 - alpha*(c.Y()-255.0)/255.0))
	return color.RGBA{gray, gray, gray, 255}
}

func (c Color) String() string {
	return fmt.Sprintf("(%v,%v,%v,%v)", c.R, c.G, c.B, c.A)
}
