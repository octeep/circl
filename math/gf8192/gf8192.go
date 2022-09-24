// Package gf8192 provides finite field arithmetic over GF(2^13).
package gf8192

// Gf is a field element of characteristic 2 modulo z^13 + z^4 + z^3 + z + 1
type Gf = uint16

const (
	gfBits = 13
	gfMask = (1 << gfBits) - 1
)

// Add two Gf elements together. Since an addition in Gf(2) is the same as XOR,
// this implementation uses a simple XOR for addition.
func Add(a, b Gf) Gf {
	return a ^ b
}

// Mul calculate the product of two Gf elements.
func Mul(a, b Gf) Gf {
	a64 := uint64(a)
	b64 := uint64(b)

	// if the LSB of b is 1, set tmp to a64, and 0 otherwise
	tmp := a64 & -(b64 & 1)

	// check if i-th bit of b64 is set, add a64 shifted by i bits if so
	for i := 1; i < gfBits; i++ {
		tmp ^= a64 * (b64 & (1 << i))
	}

	// polynomial reduction
	t := tmp & 0x1FF0000
	tmp ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)

	t = tmp & 0x000E000
	tmp ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)

	return uint16(tmp & gfMask)
}

// Sqr2 calculates a^4
func Sqr2(a Gf) Gf {
	a64 := uint64(a)
	a64 = (a64 | (a64 << 24)) & 0x000000FF000000FF
	a64 = (a64 | (a64 << 12)) & 0x000F000F000F000F
	a64 = (a64 | (a64 << 6)) & 0x0303030303030303
	a64 = (a64 | (a64 << 3)) & 0x1111111111111111

	t := a64 & 0x0001FF0000000000
	a64 ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	t = a64 & 0x000000FF80000000
	a64 ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	t = a64 & 0x000000007FC00000
	a64 ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	t = a64 & 0x00000000003FE000
	a64 ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)

	return uint16(a64 & gfMask)
}

// SqrMul calculates the product of a^2 and b
func SqrMul(a, b Gf) Gf {
	a64 := uint64(a)
	b64 := uint64(b)

	x := (b64 << 6) * (a64 & (1 << 6))
	a64 ^= a64 << 7
	x ^= b64 * (a64 & (0x04001))
	x ^= (b64 * (a64 & (0x08002))) << 1
	x ^= (b64 * (a64 & (0x10004))) << 2
	x ^= (b64 * (a64 & (0x20008))) << 3
	x ^= (b64 * (a64 & (0x40010))) << 4
	x ^= (b64 * (a64 & (0x80020))) << 5

	t := x & 0x0000001FF0000000
	x ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	t = x & 0x000000000FF80000
	x ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	t = x & 0x000000000007E000
	x ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)

	return uint16(x & gfMask)
}

// Sqr2Mul calculates the product of a^4 and b
func Sqr2Mul(a, b Gf) Gf {
	a64 := uint64(a)
	b64 := uint64(b)

	x := (b64 << 18) * (a64 & (1 << 6))
	a64 ^= a64 << 21
	x ^= b64 * (a64 & (0x010000001))
	x ^= (b64 * (a64 & (0x020000002))) << 3
	x ^= (b64 * (a64 & (0x040000004))) << 6
	x ^= (b64 * (a64 & (0x080000008))) << 9
	x ^= (b64 * (a64 & (0x100000010))) << 12
	x ^= (b64 * (a64 & (0x200000020))) << 15

	t := x & 0x1FF0000000000000
	x ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	t = x & 0x000FF80000000000
	x ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	t = x & 0x000007FC00000000
	x ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	t = x & 0x00000003FE000000
	x ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	t = x & 0x0000000001FE0000
	x ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)
	t = x & 0x000000000001E000
	x ^= (t >> 9) ^ (t >> 10) ^ (t >> 12) ^ (t >> 13)

	return uint16(x & gfMask)
}

// Inv calculates the multiplicative inverse of Gf element a
func Inv(a Gf) Gf {
	return Div(a, 1)
}

// Div calculates a / b
func Div(b, a Gf) Gf {
	tmp3 := SqrMul(b, b)         // b^3
	tmp15 := Sqr2Mul(tmp3, tmp3) // b^15 = b^(3*2*2+3)
	out := Sqr2(tmp15)
	out = Sqr2Mul(out, tmp15) // b^255 = b^(15*4*4+15)
	out = Sqr2(out)
	out = Sqr2Mul(out, tmp15) // b^4095 = b^(255*2*2*2*2+15)

	return SqrMul(out, a) // b^8190 = b^(4095*2) = b^-1
}
