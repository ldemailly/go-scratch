package main_test

import (
	"testing"
)

func makeTestSlice(n int) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = i
	}
	return s
}

func sum(numbers []int, start int, end int, c chan int) {
	sum := 0
	for i := start; i < end; i++ {
		sum += numbers[i]
	}
	c <- sum
}

func directSum(arr []int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return sum
}

// Sum the numbers in a slice using numberOfSplit go routines
func splitSum(c chan int, arr []int, numberOfSplit int) int {
	splitSize := len(arr) / numberOfSplit
	for i := 0; i < numberOfSplit; i++ {
		start := i * splitSize
		end := (i + 1) * splitSize
		if i == numberOfSplit-1 {
			end = len(arr)
		}
		go sum(arr, start, end, c)
	}
	sum := 0
	for i := 0; i < numberOfSplit; i++ {
		sum += <-c
	}
	return sum
}

const chanCapacity = 24

var c = make(chan int)

func TestSliceSplit(t *testing.T) {
	arr := makeTestSlice(1007)
	sum := splitSum(c, arr, 10)
	t.Log("Sum from go routines: ", sum)
	// check:
	dSum := directSum(arr)
	if dSum != sum {
		t.Errorf("Mismatch in sums: %d %d", dSum, sum)
	}
}

const howMany = 1000023

func BenchmarkDirectSum(b *testing.B) {
	arr := makeTestSlice(howMany)
	for i := 0; i < b.N; i++ {
		directSum(arr)
	}
}

func BenchmarkSplitSlice(b *testing.B) {
	arr := makeTestSlice(howMany)
	for i := 0; i < b.N; i++ {
		splitSum(c, arr, chanCapacity)
	}
}
