package main

import "fmt"

// Enlarge 将输入的切片变大factor倍
func Enlarge(s []int, factor int) []int {
	ns := make([]int, len(s)*factor)
	copy(ns, s)
	// 让s指向ns
	s = ns
	return s
}

// RemoveSeg 从输入的切片中移除[start, end)的段落
func RemoveSeg(s []string, start, end int) []string {
	ret := make([]string, len(s)-(end-start))
	at := copy(ret, s[:start])
	copy(ret[at:], s[end:])
	return ret
}

func main() {
	s1 := make([]int, 5, 20)
	fmt.Println(len(s1), s1)
	s1 = Enlarge(s1, 2)
	fmt.Println(len(s1), s1)

	s2 := []string{"abc", "def", "ghi", "jkl", "mno"}
	fmt.Println(len(s2), s2)
	s2 = RemoveSeg(s2, 1, 3)
	fmt.Println(len(s2), s2)
}
