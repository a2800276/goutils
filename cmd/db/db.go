package main

import (
	"flag"
	"log"
	"os"

	"github.com/a2800276/goutils/db"
)

var (
	user     = flag.String("user", "", "database user")
	password = flag.String("password", "", "db password")
	host     = flag.String("host", "", "db hostname")
	port     = flag.String("port", "", "port")
	dbname   = flag.String("dbname", "", "dbname")

	table = flag.String("table", "", "table to dump")
)

func main() {
	log.Println("Welcome.")
	flag.Parse()
	sm, err := db.NewStructMaker(*user, *password, *host, *port, *dbname)
	if err != nil {
		log.Printf("%v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	str, err := sm.MakeStructForTable(*table)
	if err != nil {
		log.Printf("%v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	println(str)

	str, err = sm.MakeLoadForTable(*table)
	if err != nil {
		log.Printf("%v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	println(str)
}