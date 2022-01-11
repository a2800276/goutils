package goutils

import "database/sql"

// SQL Helper functions
// the following functions are intended to cut down on/ centralize
// sql boilerplate code.
// These functions assume:
// - on open db connection
// - that a sql statement is compiled to a PreparedStatement
// - which is either an INSERT or a SELECT

// `execFunction`s carry out the work needed to be done
// with the compiled PreparedStatement, i.e. extract results
// or Scan values into an object.
type ExecFunc func(*sql.Stmt) (interface{}, error)

// Compiles the passed SQL statement to a PreparedStatement which
// is passed off to the provided execFunc and closed once that
// function returns.
func Execute(db *sql.DB, sqls string, exec ExecFunc) (interface{}, error) {
	stmt, err := db.Prepare(sqls)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return exec(stmt)
}

// Same as `execute, but assumes that each INSERT statement is assigned
// an automated primary key which is retrieved via Result.LastInsertId()
func Insert(db *sql.DB, sqls string, exec ExecFunc) (int64, error) {
	result_, err := Execute(db, sqls, exec)

	if err != nil {
		return -1, err
	}
	result := result_.(sql.Result)

	return result.LastInsertId()
}
