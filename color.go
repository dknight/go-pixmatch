package pixmatch

type Color struct {
	R, G, B, A uint32
}

func NewColor(r, g, b, a uint32) *Color {
	return &Color{r, g, b, a}
}

func (c Color) RGBA() (r, g, b, a uint32) {
	return c.R, c.G, c.B, c.A
}

// func (c Color) Bytes() []byte {
// 	return []byte{c.R, c.G, c.B, c.A}
// }

// func (c Color) String() string {
// 	return fmt.Sprintf("rgba(%s, %s, %s, %s)", c.R, c.G, c.B, c.A)
// }

// func (c Color) HexString() string {
// 	return fmt.Sprintf("%x%x%x%x", c.R, c.G, c.B, c.A)
// }
