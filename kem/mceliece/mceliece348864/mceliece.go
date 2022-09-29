package mceliece348864

import (
	cryptoRand "crypto/rand"
	"io"

	"github.com/cloudflare/circl/kem/mceliece"

	"github.com/cloudflare/circl/internal/sha3"

	"github.com/cloudflare/circl/math/gf4096"
)

const (
	sysT       = 64 // F(y) is 64 degree
	gfBits     = gf4096.GfBits
	gfMask     = gf4096.GfMask
	unusedBits = 16 - gfBits
	sysN       = 3488
	condBytes  = (1 << (gfBits - 4)) * (2*gfBits - 1)
	irrBytes   = sysT * 2
	pkNRows    = sysT * gfBits
	pkNCols    = sysN - pkNRows
	pkRowBytes = (pkNCols + 7) / 8
	syndBytes  = (pkNRows + 7) / 8
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

func generateKeyPair(pk, sk []byte, rand io.Reader) error {
	const (
		irrPolys  = sysN/8 + (1<<gfBits)*4
		seedIndex = sysN/8 + (1<<gfBits)*4 + sysT*2
		permIndex = sysN / 8
		sBase     = 32 + 8 + irrBytes + condBytes
	)

	seed := [33]byte{64}
	r := [sysN/8 + (1<<gfBits)*4 + sysT*2 + 32]byte{}

	f := [sysT]Gf{}
	irr := [sysT]Gf{}
	perm := [1 << gfBits]uint32{}
	pi := [1 << gfBits]int16{}
	pivots := uint64(0)

	if rand == nil {
		rand = cryptoRand.Reader
	}
	if _, err := io.ReadFull(rand, seed[1:]); err != nil {
		return err
	}

	for {
		// expanding and updating the seed
		err := shake256(r[:], seed[0:33])
		if err != nil {
			return err
		}

		copy(sk[:32], seed[1:])
		copy(seed[1:], r[len(r)-32:])

		temp := r[irrPolys:seedIndex]
		for i := 0; len(temp) > 0; i++ {
			f[i] = loadGf(temp)
			temp = temp[2:]
		}

		if !minimalPolynomial(&irr, &f) {
			continue
		}

		temp = sk[32+8 : 32+8+2*sysT]
		for i := 0; len(temp) > 0; i++ {
			storeGf(temp, irr[i])
			temp = temp[2:]
		}

		// generating permutation
		temp = r[permIndex:irrPolys]
		for i := 0; len(temp) > 0; i++ {
			perm[i] = load4(temp)
			temp = temp[4:]
		}

		// TODO: pk_gen
		mceliece.ControlBitsFromPermutation(sk[32+8+irrBytes:], pi[:], gfBits, 1<<gfBits)
		copy(sk[sBase:sBase+sysN/8], r[0:sysN/8])
		store8(sk[32:40], pivots)
		return nil
	}
}

func shake256(output []byte, input []byte) error {
	shake := sha3.NewShake256()
	_, err := shake.Write(input)
	if err != nil {
		return err
	}
	_, err = shake.Read(output)
	if err != nil {
		return err
	}
	return nil
}

func storeGf(dest []byte, a Gf) {
	dest[0] = byte(a & 0xFF)
	dest[1] = byte(a >> 8)
}

func loadGf(src []byte) Gf {
	a := uint16(src[1])
	a <<= 8
	a |= uint16(src[0])
	return a & gfMask
}

func load4(in []byte) uint32 {
	ret := uint32(in[3])
	for i := 2; i >= 0; i-- {
		ret <<= 8
		ret |= uint32(in[i])
	}
	return ret
}

func store8(out []byte, in uint64) {
	out[0] = byte((in >> 0x00) & 0xFF)
	out[1] = byte((in >> 0x08) & 0xFF)
	out[2] = byte((in >> 0x10) & 0xFF)
	out[3] = byte((in >> 0x18) & 0xFF)
	out[4] = byte((in >> 0x20) & 0xFF)
	out[5] = byte((in >> 0x28) & 0xFF)
	out[6] = byte((in >> 0x30) & 0xFF)
	out[7] = byte((in >> 0x38) & 0xFF)
}

func load8(in []byte) uint64 {
	ret := uint64(in[7])
	for i := 6; i >= 0; i-- {
		ret <<= 8
		ret |= uint64(in[i])
	}
	return ret
}

func bitRev(a Gf) Gf {
	a = ((a & 0x00FF) << 8) | ((a & 0xFF00) >> 8)
	a = ((a & 0x0F0F) << 4) | ((a & 0xF0F0) >> 4)
	a = ((a & 0x3333) << 2) | ((a & 0xCCCC) >> 2)
	a = ((a & 0x5555) << 1) | ((a & 0xAAAA) >> 1)

	return a >> unusedBits
}
