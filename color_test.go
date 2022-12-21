package pixmatch

import (
	"reflect"
	"testing"
)

func TestNewColor(t *testing.T) {
	c := NewColor(0, 0, 0, 0)
	exp := "*pixmatch.Color"
	typ := reflect.TypeOf(c).String()

	if typ != exp {
		t.Error("Expected type", exp, "got", typ)
	}
}

func TestColorRGBA(t *testing.T) {
	c := NewColor(0, 255, 255, 0)
	r, g, b, a := c.RGBA()
	if r != c.R || g != c.G || b != c.B || a != c.A {
		t.Errorf("Expected %v got %v", []uint32{r, g, b, a}, c)
	}
}
