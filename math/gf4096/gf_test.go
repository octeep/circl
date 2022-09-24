package gf4096

import (
	"fmt"
	"testing"

	"github.com/cloudflare/circl/internal/test"
)

type (
	tadd func(x, y Gf) Gf
	tmul func(x, y Gf) Gf
	tsqr func(x Gf) Gf
	tinv func(x Gf) Gf
	tdiv func(x, y Gf) Gf
)

func assertEq(t *testing.T, a, b Gf) {
	test.CheckOk(a == b, fmt.Sprintf("%d != %d", a, b), t)
}

func TestGeneric(t *testing.T) {
	t.Run("Add", func(t *testing.T) { testAdd(t, Add) })
	t.Run("Mul", func(t *testing.T) { testMul(t, Mul) })
	t.Run("Sqr", func(t *testing.T) { testSqr(t, Sqr) })
	t.Run("Inv", func(t *testing.T) { testInv(t, Inv) })
	t.Run("Div", func(t *testing.T) { testDiv(t, Div) })
}

func testDiv(t *testing.T, div tdiv) {
	assertEq(t, div(1, 6733), 2637)
	assertEq(t, div(2, 0), 0)
	assertEq(t, div(2, 4), 2)
	assertEq(t, div(2, 4096), 0)
	assertEq(t, div(3, 9), 7)
	assertEq(t, div(5, 4591), 99)
	assertEq(t, div(550, 10), 3344)
	assertEq(t, div(5501, 3), 1763)
}

func testInv(t *testing.T, inv tinv) {
	assertEq(t, inv(0), 0)
	assertEq(t, inv(1), 1)
	assertEq(t, inv(2), 2052)
	assertEq(t, inv(3), 4088)
	assertEq(t, inv(4), 1026)
	assertEq(t, inv(4095), 1539)
	assertEq(t, inv(4096), 0)
	assertEq(t, inv(8191), 1539)
	assertEq(t, inv(8192), 0)
	assertEq(t, inv(0xFFFF), 1539)
}

func testSqr(t *testing.T, sqr tsqr) {
	assertEq(t, sqr(0), 0)
	assertEq(t, sqr(1), 1)
	assertEq(t, sqr(2), 4)
	assertEq(t, sqr(3), 5)
	assertEq(t, sqr(4), 16)
	assertEq(t, sqr(4095), 2746)
	assertEq(t, sqr(4096), 0)
	assertEq(t, sqr(8191), 2746)
	assertEq(t, sqr(8192), 0)
	assertEq(t, sqr(0xFFFF), 2746)
}

func testMul(t *testing.T, mul tmul) {
	assertEq(t, mul(0, 0), 0)
	assertEq(t, mul(0, 1), 0)
	assertEq(t, mul(1, 0), 0)
	assertEq(t, mul(0, 5), 0)
	assertEq(t, mul(5, 0), 0)
	assertEq(t, mul(0, 1024), 0)
	assertEq(t, mul(1024, 0), 0)
	assertEq(t, mul(2, 6), 12)
	assertEq(t, mul(6, 2), 12)
	assertEq(t, mul(3, 8), 24)
	assertEq(t, mul(8, 3), 24)
	assertEq(t, mul(125, 19), 1879)
	assertEq(t, mul(19, 125), 1879)
	assertEq(t, mul(125, 37), 3625)
	assertEq(t, mul(37, 125), 3625)
	assertEq(t, mul(4095, 1), 4095)
	assertEq(t, mul(1, 4095), 4095)
	assertEq(t, mul(8191, 1), 4086)
	assertEq(t, mul(1, 8191), 4095)
	assertEq(t, mul(550, 3344), 10)
	assertEq(t, mul(3344, 550), 10)
}

func testAdd(t *testing.T, add tadd) {
	assertEq(t, add(0x0000, 0x0000), 0x0000)
	assertEq(t, add(0x0000, 0x0001), 0x0001)
	assertEq(t, add(0x0001, 0x0000), 0x0001)
	assertEq(t, add(0x0001, 0x0001), 0x0000)
	assertEq(t, add(0x000F, 0x0000), 0x000F)
	assertEq(t, add(0x000F, 0x0001), 0x000E)
	assertEq(t, add(0x00FF, 0x0100), 0x01FF)
	assertEq(t, add(0xF0F0, 0x0F0F), 0xFFFF)
}
