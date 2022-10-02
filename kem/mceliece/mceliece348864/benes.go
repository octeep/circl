package mceliece348864

func transpose64x64Inplace(out *[64]uint64) {
	masks := [6][2]uint64{
		{0x5555555555555555, 0xAAAAAAAAAAAAAAAA},
		{0x3333333333333333, 0xCCCCCCCCCCCCCCCC},
		{0x0F0F0F0F0F0F0F0F, 0xF0F0F0F0F0F0F0F0},
		{0x00FF00FF00FF00FF, 0xFF00FF00FF00FF00},
		{0x0000FFFF0000FFFF, 0xFFFF0000FFFF0000},
		{0x00000000FFFFFFFF, 0xFFFFFFFF00000000},
	}

	for d := 5; d >= 0; d-- {
		s := 1 << d
		for i := 0; i < 64; i += s * 2 {
			for j := i; j < i+s; j++ {
				x := (out[j] & masks[d][0]) | ((out[j+s] & masks[d][0]) << s)
				y := ((out[j] & masks[d][1]) >> s) | (out[j+s] & masks[d][1])

				out[j+0] = x
				out[j+s] = y
			}
		}
	}
}

func layer(data, bits []uint64, lgs int) {
	index := 0
	s := 1 << lgs
	for i := 0; i < 64; i += s * 2 {
		for j := i; j < i+s; j++ {
			d := data[j] ^ data[j+s]
			d &= bits[index]
			index++
			data[j] ^= d
			data[j+s] ^= d
		}
	}
}

func applyBenes(r *[512]byte, bits *[5888]byte) {
	bs := [64]uint64{}
	cond := [64]uint64{}
	for i := 0; i < 64; i++ {
		bs[i] = load8(r[i*8:])
	}

	transpose64x64Inplace(&bs)

	for low := 0; low <= 5; low++ {
		for i := 0; i < 64; i++ {
			cond[i] = uint64(load4(bits[low*256+i*4:]))
		}
		transpose64x64Inplace(&cond)
		layer(bs[:], cond[:], low)
	}

	transpose64x64Inplace(&bs)

	for low := 0; low <= 5; low++ {
		for i := 0; i < 32; i++ {
			cond[i] = load8(bits[(low+6)*256+i*8:])
		}
		layer(bs[:], cond[:], low)
	}
	for low := 4; low >= 0; low-- {
		for i := 0; i < 32; i++ {
			cond[i] = load8(bits[(4-low+6+6)*256+i*8:])
		}
		layer(bs[:], cond[:], low)
	}

	transpose64x64Inplace(&bs)

	for low := 5; low >= 0; low-- {
		for i := 0; i < 64; i++ {
			cond[i] = uint64(load4(bits[(5-low+6+6+5)*256+i*4:]))
		}
		transpose64x64Inplace(&cond)
		layer(bs[:], cond[:], low)
	}
	transpose64x64Inplace(&bs)

	for i := 0; i < 64; i++ {
		store8(r[i*8:], bs[i])
	}
}
