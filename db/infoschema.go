package db

import (
	"database/sql"

	util "github.com/a2800276/goutils"
)

type InfoSchema interface {
	Tables() ([]string, error)
	DumpTable(tbl string) ([]col, error)
	FindPrimaryKey(tbl string) ([]col, error)
	Db2GoType(dbType string) string
}

func loadTable(db *sql.DB, tableSql string) ([]string, error) {

	tableFunc := func(stmt *sql.Stmt) (interface{}, error) {
		var tables []string
		rows, err := stmt.Query()
		if err != nil {
			return nil, err
		}
		var table string
		for rows.Next() {
			if err = rows.Scan(&table); err != nil {
				return nil, err
			} else {
				tables = append(tables, table)
			}
		}
		return tables, nil
	}
	tables_, err := util.Execute(db, tableSql, tableFunc)
	if err != nil {
		panic(err)
	}
	tables, ok := tables_.([]string)
	if !ok {
		panic("!?!")
	}
	return tables, nil
}

type col struct {
	name, dataType string
	pk_ordinal     int // odinal position of primary key in case this col is part.
}

func loadColumns(db *sql.DB, tableSql string, tableName string) ([]col, error) {

	tableFunc := func(stmt *sql.Stmt) (interface{}, error) {
		cols := []col{}
		rows, err := stmt.Query(tableName)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var c col
			if err := rows.Scan(&c.name, &c.dataType); err != nil {
				return cols, err
			}
			cols = append(cols, c)
		}
		return cols, nil
	}
	cols_, err := util.Execute(db, tableSql, tableFunc)
	if err != nil {
		panic(err)
	}
	cols, ok := cols_.([]col)
	if !ok {
		panic("!?!")
	}
	return cols, nil
}
