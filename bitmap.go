package goutils

import "fmt"

func BitmapAnnotateUint32(b uint32, annotation map[int]string) string {
	const l = 31
	str := "\n31               |               0\n"
	str += " ----------------|----------------\n"
	var i uint8
	for i = 0; i <= l; i += 1 {
		prefix := (l - i) + 2
		if prefix > 17 {
			prefix += 1
		}
		suffix := i + 1
		if suffix > 16 {
			suffix += 1
		}
		bit := (b & (1 << i)) >> i
		a := annotation[int(i+1)]
		str += fmt.Sprintf("%*d%*s|%2d %s\n", prefix, bit, suffix, "", i+1, a)
	}
	str += " ----------------|----------------\n"
	str += "31               |               0\n"
	return str
}

func BitmapAnnotateUint8(b uint8, annotation map[int]string) string {
	const l = 7
	str := "\n7   |   0\n"
	str += "----|----\n"
	var i uint8
	for i = 0; i <= 7; i += 1 {
		prefix := (l - i) + 1
		if prefix > 4 {
			prefix += 1
		}
		suffix := i
		if suffix > 3 {
			suffix += 1
		}
		bit := (b & (1 << i)) >> i
		a := annotation[int(i+1)]
		str += fmt.Sprintf("%*d%*s|%d %s\n", prefix, bit, suffix, "", i+1, a)
	}
	str += "----|----\n"
	str += "7   |   0"
	return str
}
