package main

import (
	"fmt"

	"github.com/a2800276/goutils"
)

func main() {

	fmt.Println("System Info:")
	si := goutils.NewSystemInfo()
	fmt.Println(si)
}
