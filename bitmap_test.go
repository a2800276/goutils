package goutils

import "testing"

func TestAnnotate(t *testing.T) {
	b := uint8(0xf0)
	annotation := map[int]string{
		1: "BIT ONE",
		4: "BIT Four",
		6: "Surprise!",
		8: "Whatever",
	}
	is := BitmapAnnotateUint8(b, nil)
	should := "\n7   |   0\n----|----\n        0|0 \n       0 |1 \n      0  |2 \n     0   |3 \n   1     |4 \n  1      |5 \n 1       |6 \n1        |7 \n----|----\n7   |   0"

	AssertEqual(t, is, should)

	is = BitmapAnnotateUint8(b, annotation)
	should = "\n7   |   0\n----|----\n        0|0 \n       0 |1 BIT ONE\n      0  |2 \n     0   |3 \n   1     |4 BIT Four\n  1      |5 \n 1       |6 Surprise!\n1        |7 \n----|----\n7   |   0"
	AssertEqual(t, is, should)

	c := uint32(0xF0A15AA5)
	is = BitmapAnnotateUint32(c, nil)
	should = "\n31               |               0\n ----------------|----------------\n                                 1 | 0 \n                                0  | 1 \n                               1   | 2 \n                              0    | 3 \n                             0     | 4 \n                            1      | 5 \n                           0       | 6 \n                          1        | 7 \n                         0         | 8 \n                        1          | 9 \n                       0           |10 \n                      1            |11 \n                     1             |12 \n                    0              |13 \n                   1               |14 \n                  0                |15 \n                1                  |16 \n               0                   |17 \n              0                    |18 \n             0                     |19 \n            0                      |20 \n           1                       |21 \n          0                        |22 \n         1                         |23 \n        0                          |24 \n       0                           |25 \n      0                            |26 \n     0                             |27 \n    1                              |28 \n   1                               |29 \n  1                                |30 \n 1                                 |31 \n ----------------|----------------\n31               |               0\n"

	AssertEqual(t, is, should)

	is = BitmapAnnotateUint32(c, annotation)
	should = "\n31               |               0\n ----------------|----------------\n                                 1 | 0 \n                                0  | 1 BIT ONE\n                               1   | 2 \n                              0    | 3 \n                             0     | 4 BIT Four\n                            1      | 5 \n                           0       | 6 Surprise!\n                          1        | 7 \n                         0         | 8 Whatever\n                        1          | 9 \n                       0           |10 \n                      1            |11 \n                     1             |12 \n                    0              |13 \n                   1               |14 \n                  0                |15 \n                1                  |16 \n               0                   |17 \n              0                    |18 \n             0                     |19 \n            0                      |20 \n           1                       |21 \n          0                        |22 \n         1                         |23 \n        0                          |24 \n       0                           |25 \n      0                            |26 \n     0                             |27 \n    1                              |28 \n   1                               |29 \n  1                                |30 \n 1                                 |31 \n ----------------|----------------\n31               |               0\n"

	AssertEqual(t, is, should)
	//t.Logf("%#v", is)
	//t.Fail()
}
