package polymap_test

/*
$ go test -bench . -benchtime 10000000x
goos: darwin
goarch: arm64
pkg: github.com/ldemailly/go-scratch/polymap
goos: darwin
goarch: arm64
pkg: github.com/ldemailly/go-scratch/polymap
BenchmarkStringSetString-11           	10000000	         9.505 ns/op
BenchmarkIntSetInt-11                 	10000000	         5.782 ns/op
BenchmarkAnySetString-11              	10000000	        18.91 ns/op
BenchmarkAnySetInt-11                 	10000000	        14.23 ns/op
BenchmarkObjMapWithStrinStruct-11     	10000000	        46.55 ns/op
BenchmarkAnyMapWithStringStruct-11    	10000000	        46.13 ns/op
BenchmarkAnyMapWithStringPtr-11       	10000000	        13.12 ns/op
BenchmarkObjMapWithIntStruct-11       	10000000	        15.63 ns/op
*/

import (
	"testing"
)

type strMap map[string]int
type intMap map[int]int
type anyMap map[any]int
type objMap map[Object]int // perf wise this is really same as `anyMap`

type Type uint8

type Object interface {
	Type() Type
	// Hash() []byte
}

const (
	INT Type = iota
	STR
)

type Int struct {
	Value int
	//fn    func() // indirect way to confirm deep equality (ie panic for not hashble)
}

func (i Int) Type() Type {
	return INT
}

/*
	func intToBytes(n int) []byte {
		size := unsafe.Sizeof(n)
		b := (*[8]byte)(unsafe.Pointer(&n))[:size:size]
		return b
	}

	func (i *Int) Hash() []byte {
		return intToBytes(i.Value)
	}
*/
type Str struct {
	Value string
	// fn    func() // indirect way to confirm deep equality (ie panic for not hashble)
}

func (s Str) Type() Type {
	return STR
}

/*
	func (s *Str) Hash() []byte {
		return []byte(s.Value)
	}
*/
func BenchmarkStringSetString(b *testing.B) {
	//b.Log("BenchmarkMapSet", b.N)
	m := make(strMap)
	for i := 0; i < b.N; i++ {
		m["foo"]++
	}
	//b.Log(m["foo"])
}

func BenchmarkIntSetInt(b *testing.B) {
	m := make(intMap)
	for i := 0; i < b.N; i++ {
		m[42]++
	}
}

func BenchmarkAnySetString(b *testing.B) {
	m := make(anyMap)
	for i := 0; i < b.N; i++ {
		m["foo"]++
	}
}

func BenchmarkAnySetInt(b *testing.B) {
	m := make(anyMap)
	for i := 0; i < b.N; i++ {
		m[42]++
	}
}

func BenchmarkObjMapWithStrinStruct(b *testing.B) {
	m := make(objMap)
	o := Str{Value: "foo"}
	for i := 0; i < b.N; i++ {
		m[o]++ // lets you pass &o too which changes correctness
	}
}

func BenchmarkAnyMapWithStringStruct(b *testing.B) {
	m := make(anyMap) // it's same as objMap in that case
	o := Str{Value: "foo"}
	for i := 0; i < b.N; i++ {
		m[o]++
	}
}

func BenchmarkAnyMapWithStringPtr(b *testing.B) {
	m := make(anyMap)
	key := "foo"
	for i := 0; i < b.N; i++ {
		m[&key]++ // but this is not correct for equality
	}
}

func BenchmarkObjMapWithIntStruct(b *testing.B) {
	m := make(objMap)
	o := Int{Value: 42}
	for i := 0; i < b.N; i++ {
		m[o]++
	}
}
