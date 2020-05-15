package goutils

import "testing"

func AssertEqual(t *testing.T, is, should interface{}) bool {
	t.Helper()
	if is != should {
		t.Errorf("(%s)\n    is: >%#v<\nshould: >%#v<", t.Name(), is, should)
		return false
	}
	return true
}
