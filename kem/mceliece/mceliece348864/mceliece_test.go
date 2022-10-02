package mceliece348864

import (
	"crypto/sha256"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/cloudflare/circl/internal/nist"
	"github.com/cloudflare/circl/internal/test"
)

const katNum = 10

func fill(t *[sysT]gf, v gf) {
	for i := 0; i < sysT; i++ {
		t[i] = v
	}
}

func assertEq(t *testing.T, a *[sysT]gf, b []gf) {
	if !reflect.DeepEqual(a[:], b) {
		test.ReportError(t, b, a[:])
	}
}

func printHex(t *testing.T, writer io.Writer, header string, bytes []byte) {
	_, err := fmt.Fprint(writer, header)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fmt.Fprintf(writer, "%X", bytes)
	if err != nil {
		t.Fatal(err)
	}

	_, err = fmt.Fprintf(writer, "\n")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateKeyPair(t *testing.T) {
	entropy := [48]byte{}
	for i := 0; i < len(entropy); i++ {
		entropy[i] = byte(i)
	}
	rng := nist.NewDRBG(&entropy)

	digest := sha256.New()
	_, err := fmt.Fprintf(digest, "# kem/mceliece348864\n\n")
	if err != nil {
		t.Fatal(err)
	}
	scheme := Scheme()

	for i := 0; i < katNum; i++ {
		fmt.Fprintf(digest, "count = %d\n", i)

		s := [48]byte{}
		rng.Fill(s[:])

		dRng := nist.NewDRBG(&s)
		sessionSeed := make([]byte, 32)
		dRng.Fill(sessionSeed)
		pk, sk := scheme.DeriveKeyPair(sessionSeed)
		if !pk.Equal(sk.Public()) {
			t.Fatal("sk.Public() does not match pk")
		}
		pkBytes, err := pk.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		skBytes, err := sk.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		ct, ss, err := scheme.EncapsulateDeterministically(pk, s[:])
		if err != nil {
			t.Fatal(err)
		}
		dss, err := scheme.Decapsulate(sk, ct)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(dss, ss) {
			t.Logf("Error at count %d\n", i)
			test.ReportError(t, fmt.Sprintf("%X", dss), fmt.Sprintf("%X", ss))
		}
		printHex(t, digest, "seed = ", s[:])
		printHex(t, digest, "pk = ", pkBytes)
		printHex(t, digest, "sk = ", skBytes)
		printHex(t, digest, "ct = ", ct)
		printHex(t, digest, "ss = ", ss)
		fmt.Fprintf(digest, "\n")
	}

	if fmt.Sprintf("%x", digest.Sum(nil)) != "083224b827fc165a0f0e395e1905d7056ca309bf88a84c9b21ca658eddcbf140" {
		t.Fatal()
	}
}

func TestMinimalPolynomial(t *testing.T) {
	// tests data generated by Sage
	out := [sysT]gf{}

	in := [sysT]gf{
		1214, 685, 3954, 4010, 38, 1628, 1012, 2473, 3205, 2682, 2677, 3425, 3110, 1093, 2185, 3846, 3695,
		767, 1637, 2096, 505, 1631, 2771, 3982, 1826, 3355, 2038, 1636, 1221, 3228, 1473, 3371, 1443, 2985,
		386, 11, 2376, 529, 828, 2615, 1517, 2414, 3324, 3951, 460, 3457, 974, 2316, 2655, 2889, 2150, 1163,
		1612, 974, 758, 517, 1874, 2819, 1257, 2746, 1559, 1596, 1795, 740,
	}
	want := [sysT]gf{
		3991, 3480, 592, 686, 1616, 2086, 804, 3006, 1220, 3868, 2339, 1195, 3235, 3101, 1893, 1285, 280, 3093,
		1919, 1048, 458, 704, 954, 2844, 3679, 3228, 2270, 3886, 717, 1133, 1363, 3026, 3005, 241, 829, 316, 3951,
		2312, 2934, 2610, 1465, 2208, 1915, 2534, 1487, 1266, 3039, 1729, 1585, 3671, 1597, 1189, 3907, 956, 3519,
		3007, 3677, 2253, 1595, 1293, 2029, 2971, 1370, 1240,
	}
	test.CheckOk(minimalPolynomial(&out, &in), "minimalPolynomial failed", t)
	assertEq(t, &want, out[:])

	in = [sysT]gf{
		2871, 2450, 516, 137, 3881, 2283, 3696, 3941, 921, 2528, 2099, 3, 3880, 333, 3277, 3787, 141, 2552, 3086,
		1178, 612, 3233, 456, 1222, 3546, 1205, 2786, 877, 2183, 1318, 300, 3583, 1996, 3838, 2263, 3690, 1449,
		3487, 1005, 2206, 525, 779, 1220, 3983, 1697, 3521, 3307, 2752, 1003, 2322, 4022, 1426, 2106, 360, 1261,
		3268, 2050, 3243, 189, 2432, 4048, 362, 2431, 2441,
	}
	want = [sysT]gf{
		2143, 2914, 2162, 2520, 1157, 4069, 1048, 237, 3123, 2684, 3638, 3812, 3771, 371, 3156, 2345, 562, 3051,
		3702, 3364, 2452, 352, 2525, 2061, 562, 2696, 868, 108, 295, 910, 1263, 1763, 266, 2752, 1784, 550, 2606,
		703, 32, 387, 2213, 4021, 2388, 2996, 325, 2429, 1533, 3940, 1817, 766, 422, 3198, 289, 3519, 3990, 2593,
		165, 3710, 3932, 3734, 3981, 1382, 2358, 2275,
	}
	test.CheckOk(minimalPolynomial(&out, &in), "minimalPolynomial failed", t)
	assertEq(t, &want, out[:])
}

//nolint
func TestPolyMul(t *testing.T) {
	res := [sysT]gf{}
	arg1 := [sysT]gf{}
	arg2 := [sysT]gf{}

	fill(&arg1, 0)
	fill(&arg2, 0)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0,
	})

	fill(&arg1, 0)
	fill(&arg2, 1)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
	})

	fill(&arg1, 1)
	fill(&arg2, 0)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1,
	})

	fill(&arg1, 0)
	fill(&arg2, 5)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		1, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
		5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
		5, 5, 5, 5, 5, 5, 5, 5, 5,
	})

	fill(&arg1, 5)
	fill(&arg2, 0)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		1, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
		5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
		5, 5, 5, 5, 5, 5, 5, 5, 5,
	})

	fill(&arg1, 0)
	fill(&arg2, 1024)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		1, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024,
		1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024,
		1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024,
		1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024,
		1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024,
	})

	fill(&arg1, 1024)
	fill(&arg2, 0)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		1, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024,
		1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024,
		1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024,
		1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024,
		1024, 1024, 1024, 1024, 1024, 1024, 1024, 1024,
	})

	fill(&arg1, 2)
	fill(&arg2, 6)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		25, 16, 28, 4, 28, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4,
		16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4,
		16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4,
	})

	fill(&arg1, 6)
	fill(&arg2, 2)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		25, 16, 28, 4, 28, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4,
		16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4,
		16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4, 16, 4,
	})

	fill(&arg1, 3)
	fill(&arg2, 8)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		49, 35, 59, 11, 59, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11,
		35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35,
		11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11,
		35, 11,
	})

	fill(&arg1, 8)
	fill(&arg2, 3)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		49, 35, 59, 11, 59, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11,
		35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35,
		11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11, 35, 11,
		35, 11,
	})

	fill(&arg1, 125)
	fill(&arg2, 19)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		3759, 2455, 3776, 110, 3776, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110,
		2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455,
		110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110,
		2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455,
		110, 2455, 110, 2455, 110,
	})

	fill(&arg1, 19)
	fill(&arg2, 125)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		3759, 2455, 3776, 110, 3776, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110,
		2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455,
		110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110,
		2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455, 110, 2455,
		110, 2455, 110, 2455, 110,
	})

	fill(&arg1, 125)
	fill(&arg2, 37)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		3162, 554, 3075, 88, 3075, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554,
		88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554,
		88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554,
		88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88,
	})

	fill(&arg1, 37)
	fill(&arg2, 125)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		3162, 554, 3075, 88, 3075, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554,
		88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554,
		88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88, 554,
		88, 554, 88, 554, 88, 554, 88, 554, 88, 554, 88,
	})

	fill(&arg1, 4095)
	fill(&arg2, 1)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		4086, 4086, 9, 4094, 9, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
		4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
		4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
		4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
		4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
	})

	fill(&arg1, 1)
	fill(&arg2, 4095)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		4086, 4086, 9, 4094, 9, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
		4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
		4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
		4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
		4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
	})

	fill(&arg1, 8191)
	fill(&arg2, 1)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		4068, 4068, 18, 4087, 18, 4087, 4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087,
		4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087,
		4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087,
		4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087,
		4068, 4087, 4068, 4087, 4068, 4087, 4068, 4087,
	})

	fill(&arg1, 1)
	fill(&arg2, 8191)
	arg1[0] = 1
	arg2[0] = 1
	polyMul(&res, &arg1, &arg2)
	assertEq(t, &res, []gf{
		4086, 4086, 9, 4094, 9, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
		4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
		4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
		4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
		4086, 4094, 4086, 4094, 4086, 4094, 4086, 4094,
	})
}
