package idiv

import (
	"math"
	"testing"
)

func testPair(t *testing.T, a, b int32) {
	q1, r1 := Div1(a, b)
	q2, r2 := Div2(a, b)
	if q1 != q2 || r1 != r2 {
		t.Errorf("Div1 and Div2 returned different values for dividend=%d, divisor=%d. Expected: %d, %d, Actual: %d, %d", a, b, q1, r1, q2, r2)
	}
}

func testSignedPair(t *testing.T, a, b int32) {
	testPair(t, a, b)
	testPair(t, a, -b)
	testPair(t, -a, b)
	testPair(t, -a, -b)
}

func TestDiv1AndDiv2(t *testing.T) {
	for i := range 100 {
		for j := range 100 {
			if j == 0 {
				continue
			}
			testSignedPair(t, int32(i), int32(j))
		}
	}
}

func TestBoundary(t *testing.T) {
	a := int32(math.MinInt32)
	b := int32(1 << 16)
	testSignedPair(t, a, b)
	q, r := Div2(a, b)
	expected := int32(-(1 << 15))
	if q != expected || r != 0 {
		t.Errorf("Div2 returned wrong values for dividend=math.MinInt64, divisor=-1. Expected: %d, %d, Actual: %d, %d", expected, 0, q, r)
	}
}

func TestBoundaryBadImpl(t *testing.T) {
	a := int32(math.MinInt32)
	b := int32(1)
	q, r := Div3(a, b)
	expected := a
	if q != expected || r != 0 {
		t.Errorf("Div3 returned wrong values for dividend=math.MinInt64, divisor=-1. Expected: %d, %d, Actual: %d, %d", expected, 0, q, r)
	}
}

func TestAbs(t *testing.T) {
	if r := abs2(math.MinInt32); r < 0 {
		t.Errorf("abs1 returned wrong value: %d", r)
	}
}
