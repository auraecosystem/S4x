package main

/*
#cgo LDFLAGS: -L. -lcompute
int computePremiumScore(int base, int multiplier);
*/
import "C"

import (
	"context"
	// ... rest of your imports stay the same
)

// Example usage inside an endpoint handler:
func CalculateScore(base, mult int) int {
	return int(C.computePremiumScore(C.int(base), C.int(mult)))
}
