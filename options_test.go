package pixmatch

import (
	"image/color"
	"reflect"
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	want := defaultOptions
	if !reflect.DeepEqual(*opts, *defaultOptions) {
		t.Errorf("Expected %+v got %+v", *want, *opts)
	}
}

func TestNewOptions(t *testing.T) {
	opts := NewOptions()
	if opts == nil {
		t.Errorf("Cannot create options")
	}
}

func TestResolveDiffColor(t *testing.T) {
	opts := DefaultOptions()
	res := opts.ResolveDiffColor()
	if opts.DiffColorAlt != nil {
		t.Errorf("Expected %+v got %+v", nil, res)
	}

	opts.DiffColorAlt = color.RGBA{255, 255, 0, 255}
	want := opts.ResolveDiffColor()
	if opts.DiffColorAlt == nil {
		t.Errorf("Expected %+v got %+v", want, res)
	}
}
