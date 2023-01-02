package pixmatch

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"strings"
)

// Color represents color structure and its components (R)ed, (G)reen,
// (B)lue, (Alpha). This is similar to color.Color from the standard library.
type Color struct {
	R, G, B, A uint32
}

// NewColor creates a new color instance.
func NewColor(r, g, b, a uint32) *Color {
	return &Color{r, g, b, a}
}

// Equals checks colors' equality, ensuring that all color channels are equal.
func (c Color) Equals(c2 *Color) bool {
	return c.R == c2.R && c.G == c2.G && c.B == c2.B && c.A == c2.A
}

// RGBA returns Red, Green, Blue, Alpha channels similar to color.RGBA()
// from the standard library, which always returns values as uint32 type.
func (c Color) RGBA() (r, g, b, a uint32) {
	return c.R, c.G, c.B, c.A
}

// YIQ converts RGB intoto YIQ color space. See the wiki page about YIQ:
// https://en.wikipedia.org/wiki/YIQ
func (c Color) YIQ() (float64, float64, float64) {
	return c.Y(), c.I(), c.Q()
}

// Y is the RBG to Y (brightness) conversion.
func (c Color) Y() float64 {
	return float64(c.R)*0.29889531 +
		float64(c.G)*0.58662247 +
		float64(c.B)*0.11448223
}

// I is the RBG to I (chrominance) conversion.
func (c Color) I() float64 {
	return float64(c.R)*0.59597799 -
		float64(c.G)*0.27417610 -
		float64(c.B)*0.32180189
}

// Q is the RBG to Q (chrominance) conversion.
func (c Color) Q() float64 {
	return float64(c.R)*0.21147017 -
		float64(c.G)*0.52261711 +
		float64(c.B)*0.31114694
}

// Blend is the procedure of blending the color with the alpha factor is
// known as blending.
func (c Color) Blend(a float64) *Color {
	r := 255 - float64(c.R)*a
	g := 255 - float64(c.G)*a
	b := 255 - float64(c.B)*a
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

// HexStringToColor converts hexadecimal string RRGGBBAA of color
// representation to color.RGBA. Input string are case-insensitive.
// Also strings can be prefixed with '0x' or '0X'.
//
// Examples values are:
//   - FF000099
//   - ff00ff00
//   - 0xff00ff00
//   - #ffFF00ff00
func HexStringToColor(hexstr string) (*color.RGBA, error) {
	s := strings.ToUpper(hexstr)
	s = strings.TrimPrefix(s, "0X")
	if len(s) != 8 {
		return nil, ErrInvalidColorFormat
	}
	bs, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return &color.RGBA{bs[0], bs[1], bs[2], bs[3]}, nil
}

func (c Color) String() string {
	return fmt.Sprintf("(%v,%v,%v,%v)", c.R, c.G, c.B, c.A)
}
