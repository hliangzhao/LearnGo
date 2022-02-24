package main

import (
	`fmt`
	`time`
)

// Fib 的DFS实现。这种实现会导致大量的重复计算！
func Fib(x int) int {
	if x < 0 {
		return 0
	} else if x <= 2 {
		return 1
	} else {
		return Fib(x-1) + Fib(x-2)
	}
}

func DPFib(idx int) int {
	// TODO: 基于缓存思想，将已经计算的结果存入切片中
	//  即，DP没有递归，但是有递推思想！
	arr := make([]int, idx+1)
	arr[1] = 1
	arr[2] = 1
	for i := 3; i < idx+1; i++ {
		arr[i] = arr[i-1] + arr[i-2]
	}
	return arr[idx]
}

func main() {
	s1 := time.Now()
	fmt.Println(Fib(44))
	fmt.Println(time.Since(s1))
	s2 := time.Now()
	fmt.Println(DPFib(64))
	fmt.Println(time.Since(s2))
}
