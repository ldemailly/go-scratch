package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX 24  // demos safe truncation

// Define a struct to hold the int and the error message
typedef struct {
    int myOtherResult;
    char myError[MAX];
} Result;

// Function that returns the struct with dynamic memory allocation
Result getResult() {
    Result res;
    res.myOtherResult = 42;
    int n = snprintf(res.myError, sizeof(res.myError), "This is a C error %d", 1234567890);
    if (n >= sizeof(res.myError)) {
        fprintf(stderr, "C Warning: snprintf was truncated %d\n", n);
    }
    return res;
}
#cgo CFLAGS: -Wno-unknown-warning-option -Wno-format-truncation
*/
import "C"
import (
	"errors"
	"fmt"
)

func getResult() (int, error) {
	// Call the C function
	res := C.getResult()

	// Access the int value
	code := int(res.myOtherResult)

	// Convert the C string to a Go string
	eMsg := C.GoString(&res.myError[0])
	if eMsg == "" {
		return code, nil
	}
	return code, errors.New(eMsg)
}

func main() {
	i, err := getResult()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Result:", i)
}
