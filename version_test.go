package pixmatch

import "testing"

func TestVersion(t *testing.T) {
	v := Version{1, 0, 0, ""}
	want := "1.0.0"
	res := v.String()
	if want != res {
		t.Errorf("Expected %v got %v", want, res)
	}

	v = Version{0, 0, 1, "alpha"}
	want = "0.0.1-alpha"
	res = v.String()
	if want != res {
		t.Errorf("Expected %v got %v", want, res)
	}
}
