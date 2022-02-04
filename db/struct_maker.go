package db

import (
	"database/sql"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	_ "github.com/lib/pq" // shouldn't care but we do, currently only pg
)

import util "github.com/a2800276/goutils"

type InfoSchema interface {
	Tables() ([]string, error)
	DumpTable(tbl string) ([]col, error)
	FindPrimaryKey(tbl string) ([]col, error)
}

type StructMaker struct {
	InfoSchema
}

type PGInfoSchema struct {
	db *sql.DB
}

type SqliteInfoSchema struct {
	db *sql.DB
}

func NewPGInfoSchema(connectionString string) (InfoSchema, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return &PGInfoSchema{db}, nil
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

func camelCase(name string) string {
	if len(name) == 0 {
		return name
	}
	var b strings.Builder
	nextUpper := false
	for i, r := range name {
		if i == 0 {
			b.WriteRune(unicode.ToUpper(r))
			continue
		}
		if r == '_' {
			nextUpper = true
		} else {
			if nextUpper {
				b.WriteRune(unicode.ToUpper(r))
				nextUpper = false
			} else {
				b.WriteRune(r)
			}
		}
	}
	return b.String()
}

func mapType(typeName string) string {
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
		return "!unknown!"
	}
}

func makeStructForTable(tbl string, cols []col) string {
	str := fmt.Sprintf("type %s struct {\n", camelCase(tbl))
	for _, c := range cols {
		fieldName := camelCase(c.name)
		typeName := mapType(c.dataType)
		str = fmt.Sprintf("%s\t%s %s\n", str, fieldName, typeName)
	}

	return fmt.Sprintf("%s}", str)

}

func (sm *StructMaker) MakeStructForTable(tbl string) (string, error) {
	cols, err := sm.DumpTable(tbl)
	if err != nil {
		return "", err
	}
	return makeStructForTable(tbl, cols), nil
}

const loadTemplate = `


func load{{.StructName}}FromRows(rows *sql.Rows)(*{{.StructName}}, error) {
	var s {{.StructName}}
	if err := rows.Scan( {{range .Fields}}
		&s.{{.}}, {{end}} 
	); err != nil {
		return &s, err
	}
	return &s, nil
}

func Load{{.StructName}}By{{.PkName}}(db *sql.DB, pk {{.PkType}}) (*{{.StructName}}, error) {
	sql:="SELECT {{.DBFields}} FROM {{.TableName}} WHERE {{.PkName}} = $1"
	var s *{{.StructName}}
	rows, err := db.Query(sql, pk)
	if err != nil {
		return s, err
	}
	defer rows.Close()
	if rows.Next() {
		if s, err = load{{.StructName}}FromRows(rows); err != nil {
			return s, err
		}
	}
	if rows.Next() {
		return s, fmt.Errorf("more than 1 row ...")
	}
	return s, nil
}

type {{.StructName}}Loader func (i int, s *{{.StructName}}) error;

func BulkLoad{{.StructName}}BySql(db *sql.DB, loader {{.StructName}}Loader, sql string, params ...interface{}) (int, error) {
	i := 0
	rows, err := db.Query(sql, params...)
	if err != nil {
		return i, err
	}
	defer rows.Close()
	for rows.Next() {
		if s, err := load{{.StructName}}FromRows(rows); err != nil {
			return i, err
		} else {
			if err = loader(i, s); err != nil {
				return i, err
			}
			i++
		}
	}
	return i, err

}

func BulkLoad{{.StructName}}(db *sql.DB, limit, offset uint)([]*{{.StructName}}, error) {
	sql:="SELECT {{.DBFields}} FROM {{.TableName}} ORDER BY {{.PkName}} LIMIT $1 OFFSET $2"
	var records []*{{.StructName}}


	loader := func(i int, s *{{.StructName}}) error {
		records = append(records, s)
		return nil
	}

	_, err := BulkLoad{{.StructName}}BySql(db, loader, sql, limit, offset)

	return records, err
}

`

type structTmplData struct {
	TableName  string
	StructName string
	PkName     string
	PkType     string
	DBFields   string // joined by ",", not very elegant, but better than cluttering template(?)
	Fields     []string
}

func (sm *StructMaker) MakeLoadForTable(tbl string) (string, error) {
	cols, err := sm.DumpTable(tbl)
	if err != nil {
		return "", err
	}
	//str := makeStructForTable(tbl, cols), nil

	fields := []string{}
	dbFields := []string{}
	for _, f := range cols {
		fields = append(fields, camelCase(f.name))
		dbFields = append(dbFields, f.name)
	}

	cs, err := sm.FindPrimaryKey(tbl)
	if err != nil {
		return "", err
	}
	if len(cs) != 1 {
		return "", fmt.Errorf("can only handle single primary key, sorry #plsfix")
	}
	c := cs[0]

	s := structTmplData{
		TableName:  tbl,
		StructName: camelCase(tbl),
		PkName:     camelCase(c.name),
		PkType:     mapType(c.dataType),
		DBFields:   strings.Join(dbFields, ", "),
		Fields:     fields,
	}

	builder := strings.Builder{}
	var t = template.Must(template.New("load").Parse(loadTemplate))
	t.Execute(&builder, s)

	return builder.String(), nil

}
