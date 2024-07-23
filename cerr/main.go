package main

/*
#include <stdlib.h>

// Define a struct to hold the int and the error message
typedef struct {
    int myOtherResult;
    const char *myError;
} Result;

// Function that returns the struct
Result getResult() {
    Result res;
    res.myOtherResult = 42;
    res.myError = "This is an error message from C";
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
	eMsg := C.GoString(res.myError)
	if eMsg == "" {
		return code, nil
	}
	return code, errors.New(eMsg)
}

func main() {
	i, err := getResult()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Result:", i)
}
