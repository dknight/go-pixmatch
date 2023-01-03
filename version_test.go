package pixmatch

import (
	"fmt"
	"testing"
)

func TestNewVersion(t *testing.T) {
	v := Version{1, 0, 0, ""}
	want := "1.0.0"
	res := v.String()
	if want != res {
		t.Errorf("Expected %v got %v", want, res)
	}
}

func TestGetVersion(t *testing.T) {
	res := GetVersion()
	want := fmt.Sprintf("%s", currentVersion)
	if want != res {
		t.Errorf("Expected %v got %v", want, res)
	}
}
