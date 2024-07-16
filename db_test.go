package goutils

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func getTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?cache=private")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1) // <- this is necessary to avoid some extremly tedious racy conditions
	// which may or may not be a driver bug for in-memory dbs during testing.

	return db
}

func TestExecuteInsert(t *testing.T) {
	sqls := `CREATE TABLE bla (i INTEGER);`
	db := getTestDB(t)
	defer db.Close()
	exec := func(stmt *sql.Stmt) (sql.Result, error) {
		return stmt.Exec()
	}
	_, err := Execute(db, sqls, exec)
	AssertNil(t, err)

	insert := `INSERT INTO bla VALUES (:value);`
	exec = func(stmt *sql.Stmt) (sql.Result, error) {
		return stmt.Exec(5)
	}
	result, err := Insert(db, insert, exec)
	AssertNil(t, err)
	Assert(t, int(result) == 1)

	selects := `SELECT * FROM bla`
	rows, err := db.Query(selects)
	cnt := 0
	AssertNil(t, err)
	for rows.Next() {
		Assert(t, cnt == 0)
		cnt += 1
		var val int
		err := rows.Scan(&val)
		AssertNil(t, err)
		AssertEqual(t, val, 5)
	}
}

type testStruct struct {
	Value int
}

func TestExecute(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	stmt := `select 5;`
	exec := func(stmt *sql.Stmt) (testStruct, error) {
		test := testStruct{}
		err := stmt.QueryRow().Scan(&test.Value)
		return test, err
	}

	testStruct, err := Execute(db, stmt, exec)
	AssertNil(t, err)
	AssertEqual(t, testStruct.Value, 5)
}
