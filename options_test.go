package pixmatch

import (
	"reflect"
	"testing"
)

func TestNewOptions(t *testing.T) {
	opts := NewOptions()
	if !reflect.DeepEqual(*opts, defaultOptions) {
		t.Errorf("Expected %+v got %+v", defaultOptions, *opts)
	}
}
