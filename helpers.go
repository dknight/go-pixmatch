package pixmatch

// I really don't want to overload code with generic, conversions, dependencies
// or ever dependencies. So let just be these 2 simple functions.

// intMin returns the minimum int of 2 numbers.
func intMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// intMax returns the maximum int of 2 numbers.
func intMax(x, y int) int {
	if x > y {
		return x
	}
	return y
}
