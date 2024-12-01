package goutils

import (
	"fmt"
	"math"
	"slices"
	"sort"
)

type ComparisonFunc func(s1, s2 string) bool

type ASCIIHist struct {
	Width        int  // the maximum width of the histogram / console
	DisplayCount bool // display the count/value of the key
	CmpFunc      ComparisonFunc
	OrderedKeys  []string // provide keys in desired order or leave empty for default order (alphabetically, taking into account numerics, e.g. key2 before key10)
}

func NewASCIIHist() *ASCIIHist {
	return &ASCIIHist{Width: 80}
}

func (h *ASCIIHist) keys(m map[string]int) []string {
	if h.OrderedKeys != nil && len(h.OrderedKeys) == len(m) {
		return h.OrderedKeys
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	if h.CmpFunc == nil {
		h.CmpFunc = CompareMixedNumericStrings
	}

	sort.Slice(keys, func(i, j int) bool {
		return h.CmpFunc(keys[i], keys[j])
	})
	return keys
}

// This function writes a precalcualted histogram to a string, i.e.:
// each key is a bucket, the value is the count of the bucket.
// This is more accurately a bar chart. :)
func (h *ASCIIHist) ToString(m map[string]int) string {
	// determine the max of the values
	max := 0
	// determine the width of the key column
	width := 0
	for k, v := range m {
		if v > max {
			max = v
		}
		if len(k) > width {
			width = len(k)
		}
	}
	// histo should be:
	// key0     |********
	// otherkey1|*****

	// determine the width of the value column
	valueWidth := h.Width - width - 1
	// determine scaling factor
	if h.DisplayCount {
		valueWidth -= (2 + len(fmt.Sprintf("%d", max)))
	}
	scale := float64(valueWidth) / float64(max)

	str := ""
	for _, k := range h.keys(m) {
		v := m[k]
		// print key
		if h.DisplayCount {
			//			dw := width + (1 + len(fmt.Sprintf("%d", max)))
			str += fmt.Sprintf("%-*s(%*d)|", width, k, len(fmt.Sprintf("%d", max)), v)
		} else {
			str += fmt.Sprintf("%-*s|", width, k)
		}
		// print value
		for i := 0; i < int(float64(v)*scale); i++ {
			str += ("*")
		}
		str += "\n"
	}
	return str
}

func CalculateHistogram[T Number](samples []T, numBuckets int) (map[string]int, T, T, float64) {
	// This function sorts through the samples, divides them into buckets
	// and returns a map of the buckets and their counts and basic min/max/avg/stdev statistics.
	// The number of buckets is determined by the ASCIIHist.Buckets field.

	// sort the samples
	slices.Sort(samples)
	// determine the min/max/avg/stdev
	min := samples[0]
	max := samples[len(samples)-1]
	println(max)
	avg := 0.0
	for _, v := range samples {
		avg += float64(v)
	}
	avg /= float64(len(samples))
	// determine the stdev
	stdev := 0.0
	for _, v := range samples {
		stdev += (float64(v) - avg) * (float64(v) - avg)
	}
	stdev = math.Sqrt(stdev / float64(len(samples)))
	// determine the number of buckets
	if numBuckets == 0 {
		numBuckets = 10
	}
	// determine the width of the key column
	width := 0
	for _, v := range samples {
		if len(fmt.Sprintf("%v", v)) > width {
			width = len(fmt.Sprintf("%v", v))
		}
	}
	// determine the buckets
	bucketWidth := (max - min) / T(numBuckets)
	//buckets := make(map[string]int)
	buckets := make([]int, numBuckets)
	//for i := 0; i < buckets; i++ {
	//	bucket["0"] = 0
	//}
	for i, v := range samples {
		b := (v - min) / bucketWidth
		//buckets[fmt.Sprintf("%.3f", float64(b*bucketWidth))]++
		fmt.Printf("%d v: %v, b: %v,  w:%v, bucket: %v\n", i, v, b, bucketWidth, int(math.Floor(float64(b))))
		buckets[int(math.Floor(float64(b)))]++

	}
	bucketsMap := make(map[string]int)
	for i, v := range buckets {
		bucketsMap[fmt.Sprintf("%v", float64(i)*float64(bucketWidth))] = v
	}
	return bucketsMap, min, max, avg // , stdev
}
