package mceliece

// Returns (min(a, b), max(a, b)), executes in constant time
func minMax(a, b int32) (int32, int32) {
	ab := b ^ a
	c := b - a
	c ^= ab & (c ^ b)
	c >>= 31
	c &= ab
	a ^= c
	b ^= c
	return a, b
}

// Reference: [djbsort](https://sorting.cr.yp.to/).
func int32Sort(x []int32) {
	n := len(x)
	if n < 2 {
		return
	}
	top := 1
	for top < n-top {
		top += top
	}
	for p := top; p > 0; p >>= 1 {
		for i := 0; i < n-p; i++ {
			if (i & p) == 0 {
				min, max := minMax(x[i], x[i+p])
				x[i] = min
				x[i+p] = max
			}
		}

		i := 0
		for q := top; q > p; q >>= 1 {
			for ; i < n-q; i++ {
				if (i & p) == 0 {
					a := x[i+p]
					for r := q; r > p; r >>= 1 {
						min, max := minMax(a, x[i+r])
						x[i+r] = max
						a = min
					}
					x[i+p] = a
				}
			}
		}
	}
}
