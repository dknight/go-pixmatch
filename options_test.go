package pixmatch

import (
	"image/color"
	"reflect"
	"testing"
)

func TestNewOptions(t *testing.T) {
	opts := NewOptions()
	want := defaultOptions
	if !reflect.DeepEqual(*opts, defaultOptions) {
		t.Errorf("Expected %+v got %+v", want, opts)
	}
}

func TestResolveDiffColor(t *testing.T) {
	opts := NewOptions()
	res := opts.ResolveDiffColor()
	if opts.DiffColorAlt != nil {
		t.Errorf("Expected %+v got %+v", nil, res) // diff color
	}

	opts.DiffColorAlt = color.RGBA{255, 255, 0, 255}
	want := opts.ResolveDiffColor()
	if opts.DiffColorAlt == nil {
		t.Errorf("Expected %+v got %+v", want, res)
	}
}
