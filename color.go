package pixmatch

// GNumber is a simplier number interface to work with colors.
// I really don't want to use golang.org/x/exp/constraints, because
// it is too general. Here is enough.
type GNumber interface {
	uint8 | uint16 | uint32 | float64
}

// Color represnts colors struct it's components (R)ed, (G)reen, (B)lue,
// (Alpha).
type Color[T GNumber] struct {
	R, G, B, A T
}

// NewColor create new color with certain type.
func NewColor[T GNumber](r, g, b, a T) *Color[T] {
	return &Color[T]{r, g, b, a}
}

// Equals checks color equality, that all color channels are equals.
// returns true if colors are equal.
func (c Color[T]) Equals(c2 *Color[T]) bool {
	return c.R == c2.R && c.G == c2.G && c.B == c2.B && c.A == c2.A
}

// RGBA implementation of color.Color interface but with generics.
// NOTE will it work?
func (c Color[T]) RGBA() (r, g, b, a T) {
	return c.R, c.G, c.B, c.A
}

// YIQ converts RGB to YIQ color space. See wiki page
// https://en.wikipedia.org/wiki/YIQ
func (c Color[T]) YIQ() (y, i, q float64) {
	return c.Y(), c.I(), c.Q()
}

// Blend is for blending colors with alpha.
func (c Color[T]) Blend(alpha T) *Color[T] {
	r := 255.0 + (c.R-255.0)*alpha
	g := 255.0 + (c.G-255.0)*alpha
	b := 255.0 + (c.B-255.0)*alpha
	a := alpha / 255.0
	return NewColor[T](r, g, b, a) // uinfere?
}

// Y is Y component of YIQ color space.
func (c Color[T]) Y() float64 {
	return float64(c.R)*float64(0.29889531) +
		float64(c.G)*float64(0.58662247) +
		float64(c.B)*float64(0.11448223)
}

// I is I component of YIQ color space.
func (c Color[T]) I() float64 {
	return float64(c.R)*float64(0.59597799) -
		float64(c.G)*float64(0.27417610) -
		float64(c.B)*float64(0.32180189)
}

// Q is Q component of YIQ color space.
func (c Color[T]) Q() float64 {
	return float64(c.R)*float64(0.21147017) -
		float64(c.G)*float64(0.52261711) +
		float64(c.B)*float64(0.31114694)
}

// func (c Color) String() string {
// 	return fmt.Sprintf("rgba(%s, %s, %s, %s)", c.R, c.G, c.B, c.A)
// }

// func (c Color) HexString() string {
// 	return fmt.Sprintf("%x%x%x%x", c.R, c.G, c.B, c.A)
// }
