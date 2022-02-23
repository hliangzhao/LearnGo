package main

import `fmt`

var finalRes [][]int

func DFS(input []int, used []bool, res []int) {
	if len(input) == len(res) {
		// fmt.Println(finalRes)
		var equal bool
		for _, fr := range finalRes {
			// 比较fr和res是否相同
			equal = true
			for idx, frVal := range fr {
				if res[idx] != frVal {
					equal = false
					break
				}
			}
			if equal == true {
				break
			}
		}
		if equal != true {
			fmt.Println(res)
			finalRes = append(finalRes, res)
		}
		return
	}

	for idx, val := range input {
		if used[idx] == false {
			res = append(res, val)
			used[idx] = true
			DFS(input, used, res)
			used[idx] = false
			res = res[:len(res)-1]
		}
	}
}

func DFS2(input []int, count map[int]int, res []int) {
	if len(res) == len(input) {
		fmt.Println(res)
		return
	}
	for idx, val := range input {
		// TODO：将bool数组改成count，可以直接杜绝重复排列！
		if count[idx] > 0 {
			res = append(res, val)
			count[idx]--
			DFS2(input, count, res)
			count[idx]++
			res = res[:len(res)-1]
		}
	}
}

func getNum(input []int) map[int]int {
	res := make(map[int]int, 0)
	for _, val := range input {
		res[val]++
	}
	return res
}

func main() {
	input := []int{1, 1, 2}
	// used := []bool{false, false, false}
	count := getNum(input)
	var res []int
	// DFS(input, used, res)
	DFS2(input, count, res)
}
