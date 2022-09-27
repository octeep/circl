package mceliece

import (
	"math/rand"
	"testing"
)

const (
	w = 8
	n = 1 << w
)

func TestCBRecursion(t *testing.T) {
	temp := [2 * n]int32{}
	out := [2 * n]uint8{}
	pi32 := rand.Perm(n)
	pi := [n]uint16{}
	for i := 0; i < n; i++ {
		pi[i] = uint16(pi32[i])
	}
	CBRecursion(out[:], 0, 1, pi[:], w, n, temp[:])
}
