package main

import (
	"math"
	"reflect"
	"runtime"
	"strings"
)

func FunctionName(f any) string {
	val := reflect.ValueOf(f)
	if val.Kind() != reflect.Func {
		return ""
	}
	fullName := runtime.FuncForPC(val.Pointer()).Name()
	return fullName
}

func ShortFunctionName(f any) string {
	fullName := FunctionName(f)
	if fullName == "" {
		return ""
	}
	lastDot := strings.LastIndex(fullName, ".")
	if lastDot == -1 {
		return fullName
	}
	return fullName[lastDot+1:]
}

func main() {
	fn := math.Cos               // for instance
	str := ShortFunctionName(fn) // "Cos"
	println(str)
}
