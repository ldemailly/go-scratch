package main_test

import (
	"testing"

	"golang.org/x/exp/slices"
)

// Just for reference as obviously that's what copy() is for.
// also to "warm up" the benchmarks.
func Clone0(src []byte) []byte {
	dst := make([]byte, 0, len(src))
	for i := range src {
		dst = append(dst, src[i])
	}
	return dst
}

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

func Clone5(src []byte) []byte {
	return append([]byte(nil), src...)
}

// Clearly this is a small slice so only one of many cases...
var byteSlice = []byte("hello world, how are you?")

func BenchmarkClone0(b *testing.B) {
	var res []byte
	for i := 0; i < b.N; i++ {
		res = Clone0(byteSlice)
	}
	if res[0] == 0 {
		b.Log(res)
	}
}

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

func BenchmarkClone5(b *testing.B) {
	var res []byte
	for i := 0; i < b.N; i++ {
		res = Clone5(byteSlice)
	}
	if res[0] == 0 {
		b.Log(res)
	}
}
