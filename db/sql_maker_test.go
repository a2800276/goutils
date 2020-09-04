package db

import "testing"

import "github.com/a2800276/goutils"

type TestStruct struct {
	One   string
	Two   int
	Three float32
}

type Data struct {
	Id                 uint64
	NeoId              int64
	Mac                string
	Timestamp          int64
	MessageType        byte
	MvIdle             int
	MvLoad             int
	Status             byte
	Uptime             uint32
	AdvertisementCount byte
}

func TestCreate(t *testing.T) {
	//is := Create(GetInfo(TestStruct{}))
	is := Create(GetInfo(Data{}))
	should := "\nCREATE TABLE IF NOT EXIST test_struct (\n\n\t one TEXT\n\n\t,two INTEGER\n\n\t,three REAL\n\n);"
	goutils.AssertEqual(t, is, should)
}
func TestInsert(t *testing.T) {
	is := Insert(GetInfo(Data{}))
	should := "\nCREATE TABLE IF NOT EXIST test_struct (\n\n\t one TEXT\n\n\t,two INTEGER\n\n\t,three REAL\n\n);"
	goutils.AssertEqual(t, is, should)
}

func TestSave(t *testing.T) {
	is := Save(GetInfo(Data{}))
	should := "\nCREATE TABLE IF NOT EXIST test_struct (\n\n\t one TEXT\n\n\t,two INTEGER\n\n\t,three REAL\n\n);"
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

type Xfer struct {
	Id        int64
	Raw       []byte
	Mac       string // Beacon in the DB
	Timestamp int64
	Rssi      int32
}

func TestCreateX(t *testing.T) {
	//is := Create(GetInfo(TestStruct{}))
	is := Create(GetInfo(Xfer{}))
	should := "\nCREATE TABLE IF NOT EXIST test_struct (\n\n\t one TEXT\n\n\t,two INTEGER\n\n\t,three REAL\n\n);"
	goutils.AssertEqual(t, is, should)
}
func TestInsertX(t *testing.T) {
	is := Insert(GetInfo(Xfer{}))
	should := "\nCREATE TABLE IF NOT EXIST test_struct (\n\n\t one TEXT\n\n\t,two INTEGER\n\n\t,three REAL\n\n);"
	goutils.AssertEqual(t, is, should)
}

func TestSaveX(t *testing.T) {
	is := Save(GetInfo(Xfer{}))
	should := "\nCREATE TABLE IF NOT EXIST test_struct (\n\n\t one TEXT\n\n\t,two INTEGER\n\n\t,three REAL\n\n);"
	goutils.AssertEqual(t, is, should)
}
