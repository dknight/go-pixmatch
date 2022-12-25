package pixmatch

import "testing"

func Test_intMin(t *testing.T) {
	a, b := 10, 2
	if intMin(a, b) != b {
		t.Errorf("Expected %v got %v", b, a)
	}

	a, b = 0, 1
	if intMin(a, b) != a {
		t.Errorf("Expected %v got %v", a, b)
	}
}

func Test_intMax(t *testing.T) {
	a, b := 10, 2
	if intMax(a, b) != a {
		t.Errorf("Expected %v got %v", a, b)
	}

	a, b = 0, 1
	if intMax(a, b) != b {
		t.Errorf("Expected %v got %v", b, a)
	}
}
