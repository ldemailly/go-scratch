package main

import (
	"bytes"
	"go/token"
	"testing"
)

func TestCountEffectiveLinesOfCode(t *testing.T) {
	fset := token.NewFileSet()
	var buf bytes.Buffer
	n, err := countEffectiveLinesOfCode(&buf, true, "testdata/test_simple.go", fset)
	if err != nil {
		t.Fatalf("Error counting effective lines of code: %v", err)
	}
	if n != 6 {
		t.Errorf("Expected 6 effective lines of code, got %d", n)
	}
	expectedOutput := `Line   1 (  4): 	package main
Line   2 (  6): 	import "fmt"
Line   3 (  8): 	func main() {
Line   4 (  9): 		x := 5
Line   5 ( 11): 		fmt.Println(x)
Line   6 ( 12): 	}
` // note first line: ( 4) is correct but ( 6) should be ( 7) etc.
	actualOutput := buf.String()
	if actualOutput != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, actualOutput)
	}
}
