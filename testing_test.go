package goutils

import "testing"

func TestAssertEquals(t *testing.T) {
	AssertEqual(t, "a", "a")
}

func TestAssertNotEquals(t *testing.T) {
	AssertNotEqual(t, "a", "")
}
