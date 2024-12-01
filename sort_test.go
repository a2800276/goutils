package goutils

import (
	"sort"
	"testing"
)

func TestCompareMixedNumericStrings(t *testing.T) {
	keys := []string{"a2b3", "a10b2", "a2b10", "b1a1", "a1b2", "a1b10"}

	sort.Slice(keys, func(i, j int) bool {
		return CompareMixedNumericStrings(keys[i], keys[j])
	})

	should := []string{"a1b2", "a1b10", "a2b3", "a2b10", "a10b2", "b1a1"}
	AssertArrayEqual(t, keys, should)

	//keys = []string{"0.27", "0.027", "0.0027"}
	//sort.Slice(keys, func(i, j int) bool {
	//	return CompareMixedNumericStrings(keys[i], keys[j])
	//})
	//should = []string{"0.0027", "0.027", "0.27"}
	//AssertArrayEqual(t, keys, should)

}
