package db

import "testing"

import "github.com/a2800276/goutils"

type TestStruct struct {
	One   string
	Two   int
	Three float32
}

func TestCreate(t *testing.T) {
	is := Create(GetInfo(TestStruct{}))
	should := "\nCREATE TABLE IF NOT EXISTS test_struct (\n\t one TEXT\n\t,two INTEGER\n\t,three REAL\n\t,id INTEGER PRIMARY KEY AUTOINCREMENT\n);"
	goutils.AssertEqual(t, is, should)
}
func TestInsert(t *testing.T) {
	is := Insert(GetInfo(TestStruct{}))
	should := "\nINSERT INTO test_struct (\n\t one\n\t,two\n\t,three\n\t,id\n)\nVALUES\n($1, $2, $3, $4)\n"
	goutils.AssertEqual(t, is, should)
}

func TestSave(t *testing.T) {
	is := Save(GetInfo(TestStruct{}))

	should := `
func SaveTestStruct (db *sql.DB, s *TestStruct)(error) {
	sql := INSERT_TestStruct
	res, err := db.Exec(sql, 
		s.One,
		s.Two,
		s.Three,
		s.Id,
	)
	if err != nil {
		return err
	}
	s.Id, err = res.LastInsertId()
	return err
}
`
	goutils.AssertEqual(t, is, should)
}

func TestUnCaml(t *testing.T) {
	var tests = []struct {
		is     string
		should string
	}{
		{"", ""},
		{"AfterBark", "after_bark"},
		{"NotoriousBEAK", "notorious_beak"},
	}

	for _, test := range tests {
		goutils.AssertEqual(t, uncamel(test.is), test.should)
	}
}
