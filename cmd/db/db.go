package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)
import (
	"github.com/a2800276/goutils/db"
)

import (
	_ "github.com/mattn/go-sqlite3"
)

var (
	user     = flag.String("user", "", "pg database user")
	password = flag.String("password", "", "pg db password")
	host     = flag.String("host", "", "pg db hostname")
	port     = flag.String("port", "", "pg port")
	dbname   = flag.String("dbname", "", "pg dbname")
	sqlite   = flag.String("sqlite", "", "sqlite connection string")

	// table = flag.String("table", "", "table to dump")
)

func main() {
	log.Println("Welcome.")
	flag.Parse()

	var infoScheme db.InfoSchema
	var err error

	if *sqlite != "" {
		infoScheme, err = db.NewSqliteInfoSchema(*sqlite)
	} else {
		// if no sqlite string is set, assume pg ...
		connectionString := fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			*user, *password,
			*host, *port, *dbname)
		infoScheme, err = db.NewPGInfoSchema(connectionString)

	}

	if err != nil {
		log.Printf("%v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	structMaker := db.StructMaker{
		infoScheme,
	}
	tables, err := structMaker.InfoSchema.Tables()
	if err != nil {
		log.Printf("%v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	for _, table := range tables {
		str, err := structMaker.MakeStructForTable(table)
		if err != nil {
			log.Printf("%v\n", err)
			flag.Usage()
			os.Exit(1)
		}

		println(str)

		str, err = structMaker.MakeLoadForTable(table)
		if err != nil {
			log.Printf("%v\n", err)
			flag.Usage()
			os.Exit(1)
		}

		println(str)
	}
}
