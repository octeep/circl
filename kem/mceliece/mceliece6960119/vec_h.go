package mceliece6960119

func vecSetBits(b uint64) uint64 {
	ret := -b
	return ret
}

func vecSet116b(v uint16) uint64 {
	ret := uint64(v)
	ret |= ret << 16
	ret |= ret << 32

	return ret
}

func vecCopy(out, in []uint64) {
	for i := 0; i < gfBits; i++ {
		out[i] = in[i]
	}
}

func vecOrReduce(a []uint64) uint64 {
	ret := a[0]
	for i := 1; i < gfBits; i++ {
		ret |= a[i]
	}

	return ret
}

func vecTestZ(a uint64) int {
	a |= a >> 32
	a |= a >> 16
	a |= a >> 8
	a |= a >> 4
	a |= a >> 2
	a |= a >> 1

	return int((a & 1) ^ 1)
}
