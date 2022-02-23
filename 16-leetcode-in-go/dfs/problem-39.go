package main

import (
	`fmt`
	`reflect`
)

var exists []map[int]struct{}

func GetArr(candidates []int, target int, res []int, sum int) {
	// 截止条件：res中元素和大于等于target
	if sum >= target {
		if sum == target {
			pattern := make(map[int]struct{})
			for _, val := range res {
				pattern[val] = struct{}{}
			}
			// 判断当前pattern是否出现在exists中
			found := false
			for _, exist := range exists {
				if reflect.DeepEqual(exist, pattern) {
					found = true
				}
			}
			if !found {
				fmt.Println(res)
				exists = append(exists, pattern)
			}
		}
		return
	}

	// 候选
	for _, candidate := range candidates {
		// 筛选条件：加上新数字之后，大于target的会被筛除
		sum += candidate
		// 下面的判断是不需要的
		// if sum <= target {
		// 	res = append(res, candidate)
		// 	GetArr(candidates, target, res, sum)
		// 	res = res[:len(res)-1]
		// 	sum -= candidate
		// } else {
		// 	return
		// }
		res = append(res, candidate)
		GetArr(candidates, target, res, sum)
		res = res[:len(res)-1]
		sum -= candidate
	}
}

func main() {
	candidates := []int{2, 3, 5}
	target := 8
	sum := 0
	var res []int
	GetArr(candidates, target, res, sum)
}
