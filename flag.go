package goutils

import (
	"fmt"
	"strings"
)

// MultiFlag is a custom flag type that allows multiple values to be set,
// e.g. -D one -D two -D three =>
//
//	flag.Var(&multiFlag, "D", "define some things")
//	flag.Parse()
//	fmt.Println(multiFlag) // [one two three]
type MultiFlag []string

func (m *MultiFlag) String() string {
	return strings.Join(*m, ",")
}

func (m *MultiFlag) Set(value string) error {
	println("set", value)
	*m = append(*m, value)
	fmt.Printf("%v\n", m)
	return nil
}

// MultiDefinition is a custom flag type that allows multiple key-value pairs to be set,
// e.g. -D one=1 -D two=2 -D three=3 =>
//
//	flag.Var(&multiFlag, "D", "define some things")
//	flag.Parse()
//	fmt.Println(multiFlag) // map[one:1 two:2 three:3]
type MultiDefinition map[string]string

func (m *MultiDefinition) String() string {
	return fmt.Sprintf("%v", *m)
}
func (m *MultiDefinition) Set(value string) error {
	// Split the value into key and value
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format: %s, expected key=value", value)
	}
	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	if *m == nil {
		*m = make(MultiDefinition)
	}
	(*m)[key] = val
	return nil
}

// MultiCount is a custom flag type that counts the number of times it is set, e.g. ssh style verbose flags: -v -v -v =>
//
//	flag.Var(&multiCount, "v", "verbose mode")
//	flag.Parse()
//	fmt.Println(multiCount) // 3
type MultiCount int

func (m *MultiCount) String() string {
	return fmt.Sprintf("%d", *m)
}
func (m *MultiCount) Set(value string) error {
	if value != "true" {
		return fmt.Errorf("invalid format: >%s<, expected empty string", value)
	}

	*m++
	return nil
}

func (m *MultiCount) IsBoolFlag() bool {
	return true
}
