package p2

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/ldemailly/go-scratch/reflect/p1"
)

func Examine(p1 *p1.P1) {
	v := reflect.ValueOf(p1).Elem()
	t := v.Type()
	fmt.Printf("Type: %s\n", t)
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fmt.Printf("Field: %s, Value (printf): %v, CanInterface: %t, CanSet: %t, Value (user code): ", field.Name, value, value.CanInterface(), value.CanSet())
		// getting value.Interface() will panic for unexported fields, even though printf prints it.
		// to get the actual value:
		ptr := unsafe.Pointer(value.UnsafeAddr())
		rf := reflect.NewAt(value.Type(), ptr).Elem()
		fmt.Println(rf.Interface())
	}
}
