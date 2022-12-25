package pixmatch

import (
	"os"
	"testing"
)

func TestNewOutput(t *testing.T) {
	tmp, err := os.CreateTemp("", "tmp")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmp.Name())

	want, err := NewOutput(tmp.Name(), 10, 10)
	if err != nil {
		t.Error(err)
	}
	if want == nil {
		t.Errorf("Expected %v got %v", want, nil)
	}
}

func TestSave(t *testing.T) {
	tmp, err := os.CreateTemp("", "tmp")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmp.Name())

	out, err := NewOutput("", 10, 10)
	if out != nil {
		t.Error(err)
	}

	out, err = NewOutput(tmp.Name(), 10, 10)
	if out == nil {
		t.Error(err)
	}

	err = out.Save()
	if err != nil {
		t.Errorf("Cannot save output")
	}
}
