package mceliece

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/cloudflare/circl/internal/test"
)

type foo []int32

func (f foo) Len() int {
	return len(f)
}

func (f foo) Less(i, j int) bool {
	return f[i] < f[j]
}

func (f foo) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

//nolint:gosec
func TestSort(t *testing.T) {
	arr := make(foo, 314)
	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Int31()
	}

	int32Sort(arr)
	if !sort.IsSorted(arr) {
		want := make(foo, len(arr))
		copy(want, arr)
		sort.Sort(want)
		test.ReportError(t, arr, want)
	}
}
