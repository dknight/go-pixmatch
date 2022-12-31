package pixmatch

import (
	"reflect"
	"testing"
)

func TestNewColor(t *testing.T) {
	c := NewColor(0, 0, 0, 0)
	want := "*pixmatch.Color"
	typ := reflect.TypeOf(c).String()

	if typ != want {
		t.Errorf("Expected type %v got %v", want, typ)
	}
}

func TestColorRGBA(t *testing.T) {
	c := NewColor(0, 255, 255, 0)
	r, g, b, a := c.RGBA()
	want := []uint32{r, g, b, a}
	if r != c.R || g != c.G || b != c.B || a != c.A {
		t.Errorf("Expected %v got %v", want, c)
	}
}

func TestColorEquals(t *testing.T) {
	c1 := NewColor(0, 255, 255, 0)
	c2 := NewColor(0, 255, 255, 0)
	c3 := NewColor(0, 1, 244, 213)

	want := true
	res := c1.Equals(c2)
	if !res {
		t.Errorf("Expected %v got %v", want, res)
	}

	want = false
	res = c1.Equals(c3)
	if res {
		t.Errorf("Expected %v got %v", want, res)
	}
}

func TestImageYQI(t *testing.T) {
	color := NewColor(11, 22, 33, 44)
	y, i, q := color.YIQ()
	wantY, wantI, wantQ := 19.97145634, -10.09557868, 1.0964444699999998
	if y != wantY && i != wantI && q != wantQ {
		t.Errorf("Expected %v,%v,%v got %v,%v,%v",
			wantY, wantI, wantQ, y, i, q)
	}
}

func TestString(t *testing.T) {
	color := NewColor(123, 233, 12, 42)
	want := "(123,233,12,42)"
	res := color.String()
	if want != res {
		t.Errorf("Expected %v got %v", want, res)
	}
}
