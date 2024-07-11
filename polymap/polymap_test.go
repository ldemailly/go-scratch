package polymap_test

import "testing"

type strMap map[string]int
type intMap map[int]int
type anyMap map[any]int

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
