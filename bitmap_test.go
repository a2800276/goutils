package goutils

import "testing"

func TestAnnotate(t *testing.T) {
	b := uint8(0xf0)
	annotation := map[uint64]string{
		1: "BIT ONE",
		4: "BIT Four",
		6: "Surprise!",
		8: "Whatever",
	}
	is := BitmapAnnotateUint8(b, nil)
	should := "\n7   |   0\n----|----\n        0|1 \n       0 |2 \n      0  |3 \n     0   |4 \n   1     |5 \n  1      |6 \n 1       |7 \n1        |8 \n----|----\n7   |   0"

	AssertEqual(t, is, should)

	is = BitmapAnnotateUint8(b, annotation)
	should = "\n7   |   0\n----|----\n        0|1 BIT ONE\n       0 |2 \n      0  |3 \n     0   |4 BIT Four\n   1     |5 \n  1      |6 Surprise!\n 1       |7 \n1        |8 Whatever\n----|----\n7   |   0"
	AssertEqual(t, is, should)

	c := uint32(0xF0A15AA5)
	is = BitmapAnnotateUint32(c, nil)
	should = "\n31               |               0\n ----------------|----------------\n                                 1 | 1 \n                                0  | 2 \n                               1   | 3 \n                              0    | 4 \n                             0     | 5 \n                            1      | 6 \n                           0       | 7 \n                          1        | 8 \n                         0         | 9 \n                        1          |10 \n                       0           |11 \n                      1            |12 \n                     1             |13 \n                    0              |14 \n                   1               |15 \n                  0                |16 \n                1                  |17 \n               0                   |18 \n              0                    |19 \n             0                     |20 \n            0                      |21 \n           1                       |22 \n          0                        |23 \n         1                         |24 \n        0                          |25 \n       0                           |26 \n      0                            |27 \n     0                             |28 \n    1                              |29 \n   1                               |30 \n  1                                |31 \n 1                                 |32 \n ----------------|----------------\n31               |               0\n"

	AssertEqual(t, is, should)

	is = BitmapAnnotateUint32(c, annotation)
	should = "\n31               |               0\n ----------------|----------------\n                                 1 | 1 BIT ONE\n                                0  | 2 \n                               1   | 3 \n                              0    | 4 BIT Four\n                             0     | 5 \n                            1      | 6 Surprise!\n                           0       | 7 \n                          1        | 8 Whatever\n                         0         | 9 \n                        1          |10 \n                       0           |11 \n                      1            |12 \n                     1             |13 \n                    0              |14 \n                   1               |15 \n                  0                |16 \n                1                  |17 \n               0                   |18 \n              0                    |19 \n             0                     |20 \n            0                      |21 \n           1                       |22 \n          0                        |23 \n         1                         |24 \n        0                          |25 \n       0                           |26 \n      0                            |27 \n     0                             |28 \n    1                              |29 \n   1                               |30 \n  1                                |31 \n 1                                 |32 \n ----------------|----------------\n31               |               0\n"

	AssertEqual(t, is, should)
	//t.Logf("%#v", is)
	//t.Fail()
}
