package pixmatch

import (
	"os"
	"testing"
)

func TestNewOutputToFile(t *testing.T) {
	tmp, err := os.CreateTemp("", "tmp")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmp.Name())

	want, err := NewOutputToFile(tmp.Name(), 10, 10)
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

	out, err := NewOutputToFile("", 10, 10)
	if out != nil {
		t.Error(err)
	}

	out, err = NewOutputToFile(tmp.Name(), 10, 10)
	if out == nil {
		t.Error(err)
	}

	err = out.Save("png")
	if err != nil {
		t.Errorf("Cannot save PNG output")
	}

	err = out.Save("gif")
	if err != nil {
		t.Errorf("Cannot save GIF output")
	}

	err = out.Save("jpeg")
	if err != nil {
		t.Errorf("Cannot save JPEG output")
	}
}
