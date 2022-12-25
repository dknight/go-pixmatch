package pixmatch

import (
	"reflect"
	"testing"
)

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	want := defaultOptions
	if !reflect.DeepEqual(*opts, defaultOptions) {
		t.Errorf("Expected %+v got %+v", want, *opts)
	}
}
