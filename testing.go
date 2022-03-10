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

func AssertBytesEqual(t *testing.T, is, should []byte) bool {
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
			return false
		}
	}

	return true
}
func AssertNotEqual(t *testing.T, is, shouldnt interface{}) bool {
	t.Helper()
	if is == shouldnt {
		t.Errorf("(%s)\n    is: >%#v<\nshould: >%#v<", t.Name(), is, shouldnt)
		return false
	}
	return true
}

func AssertNotNil(t *testing.T, val interface{}) bool {
	t.Helper()
	if val == nil {
		t.Errorf("(%s)\n    not nil: >%#v<\n", t.Name(), val)
		return false
	}
	return true
}

func AssertNil(t *testing.T, val interface{}) bool {
	t.Helper()
	if val != nil {
		t.Errorf("(%s)\n    not nil: >%#v<\n", t.Name(), val)
		return false
	}
	return true
}
