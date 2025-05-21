package goutils

import "testing"

func Assert(t *testing.T, value bool) bool {
	t.Helper()
	if !value {
		t.Errorf("(%s)\n    not true: >%#v<\n", t.Name(), value)
		return false

	}
	return true
}

func AssertEqual(t *testing.T, is, should interface{}) bool {
	t.Helper()
	if is != should {
		t.Errorf("(%s)\n    is: >%#v<\nshould: >%#v<", t.Name(), is, should)
		return false
	}
	return true
}

func AssertArrayEqual[C comparable](t *testing.T, is, should []C) bool {
	t.Helper()
	if is == nil && should == nil {
		return true
	}
	if is == nil {
		t.Errorf("arr is nil, should: %v", should)
		return false
	} else if should == nil {
		t.Errorf("should be nil, is: %v", is)
		return false
	} else {
		if len(is) != len(should) {
			t.Errorf("unequal length: %d (%d)", len(is), len(should))
			return false
		}
	}
	for i, _ := range is {
		if is[i] != should[i] {
			t.Errorf("element #%v is:%v should:%v", i, is[i], should[i])
			return false
		}
	}
	return true

}
func AssertBytesEqual(t *testing.T, is, should []byte) bool {
	t.Helper()
	return AssertArrayEqual(t, is, should)
}

func AssertNotEqual(t *testing.T, is, shouldnt interface{}) bool {
	t.Helper()
	if is == shouldnt {
		t.Errorf("(%s)\n    is: >%#v<\nshould: >%#v<", t.Name(), is, shouldnt)
		return false
	}
	return true
}

// Assert that the provided value is not the zero value of its type.
func AssertNotZeroValue[N comparable](t *testing.T, val N) bool {
	t.Helper()
	var zero N
	if val == zero {
		t.Errorf("(%s)\n    not zero value: >%#v<\n", t.Name(), val)
		return false
	}
	return true
}

func AssertNotNil[N comparable](t *testing.T, val N) bool {
	// Deprecated: renamed to AssertNotZeroValue in the course of
	// moving to generics to better reflect the functionality.
	return AssertNotZeroValue(t, val)
}

func AssertNil[N comparable](t *testing.T, val N) bool {
	// Deprecated: renamed to AssertZeroValue in the course of
	// moving to generics to better reflect the functionality.
	return AssertZeroValue(t, val)
}

// more specifically, checks that the value is zero value of its type,
// e.g. nil or 0 for int, "" for string, false for bool, etc.
func AssertZeroValue[N comparable](t *testing.T, val N) bool {
	t.Helper()
	var zero N
	if val != zero {
		t.Errorf("(%s)\n    not nil: >%#v<\n", t.Name(), val)
		return false
	}
	return true
}

type Number interface {
	float32 | float64 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

func notWithin[N Number](is, should, eps N) bool {
	return is < should-eps || is > should+eps
}

func AssertWithin[N Number](t *testing.T, is, should, eps N) bool {
	t.Helper()
	if notWithin(is, should, eps) {
		t.Errorf("(%s)\n    is: >%#v<\nshould: >%#v<  ε=%#v", t.Name(), is, should, eps)
		return false
	}
	return true
}

func AssertArrayWithin[N Number](t *testing.T, is, should []N, eps N) bool {
	t.Helper()
	if is == nil && should == nil {
		return true
	}
	if is == nil {
		t.Errorf("arr is nil, should: %v", should)
		return false
	} else if should == nil {
		t.Errorf("should be nil, is: %v", is)
		return false
	} else {
		if len(is) != len(should) {
			t.Errorf("unequal length: %d (%d)", len(is), len(should))
			return false
		}
	}
	for i, _ := range is {
		//if is[i] != should[i] {
		if notWithin(is[i], should[i], eps) {
			t.Errorf("element #%v is:%#v should:%#v ε=%#v", i, is[i], should[i], eps)
			return false
		}
	}
	return true

}
