package db

import (
	"fmt"
	"strings"
	"text/template"
	"unicode"

	_ "github.com/lib/pq" // shouldn't care but we do, currently only pg
)

type StructMaker struct {
	InfoSchema
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

func makeStructForTable(tbl string, cols []col, mapper func(string) string) string {
	str := fmt.Sprintf("type %s struct {\n", camelCase(tbl))
	for _, c := range cols {
		fieldName := camelCase(c.name)
		typeName := mapper(c.dataType)
		str = fmt.Sprintf("%s\t%s %s\n", str, fieldName, typeName)
	}

	return fmt.Sprintf("%s}", str)

}

func (sm *StructMaker) MakeStructForTable(tbl string) (string, error) {
	cols, err := sm.DumpTable(tbl)
	if err != nil {
		return "", err
	}
	return makeStructForTable(tbl, cols, sm.Db2GoType), nil
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
		PkType:     sm.InfoSchema.Db2GoType(c.dataType),
		DBFields:   strings.Join(dbFields, ", "),
		Fields:     fields,
	}

	builder := strings.Builder{}
	var t = template.Must(template.New("load").Parse(loadTemplate))
	t.Execute(&builder, s)

	return builder.String(), nil

}
