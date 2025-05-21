package goutils

import (
	"flag"
	"testing"
)

func TestMultiFlag(t *testing.T) {
	args := []string{"-D", "one", "-D", "two", "-D", "three"}
	// Create a new flag set
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	// Create a new multi flag
	var multiFlag MultiFlag
	// Add the multi flag to the flag set
	flagSet.Var(&multiFlag, "D", "test")
	// Parse the command line arguments
	err := flagSet.Parse(args)
	if err != nil {
		t.Fatalf("Error parsing flags: %v", err)
	}
	// Check the values of the multi flag
	if len(multiFlag) != 3 {
		t.Fatalf("Expected 3 values, got %d", len(multiFlag))
	}
	AssertArrayEqual(t, []string{"one", "two", "three"}, multiFlag)

	args = append(args, "-D")
	err = flagSet.Parse(args)
	if err == nil {
		t.Fatalf("Expected error parsing flags, got nil %v", err)
	}
}

func TestMultiDefinitions(t *testing.T) {
	args := []string{"-D", "one=1", "-D", "two=2", "-D", "three=3"}
	// Create a new flag set
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	// Create a new MultiDefinition flag
	var definitions MultiDefinition
	flagSet.Var(&definitions, "D", "test")
	err := flagSet.Parse(args)
	if err != nil {
		t.Fatalf("Error parsing flags: %v", err)
	}
	AssertEqual(t, len(definitions), 3)
	AssertEqual(t, definitions["one"], "1")
	AssertEqual(t, definitions["two"], "2")
	AssertEqual(t, definitions["three"], "3")

	args2 := append(args, "-D", "four")
	err = flagSet.Parse(args2)
	if err == nil {
		t.Fatalf("Expected error parsing flags, got nil %v", err)
	}

	args2 = append(args, "-D")
	err = flagSet.Parse(args2)
	if err == nil {
		t.Fatalf("Expected error parsing flags, got nil %v", err)
	}
}

func TestMultiCount(t *testing.T) {
	args := []string{"-v", "-v", "-v"}
	// Create a new flag set
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	// Create a new MultiDefinition flag
	var count MultiCount
	flagSet.Var(&count, "v", "test")
	err := flagSet.Parse(args)
	if err != nil {
		t.Fatalf("Error parsing flags: %v", err)
	}
	AssertEqual(t, int(count), 3)
}
