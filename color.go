package pixmatch

import (
	"fmt"
	"image/color"
)

// Color represents color structure and it's components (R)ed, (G)reen,
// (B)lue, (Alpha). This is similar to color.Color from standard library.
type Color struct {
	R, G, B, A uint32
}

// NewColor creates new color instance.
func NewColor(r, g, b, a uint32) *Color {
	return &Color{r, g, b, a}
}

// Equals checks colors' equality, that all color's channels are equal.
func (c Color) Equals(c2 *Color) bool {
	return c.R == c2.R && c.G == c2.G && c.B == c2.B && c.A == c2.A
}

// RGBA returns Red, Green, Blue, Alpha channels similar to color.RGBA()
// from standard library, but always returns values as uint32 type.
func (c Color) RGBA() (r, g, b, a uint32) {
	return c.R, c.G, c.B, c.A
}

// YIQ converts RGB to YIQ color space. See wiki page about YIQ
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

// Blend is the procedure of blending the color with the alpha factor is
// known as blending.
func (c Color) Blend(a float64) *Color {
	r := 255 + float64(c.R-255)*(a)
	g := 255 + float64(c.G-255)*(a)
	b := 255 + float64(c.B-255)*(a)
	return NewColor(uint32(r), uint32(g), uint32(b), c.A)
}

// BlendToGray draws gray-scaled color with gray-scaled blending.
func (c Color) BlendToGray(a float64) color.Color {
	y := uint32(c.Y()) >> 8
	gray := uint8(255 + (float64(y)-255)*a)
	if c.A == 0 {
		gray = 255
	}
	return color.RGBA{gray, gray, gray, 255}
}

func (c Color) String() string {
	return fmt.Sprintf("(%v,%v,%v,%v)", c.R, c.G, c.B, c.A)
}
