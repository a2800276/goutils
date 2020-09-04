package db

import "database/sql"

type Loader func(i int, row *sql.Rows) error

func LoadBySql(db *sql.DB, sql string, loader Loader) (int, error) {
	i := 0
	rows, err := db.Query(sql)
	if err != nil {
		return i, err
	}
	for rows.Next() {
		if err = loader(i, rows); err != nil {
			return i, err
		}
		i++

	}
	return i, nil
}
