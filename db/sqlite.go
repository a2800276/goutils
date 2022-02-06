package db

import (
	"database/sql"
	"strings"
)

type SqliteInfoSchema struct {
	db *sql.DB
}

func NewSqliteInfoSchema(connectionString string) (InfoSchema, error) {
	// connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return nil, err
	}
	return &SqliteInfoSchema{db}, nil
	//return nil, nil
}

func (info *SqliteInfoSchema) Tables() ([]string, error) {
	sql := `
	SELECT 
		name
	FROM
		sqlite_master -- TODO this will be sqlite_schema
	WHERE
		type = 'table'
	`

	return loadTable(info.db, sql)
}

func (info *SqliteInfoSchema) Db2GoType(typeName string) string {
	// https://www.sqlite.org/datatype3.html
	TYPENAME := strings.ToUpper(typeName)
	switch TYPENAME {
	case "INT":
		fallthrough
	case "INTEGER":
		fallthrough
	case "TINYINT":
		fallthrough
	case "SMALLINT":
		fallthrough
	case "MEDIUMINT":
		fallthrough
	case "BIGINT":
		fallthrough
	case "UNSIGNED BIG INT":
		fallthrough
	case "INT2":
		fallthrough
	case "INT8":
		return "int64" // TODO optimize obv. smaller datatypes
	case "NUMERIC":
		fallthrough
	case "DECIMAL":
		return "big.Int" // TODO for now, compile fails, add import manually
	case "CHARACTER":
		fallthrough
	case "VARCHAR":
		fallthrough
	case "VARYING CHARACTER":
		fallthrough
	case "NCHAR":
		fallthrough
	case "NATIVE CHARACTER":
		fallthrough
	case "NVARCHAR":
		fallthrough
	case "TEXT":
		fallthrough
	case "CLOB":
		return "string"
	case "BLOB":
		return "[]byte"

	case "REAL":
		fallthrough
	case "DOUBLE":
		fallthrough
	case "DOUBLE PRECISION":
		fallthrough
	case "FLOAT":
		return "double"
	case "BOOLEAN":
		return "boolean"
	case "DATE":
		fallthrough
	case "DATETIME":
		return "time.Time"
	default:
		return "!unknown!"
	}
}
func (info *SqliteInfoSchema) DumpTable(tbl string) ([]col, error) {

	// sqlite> .headers on
	// sqlite> PRAGMA TABLE_INFO(sqlite_master);
	// cid|name|type|notnull|dflt_value|pk
	// 0|type|text|0||0
	// 1|name|text|0||0
	// 2|tbl_name|text|0||0
	// 3|rootpage|int|0||0
	// 4|sql|text|0||0

	pragma := `
	SELECT
		name, type
	FROM
		pragma_table_info(:TABLE)
	`
	return loadColumns(info.db, pragma, tbl)
}

func (info *SqliteInfoSchema) FindPrimaryKey(tbl string) ([]col, error) {
	pragma := `
		SELECT
			name, type
		FROM
			pragma_table_info(:TABLE)
		WHERE 
			pk > 0
		ORDER BY
			pk
	`
	return loadColumns(info.db, pragma, tbl)
}
