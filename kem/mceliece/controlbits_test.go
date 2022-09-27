package mceliece

import (
	"fmt"
	"testing"

	"github.com/cloudflare/circl/internal/test"
)

const (
	w = 8
	n = 1 << w
)

func TestCBRecursion(t *testing.T) {
	pi, err := FindTestDataI16("controlbits_kat3_mceliece348864_pi")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	out := make([]byte, 5888)
	fmt.Println("a")
	controlBitsFromPermutation(out, pi, 12, 4096)
	fmt.Println("b")
	want, err := FindTestDataByte("controlbits_kat3_mceliece348864_out_ref")
	test.ReportError(t, out, want)
}
