package main

import `fmt`

func Parenthesis(count, left, right int, res string) {
	if len(res) == 2*count {
		fmt.Println(res)
		return
	}

	// TODO：如下的判断过程很冗余，可以简化
	if left != right {
		// 可以选择左，也可以选择右
		if left <= 0 {
			if right <= 0 {
				// 结束了
				return
			} else {
				// 只能选择右
				right--
				res += ")"
				Parenthesis(count, left, right, res)
				res = res[:len(res)-1]
				right++
			}
		} else {
			if right <= 0 {
				// 只能选择左
				left--
				res += "("
				Parenthesis(count, left, right, res)
				res = res[:len(res)-1]
				left++
			} else {
				// 两边都能选
				for _, c := range "()" {
					res += string(c)
					if c == '(' {
						left--
					} else {
						right--
					}
					Parenthesis(count, left, right, res)
					if c == '(' {
						left++
					} else {
						right++
					}
					res = res[:len(res)-1]
				}
			}
		}
	} else {
		// 只能选择左
		if left <= 0 {
			return
		} else {
			left--
			res += "("
			Parenthesis(count, left, right, res)
			res = res[:len(res)-1]
			left++
		}
	}
}

func Parenthesis2(count, left, right int, res string) {
	if len(res) == 2*count {
		fmt.Println(res)
		return
	}

	// TODO: 只需要分别判断一次即可！
	if left > 0 {
		left--
		res += "("
		Parenthesis2(count, left, right, res)
		res = res[:len(res)-1]
		left++
	}
	if right > 0 && right != left {
		right--
		res += ")"
		Parenthesis2(count, left, right, res)
		res = res[:len(res)-1]
		right++
	}
}

func main() {
	n := 3
	var res string
	Parenthesis(n, 3, 3, res)
	Parenthesis2(n, 3, 3, res)
}
