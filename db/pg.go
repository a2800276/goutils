package db

import "database/sql"

type PGInfoSchema struct {
	db *sql.DB
}

func NewPGInfoSchema(connectionString string) (InfoSchema, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return &PGInfoSchema{
		db,
	}, nil
}

// retrieve tables...
func (info *PGInfoSchema) Tables() ([]string, error) {
	sql := `SELECT
		table_name 
		FROM
		information_schema.tables
		WHERE
		table_schema = 'public'`
	return loadTable(info.db, sql)
}

func (info *PGInfoSchema) Db2GoType(typeName string) string {

	switch typeName {
	case "character varying":
		return "string"
	case "timestamp without time zone":
		return "time.Time"
	case "bytea":
		return "[]byte"
	case "bigint":
		return "int64"
	case "integer":
		return "int32"
	default:
		return "!unknown!" // TODO consider err!
	}
}

func (info *PGInfoSchema) DumpTable(tbl string) ([]col, error) {

	sql := `SELECT 
		column_name, 
		data_type 
	FROM 
		information_schema.columns 
	WHERE 
		table_name =$1 
	AND 
		table_schema='public'
	`

	return loadColumns(info.db, sql, tbl)
}

func (info *PGInfoSchema) FindPrimaryKey(tbl string) ([]col, error) {
	sql := `
	SELECT 
		column_name, data_type -- ordinal_position 
	FROM
		information_schema.columns
	JOIN
		information_schema.key_column_usage u
	USING
		(table_schema, table_name, column_name)
	JOIN
		information_schema.table_constraints t
	USING
		(constraint_name,table_schema,table_name)
	WHERE 
		u.table_name = $1
	AND
		u.table_schema = 'public'
	AND
		t.constraint_type='PRIMARY KEY'
	ORDER BY
		u.ordinal_position
	`

	return loadColumns(info.db, sql, tbl)
}
