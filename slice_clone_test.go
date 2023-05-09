package main_test

import (
	"testing"

	"golang.org/x/exp/slices"
)

func Clone1(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

func Clone2(src []byte) []byte {
	return append([]byte{}, src...)
}

var Clone3 = slices.Clone[[]byte]

func Clone4(src []byte) []byte {
	dst := make([]byte, 0, len(src))
	dst = append(dst, src...)
	return dst
}

// benchmark compare Clone1, Clone2 and Clone3
// go test -bench=. -benchmem

var byteSlice = []byte("hello world, how are you?")

func BenchmarkClone1(b *testing.B) {
	var res []byte
	for i := 0; i < b.N; i++ {
		res = Clone1(byteSlice)
	}
	if res[0] == 0 {
		b.Log(res)
	}
}

func BenchmarkClone2(b *testing.B) {
	var res []byte
	for i := 0; i < b.N; i++ {
		res = Clone2(byteSlice)
	}
	if res[0] == 0 {
		b.Log(res)
	}
}

func BenchmarkClone3(b *testing.B) {
	var res []byte
	for i := 0; i < b.N; i++ {
		res = slices.Clone(byteSlice)
	}
	if res[0] == 0 {
		b.Log(res)
	}
}

func BenchmarkClone4(b *testing.B) {
	var res []byte
	for i := 0; i < b.N; i++ {
		res = Clone4(byteSlice)
	}
	if res[0] == 0 {
		b.Log(res)
	}
}
