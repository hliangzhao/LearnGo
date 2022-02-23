package main

import (
	"fmt"
	"testing"
)

// BenchmarkMyConcurrentAtomicAdd 基准测试主要测试函数的运行时长
func BenchmarkMyConcurrentAtomicAdd(b *testing.B) {
	// TODO：首先重制计时器，然后运行多组取平均值
	b.ResetTimer()
	fmt.Println(b.N)
	// 运行N次取平均值
	for i := 0; i < b.N; i++ {
		MyConcurrentAtomicAdd()
	}
}

func BenchmarkMyConcurrentMutexAdd(b *testing.B) {
	b.ResetTimer()
	fmt.Println(b.N)
	for i := 0; i < b.N; i++ {
		MyConcurrentMutexAdd()
	}
}

// (base) ➜  concurrent go test -v -run="none" -bench=.
// goos: darwin
// goarch: amd64
// pkg: 7-testing-benchmarking/concurrent
// BenchmarkMyConcurrentAtomicAdd
// 1
// 100
// 4498
// 5397
// BenchmarkMyConcurrentAtomicAdd-8            5397            216217 ns/op
// BenchmarkMyConcurrentMutexAdd
// 1
// 100
// 3552
// BenchmarkMyConcurrentMutexAdd-8             3552            332095 ns/op
// PASS
// ok      7-testing-benchmarking/concurrent       3.796s
