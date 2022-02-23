package main

import `fmt`

func Arrangements(input []int, used []bool, res []int) {
	if len(input) == len(res) {
		fmt.Println(res)
		return
	}
	for idx, val := range input {
		if used[idx] == false {
			res = append(res, val)
			used[idx] = true
			Arrangements(input, used, res)
			used[idx] = false
			res = res[:len(res)-1]
		}
	}
}

func main() {
	input := []int{1, 2, 3}
	used := []bool{false, false, false}
	var res []int
	Arrangements(input, used, res)
}
