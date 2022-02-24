package main

import (
	`fmt`
	`strconv`
)

// 19216801。问：有多少允许的IP地址（只需返回个数）
// 问题抽象：给定一个字符串，将其拆分成4个255以内的数的组合，有多少种拆法？
// TODO：dp[i][j]的含义为：str钱j个字符被拆分成i个255以内的数的组合的个数。dp的下标将用于递推。

func DPForIP(str string) int {
	// 首先分配数组大小
	length := len(str)
	dp := make([][]int, 5)
	for idx := range dp {
		dp[idx] = make([]int, length+1)
	}

	// 赋初值：0个数字拆分成8个数的组合，结果为1
	dp[0][0] = 1

	// 递推赋值
	for i := 0; i <= 4; i++ {
		// TODO: 当行号大于列号的时候（例如，将3个数字拆分成4个数的组合），显然个数为0，因为j只需要从i开始遍历即可
		for j := i; j <= length; j++ {
			// 赋初值
			if i == 0 || j == 0 {
				if !(i == 0 && j == 0) {
					dp[i][j] = 0
				}
			}
			if i >= 1 {
				// TODO：递推从第1行开始，最多在上一行往前推三列
				for x := 1; x <= 3; x++ {
					if j-x >= 0 {
						subStr := str[j-x : j]
						if validate(subStr) {
							dp[i][j] += dp[i-1][j-x]
						}
					}
				}
			}
		}
	}
	return dp[4][length]
}

func validate(s string) bool {
	if s == "0" {
		return true
	}
	if s[0] == '0' {
		return false
	}
	if val, _ := strconv.Atoi(s); val < 255 {
		return true
	}
	return false
}

func main() {
	fmt.Println(DPForIP("19216801"))
}
