package goutils

import "testing"

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

}
