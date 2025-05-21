package goutils

import (
	"fmt"
	"testing"
)

func TestAssertEquals(t *testing.T) {
	AssertEqual(t, "a", "a")
}

func TestAssertNotEquals(t *testing.T) {
	AssertNotEqual(t, "a", "")
}

func TestArrayEqual(t *testing.T) {
	arrIs := []byte{1, 2, 3}
	AssertBytesEqual(t, arrIs, arrIs)
	arrShould := []byte{1, 2, 3}
	AssertBytesEqual(t, arrIs, arrShould)
	arrShouldnt := []byte{3, 2, 1}
	t2 := &testing.T{}
	Assert(t, !AssertBytesEqual(t2, arrIs, arrShouldnt))
	Assert(t, !AssertBytesEqual(t2, arrIs, nil))
	Assert(t, !AssertBytesEqual(t2, nil, arrIs))
	empty := []byte{}
	Assert(t, !AssertBytesEqual(t2, arrIs, empty))
	AssertBytesEqual(t2, empty, empty)
}

func TestAssertWithin(t *testing.T) {
	b := float32(0)
	for i := 0; i < 10; i++ {
		b += float32(0.1)
	}
	Assert(t, !(b == float32(1.0)))
	AssertWithin(t, b, 1.0, 0.0001)

	AssertWithin(t, 10, 12, 2)
}

func TestAssertArrayWithin(t *testing.T) {
	b := []float32{0.0, 1.2, 2.3, 4.5, 6.7}
	for i := 0; i != 10; i++ {
		for j := 0; j < 5; j++ {
			b[j] += 0.1
		}
	}
	Assert(t, !(b[0] == float32(1.0)))
	AssertArrayWithin(t, b, []float32{1.0, 2.2, 3.3, 5.5, 7.7}, 0.0001)

	AssertArrayWithin(t, []uint64{10, 10, 10, 10, 10}, []uint64{12, 12, 12, 12, 12}, 2)
}

func TestAssertZeroValue(t *testing.T) {
	var a *int = nil
	var s *string = nil
	var e error = nil
	var i int = 0
	AssertZeroValue(t, a)
	AssertZeroValue(t, s)
	AssertZeroValue(t, e)
	AssertZeroValue(t, i)
	AssertZeroValue(t, "")
	AssertZeroValue(t, false)
}
func TestAssertNil(t *testing.T) {
	var a *int = nil
	var s *string = nil
	var e error = nil
	var i int = 0
	AssertNil(t, a)
	AssertNil(t, s)
	AssertNil(t, e)
	AssertNil(t, i)
}

func TestAssertNotZeroValue(t *testing.T) {
	a := 1
	s := "a"
	e := fmt.Errorf("a")
	AssertNotZeroValue(t, a)
	AssertNotZeroValue(t, s)
	AssertNotZeroValue(t, e)
	AssertNotZeroValue(t, true)
}
func TestAssertNotNil(t *testing.T) {
	a := 1
	s := "a"
	e := fmt.Errorf("a")
	AssertNotNil(t, a)
	AssertNotNil(t, s)
	AssertNotNil(t, e)
	AssertNotNil(t, true)
}
