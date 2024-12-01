package goutils

import (
	"math/rand"
	"strconv"
	"testing"
)

// 123456789|123456789|123456789|123456789|123456789|123456789|123456789|123456789|123456789|123456789
func ruler() {
	for i := 1; i < 100; i++ {
		if i%10 == 0 {
			print("|")
		} else {
			print(i % 10)
		}
	}
	println()
}

func TestHistogram(t *testing.T) {
	h := NewASCIIHist()
	samples := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		0, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		0, 0, 3, 4, 5, 6, 7, 8, 9, 10,
		1, 0, 0, 4, 5, 6, 7, 8, 9, 10,
		1, 0, 0, 0, 5, 6, 7, 8, 9, 10, 11}
	m, _, _, _ := CalculateHistogram(samples, 5)
	//ruler()
	h.DisplayCount = true
	should := `0 (11)|*************************************************************************
2 ( 5)|*********************************
4 ( 9)|***********************************************************
6 (10)|******************************************************************
8 (10)|******************************************************************
10( 6)|***************************************
`

	AssertEqual(t, h.ToString(m), should)
}

func parseFloat32(s string) float32 {
	f32, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	}
	return float32(f32)
}

func TestRandomHisto(t *testing.T) {
	samples := []float32{}
	for i := 0; i < 1000; i++ {
		samples = append(samples, rand.Float32())
	}
	m, _, _, _ := CalculateHistogram(samples, 7)
	ruler()
	h := NewASCIIHist()
	h.DisplayCount = true
	h.CmpFunc = func(a, b string) bool {
		return parseFloat32(a) < parseFloat32(b)
	}

	println(h.ToString(m))
}

func TestBars(t *testing.T) {
	h := NewASCIIHist()
	m := make(map[string]int)
	m["key0"] = 10
	m["key1"] = 5
	m["key2"] = 3
	m["key3"] = 7
	m["key4"] = 2
	m["key5"] = 1
	m["key6"] = 9
	m["key7"] = 4
	m["key8"] = 6
	m["key999"] = 8
	should := `key0  |*************************************************************************
key1  |************************************
key2  |*********************
key3  |***************************************************
key4  |**************
key5  |*******
key6  |*****************************************************************
key7  |*****************************
key8  |*******************************************
key999|**********************************************************
`
	//ruler()
	//println(h.ToString(m))
	AssertEqual(t, h.ToString(m), should)

	h.DisplayCount = true
	//ruler()
	//println(h.ToString(m))

	should = `key0  (10)|*********************************************************************
key1  ( 5)|**********************************
key2  ( 3)|********************
key3  ( 7)|************************************************
key4  ( 2)|*************
key5  ( 1)|******
key6  ( 9)|**************************************************************
key7  ( 4)|***************************
key8  ( 6)|*****************************************
key999( 8)|*******************************************************
`
	AssertEqual(t, h.ToString(m), should)

	h.DisplayCount = false
	h.Width = 15
	//ruler()
	//println(h.ToString(m))
	should = `key0  |********
key1  |****
key2  |**
key3  |*****
key4  |*
key5  |
key6  |*******
key7  |***
key8  |****
key999|******
`
	AssertEqual(t, h.ToString(m), should)

	for k, v := range m {
		m[k] = v * 100
	}
	h.Width = 72
	//ruler()
	//println(h.ToString(m))
	should = `key0  |*****************************************************************
key1  |********************************
key2  |*******************
key3  |*********************************************
key4  |*************
key5  |******
key6  |**********************************************************
key7  |**************************
key8  |***************************************
key999|****************************************************
`
	AssertEqual(t, h.ToString(m), should)

}
