package mceliece

// Returns (min(a, b), max(a, b)), executes in constant time
func minMaxI32(a, b int32) (int32, int32) {
	ab := b ^ a
	c := b - a
	c ^= ab & (c ^ b)
	c >>= 31
	c &= ab
	a ^= c
	b ^= c
	return a, b
}

// Returns (min(a, b), max(a, b)), executes in constant time
///
/// This differs from the C implementation, because the C implementation
/// only works for 63-bit integers.
///
/// Instead this implementation is based on
/// “side-channel effective overflow check of variable c”
/// from the book “Hacker's Delight” 2–13 Overflow Detection,
/// Section Unsigned Add/Subtract p. 40
func minMaxU64(a, b uint64) (uint64, uint64) {
	c := (^b & a) | ((^b | a) & (b - a))
	c = -(c >> 63)
	c &= a ^ b
	a ^= c
	b ^= c
	return a, b
}

// Reference: [djbsort](https://sorting.cr.yp.to/).
func int32Sort(x []int32, n int32) {
	if n < 2 {
		return
	}
	top := int32(1)
	for top < n-top {
		top += top
	}
	for p := top; p > 0; p >>= 1 {
		for i := int32(0); i < n-p; i++ {
			if (i & p) == 0 {
				min, max := minMaxI32(x[i], x[i+p])
				x[i] = min
				x[i+p] = max
			}
		}

		i := int32(0)
		for q := top; q > p; q >>= 1 {
			for ; i < n-q; i++ {
				if (i & p) == 0 {
					a := x[i+p]
					for r := q; r > p; r >>= 1 {
						min, max := minMaxI32(a, x[i+r])
						x[i+r] = max
						a = min
					}
					x[i+p] = a
				}
			}
		}
	}
}

// UInt64Sort sorts a slice of uint64
// Reference: [djbsort](https://sorting.cr.yp.to/).
func UInt64Sort(x []uint64, n int) {
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
				min, max := minMaxU64(x[i], x[i+p])
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
						min, max := minMaxU64(a, x[i+r])
						x[i+r] = max
						a = min
					}
					x[i+p] = a
				}
			}
		}
	}
}
