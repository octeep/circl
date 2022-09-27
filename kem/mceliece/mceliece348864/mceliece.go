package mceliece348864

import (
	"github.com/cloudflare/circl/math/gf4096"
)

const (
	sysT   = 64 // F(y) is 64 degree
	gfBits = gf4096.GfBits
)

type Gf = gf4096.Gf

// check if element is 0, returns a mask with all bits set if so, and 0 otherwise
func isZeroMask(element Gf) uint16 {
	t := uint32(element) - 1
	t >>= 19
	return uint16(t)
}

// calculate the minimal polynomial of f and store it in out
func minimalPolynomial(out *[sysT]Gf, f *[sysT]Gf) bool {
	mat := [sysT + 1][sysT]Gf{}
	mat[0][0] = 1
	for i := 1; i < sysT; i++ {
		mat[0][i] = 0
	}

	for i := 0; i < sysT; i++ {
		mat[1][i] = f[i]
	}

	for i := 2; i <= sysT; i++ {
		polyMul(&mat[i], &mat[i-1], f)
	}

	for j := 0; j < sysT; j++ {
		for k := j + 1; k < sysT; k++ {
			mask := isZeroMask(mat[j][j])
			// if mat[j][j] is not zero, add mat[c..sysT+1][k] to mat[c][j]
			// do nothing otherwise
			for c := j; c <= sysT; c++ {
				mat[c][j] ^= mat[c][k] & mask
			}
		}

		if mat[j][j] == 0 {
			return false
		}

		inv := gf4096.Inv(mat[j][j])
		for c := 0; c <= sysT; c++ {
			mat[c][j] = gf4096.Mul(mat[c][j], inv)
		}

		for k := 0; k < sysT; k++ {
			if k != j {
				t := mat[j][k]
				for c := 0; c <= sysT; c++ {
					mat[c][k] ^= gf4096.Mul(mat[c][j], t)
				}
			}
		}
	}

	for i := 0; i < sysT; i++ {
		out[i] = mat[sysT][i]
	}

	return true
}

// calculate the product of a and b in Fq^t
func polyMul(out *[sysT]Gf, a *[sysT]Gf, b *[sysT]Gf) {
	product := [sysT*2 - 1]Gf{}
	for i := 0; i < sysT; i++ {
		for j := 0; j < sysT; j++ {
			product[i+j] ^= gf4096.Mul(a[i], b[j])
		}
	}

	for i := (sysT - 1) * 2; i >= sysT; i-- {
		// polynomial reduction
		product[i-sysT+3] ^= product[i]
		product[i-sysT+1] ^= product[i]
		product[i-sysT] ^= gf4096.Mul(product[i], 2)
	}

	for i := 0; i < sysT; i++ {
		out[i] = product[i]
	}
}
