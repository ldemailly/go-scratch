package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX 256  // Define a maximum buffer size

// Define a struct to hold the int and the error message
typedef struct {
    int myOtherResult;
    char myError[MAX];
} Result;

// Function that returns the struct with dynamic memory allocation
Result getResult() {
    Result res;
    res.myOtherResult = 42;
    snprintf(res.myError, sizeof(res.myError), "This is a C error %d", 23);
    return res;
}
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
