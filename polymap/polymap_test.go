package polymap_test

/*
$ go test -bench . -benchtime 10000000x
goos: darwin
goarch: arm64
pkg: github.com/ldemailly/go-scratch/polymap
BenchmarkMapSetStr-11              	10000000	         8.908 ns/op
BenchmarkMapSetInt-11              	10000000	         5.742 ns/op
BenchmarkMapSetAnyWithString-11    	10000000	        20.57 ns/op
BenchmarkMapSetAnyWithInt-11       	10000000	        14.20 ns/op
BenchmarkObjMapWithString-11       	10000000	        13.41 ns/op
BenchmarkObjMapWithInt-11          	10000000	        13.44 ns/op
*/

import (
	"testing"
)

type strMap map[string]int
type intMap map[int]int
type anyMap map[any]int
type objMap map[Object]int

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
}

func (i *Int) Type() Type {
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
}

func (s *Str) Type() Type {
	return STR
}

/*
	func (s *Str) Hash() []byte {
		return []byte(s.Value)
	}
*/
func BenchmarkMapSetStr(b *testing.B) {
	//b.Log("BenchmarkMapSet", b.N)
	m := make(strMap)
	for i := 0; i < b.N; i++ {
		m["foo"]++
	}
	//b.Log(m["foo"])
}

func BenchmarkMapSetInt(b *testing.B) {
	//b.Log("BenchmarkMapSet", b.N)
	m := make(intMap)
	for i := 0; i < b.N; i++ {
		m[42]++
	}
	//b.Log(m["foo"])
}

func BenchmarkMapSetAnyWithString(b *testing.B) {
	//b.Log("BenchmarkMapSet", b.N)
	m := make(anyMap)
	// m[42] = 42
	for i := 0; i < b.N; i++ {
		m["foo"]++
	}
	//b.Log(m["foo"])
}

func BenchmarkMapSetAnyWithInt(b *testing.B) {
	m := make(anyMap)
	for i := 0; i < b.N; i++ {
		m[42]++
	}
}

func BenchmarkObjMapWithString(b *testing.B) {
	//b.Log("BenchmarkMapSet", b.N)
	m := make(objMap)
	o := &Str{Value: "foo"}
	// m[42] = 42
	for i := 0; i < b.N; i++ {
		m[o]++
	}
	//b.Log(m["foo"])
}

func BenchmarkObjMapWithInt(b *testing.B) {
	m := make(objMap)
	o := &Int{Value: 42}
	for i := 0; i < b.N; i++ {
		m[o]++
	}
}
