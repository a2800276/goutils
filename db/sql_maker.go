package db

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
	"unicode"
)

// attempt to convert a go struct into a db table

type Field struct {
	Name string
	Type string
}

type Info struct {
	StructName  string
	TableName   string
	Fields      []Field
	Cols        []Field
	Placeholder []string
}

func uncamel(identifier string) string {
	builder := strings.Builder{}
	prevUpper := false // keep track if previous rune was captilized as not to munge acronyms
	for i, r := range identifier {
		if i == 0 {
			if unicode.IsLower(r) {
				panic(fmt.Sprintf("only public stuff for now, sorry! (%s)", identifier))
			} else {
				builder.WriteRune(unicode.ToLower(r))
				prevUpper = true
			}
		} else if unicode.IsUpper(r) {
			if !prevUpper {
				builder.WriteRune('_')
				prevUpper = true
			}
			builder.WriteRune(unicode.ToLower(r))
		} else {
			prevUpper = false
			builder.WriteRune(r)
		}

	}
	return builder.String()
}

func mapG2S(goType string) string {
	println(goType)
	switch {
	case strings.Index(goType, "int") != -1:
		return "INTEGER"
	case strings.Index(goType, "float") != -1:
		return "REAL"
	case goType == "string":
		return "TEXT"
	default:
		return "!UNKNOWN!"
	}

}

var PK = Field{"id", "INTEGER PRIMARY KEY AUTOINCREMENT"}

func GetInfo(i interface{}) Info {
	t := reflect.TypeOf(i)
	info := Info{}
	info.StructName = t.Name()
	info.TableName = uncamel(info.StructName)

	havePK := false
	placeholder := 1
	for i := 0; i != t.NumField(); i++ {
		f := t.Field(i)

		var col = Field{}
		if f.Name == "Id" {
			havePK = true
			col = PK
		} else {
			info.Fields = append(info.Fields, Field{f.Name, f.Type.Name()})
			col = Field{uncamel(f.Name), mapG2S(f.Type.Kind().String())}
			info.Placeholder = append(info.Placeholder, fmt.Sprintf("$%d", placeholder)) // not addition in templates
			placeholder++
		}
		info.Cols = append(info.Cols, col)
	}
	if !havePK {
		info.Fields = append(info.Fields, Field{"Id", "int64"})
		info.Cols = append(info.Cols, PK)
	}
	return info
}

func fillTemplate(tName string, tmplt string, info Info) string {
	t := template.Must(template.New(tName).Parse(tmplt))
	builder := strings.Builder{}
	t.Execute(&builder, info)

	return builder.String()
}

// pass in a struct
func Create(info Info) string {
	return fillTemplate("create", CREATE_TEMPLATE, info)
}

func Insert(info Info) string {
	return fillTemplate("insert", INSERT_TEMPLATE, info)

}

func Save(info Info) string {
	return fillTemplate("save", SAVE_TEMPLATE, info)
}

const CREATE_TEMPLATE = `
CREATE TABLE IF NOT EXISTS {{.TableName}} ({{range $idx, $f := .Cols}}
	{{if $idx}},{{else}} {{end}}{{.Name}} {{.Type}}{{end}}
);`

const INSERT_TEMPLATE = `
INSERT INTO {{.TableName}} ({{range $idx, $f := .Cols}}
	{{if $idx}},{{else}} {{end}}{{.Name}}{{end}}
)
VALUES
({{range $idx, $f := .Placeholder}}{{if $idx}}, {{else}}{{end}}{{.}}{{end}})
`

const SAVE_TEMPLATE = `
func Save{{.StructName}} (db *sql.DB, s *{{.StructName}})(error) {
	sql := INSERT_{{.StructName}}
	res, err := db.Exec(sql, {{range .Fields}}
		s.{{.Name}},{{end}}
	)
	if err != nil {
		return err
	}
	s.Id, err = res.LastInsertId()
	return err
}
`
