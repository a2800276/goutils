package db

import (
	"database/sql"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	_ "github.com/lib/pq"
)

type StructMaker struct {
	db *sql.DB
}

func NewStructMaker(user, password, host, port, dbname string) (sm StructMaker, err error) {
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return
	}
	sm.db = db
	return
}

// retrieve tables...
func (sm *StructMaker) Tables() ([]string, error) {
	sql := `SELECT
		table_name 
		FROM
		information_schema.tables
		WHERE
		table_schema = 'public'`
	tables := []string{}
	rows, err := sm.db.Query(sql)
	if err != nil {
		return tables, err
	}
	defer rows.Close()
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return tables, err
		}

		tables = append(tables, table)
	}
	return tables, nil
}

type col struct {
	name, dataType string
}

func (sm *StructMaker) DumpTable(tbl string) ([]col, error) {
	sql := `SELECT 
		column_name, 
		data_type 
	FROM 
		information_schema.columns 
	WHERE 
		table_name =$1 
	AND 
		table_schema='public'`
	cols := []col{}
	rows, err := sm.db.Query(sql, tbl)
	if err != nil {
		return cols, err
	}
	defer rows.Close()
	for rows.Next() {
		var c col
		if err := rows.Scan(&c.name, &c.dataType); err != nil {
			return cols, err
		}
		cols = append(cols, c)
	}
	return cols, nil

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
func Load{{.StructName}}By{{.PkName}}(db *sql.DB, pk {{.PkType}}) ({{.StructName}}, error) {
	sql:="SELECT * FROM {{.TableName}} WHERE {{.PkName}} = $1"
	var s {{.StructName}}
	rows, err := db.Query(sql, pk)
	if err != nil {
		return s, err
	}
	defer rows.Close()
	if rows.Next() {
		// template for load all ...
		if err := rows.Scan( {{range .Fields}}
			&s.{{.}}, {{end}} 
		); err != nil {
			return err
		}
	}
	if rows.Next() {
		return fmt.Errorf("more than 1 row ...")
	}
	return s, nil
`

type structTmplData struct {
	TableName  string
	StructName string
	PkName     string
	PkType     string
	Fields     []string
}

func (sm *StructMaker) MakeLoadForTable(tbl string) (string, error) {
	cols, err := sm.DumpTable(tbl)
	if err != nil {
		return "", err
	}
	//str := makeStructForTable(tbl, cols), nil

	fields := []string{}
	for _, f := range cols {
		fields = append(fields, camelCase(f.name))
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
		Fields:     fields,
	}

	builder := strings.Builder{}
	var t = template.Must(template.New("load").Parse(loadTemplate))
	t.Execute(&builder, s)

	return builder.String(), nil

}

func (sm *StructMaker) FindPrimaryKey(tbl string) ([]col, error) {
	sql := `SELECT 
		column_name --, ordinal_position 
	FROM 
		information_schema.key_column_usage u
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
		ordinal_position
	`

	var cols []col
	rows, err := sm.db.Query(sql, tbl)
	if err != nil {
		return cols, err
	}
	defer rows.Close()
	for rows.Next() {
		var c col
		if err := rows.Scan(&c.name); err != nil {
			return cols, err
		}
		cols = append(cols, c)
	}
	return cols, nil

}
