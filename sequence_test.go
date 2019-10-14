package main

import (
	"fmt"
	"testing"
)

func TestFindSequenceAt(t *testing.T) {
	x := findSequenceAt(0)
	y := findSequenceAt(5)
	z := findSequenceAt(6)

	fmt.Printf("value of x is %d\n", x)
	fmt.Printf("value of y is %d\n", y)
	fmt.Printf("value of z is %d\n", z)

	if 3 != findSequenceAt(0) {
		t.Error("findSequenceAt of 0 should be '1' but have", x)
	}
	if 33 != findSequenceAt(5) {
		t.Error("findSequenceAt of 0 should be '1' but have", y)
	}
	if 45 != findSequenceAt(6) {
		t.Error("findSequenceAt of 0 should be '1' but have", z)
	}
}
