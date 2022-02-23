package main

import (
	`fmt`
	`strconv`
	`strings`
)

/*
「候选」：1、19、192、...
「筛选」：不能以0开头，但是0本身允许；IP地址：0-255
「截止」：已经取得了4个数字（进入到第5个坑位，如果仍有数字剩余，则非法；否则是一个可行解
*/

// GetIP 这是我的版本，和大神的还有很大差距！
func GetIP(inStr string, pos int, res []int) {
	// 截至条件
	if len(res) == 4 {
		if pos == len(inStr)-1 {
			// 所有数字均遍历过了，是一个合法的解
			// 从[]int再转换为[]string就很冗余！
			var resStr []string
			for _, num := range res {
				resStr = append(resStr, strconv.Itoa(num))
			}
			fmt.Println(strings.Join(resStr, "."))
		}
		return
	}

	// 遍历
	for idx := range inStr {
		if idx <= pos {
			continue
		}
		offset := 0
		for {
			if idx+offset+1 > len(inStr) {
				break
			}
			newNumStr := inStr[pos+1 : pos+offset+2]
			if newNumStr[0] == '0' && len(newNumStr) > 1 {
				return
			}
			newNum, _ := strconv.Atoi(newNumStr)
			if newNum <= 255 {
				// 合法
				res = append(res, newNum)
				pos = idx + offset

				GetIP(inStr, pos, res)

				res = res[:len(res)-1]
				pos = idx - 1
			} else {
				return
			}
			offset++
		}
	}
}

// GetIP2 这是大神的写法
func GetIP2(inStr string, pos, level int, res []string) {
	// 当已经进入第5层或者已经遍历完全部元素的时候，就会触发截止条件
	if level == 5 || pos == len(inStr)-1 {
		// 触发后，判断是否为一个合法解
		if level == 5 && pos == len(inStr)-1 {
			fmt.Println(strings.Join(res, "."))
		}
		return
	}
	// 此处直接遍历offset！
	for offset := 1; offset < 4; offset++ {
		if pos+offset > len(inStr)-1 {
			break
		}
		nextNumStr := inStr[pos+1 : pos+offset+1]
		nextNum, _ := strconv.Atoi(nextNumStr)
		if nextNum < 256 && (nextNumStr == "0" || nextNumStr[0] != '0') {
			res = append(res, nextNumStr)
			GetIP2(inStr, pos+offset, level+1, res)
			res = res[:len(res)-1]
		}
	}
}

func main() {
	inStr := "19216801"
	pos := -1

	res := make([]int, 0)
	GetIP(inStr, pos, res)

	resStr := make([]string, 0)
	GetIP2(inStr, pos, 1, resStr)
}
