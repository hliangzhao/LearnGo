package main

import `fmt`

func dfs(inStr string, used []bool, level int, res string) {
	if level == len(inStr) {
		fmt.Println(res)
		return
	}
	for idx, c := range inStr {
		if used[idx] == false {
			res += string(c)
			used[idx] = true
			dfs(inStr, used, level+1, res)
			res = res[:len(res)-1]
			used[idx] = false
		}
	}
}

func dfs2(inStr string, used []bool, res string) {
	if len(inStr) == len(res) {
		fmt.Println(res)
		return
	}
	for idx, c := range inStr {
		if used[idx] == false {
			res += string(c)
			used[idx] = true
			dfs2(inStr, used, res)
			res = res[:len(res)-1]
			used[idx] = false
		}
	}
}

func main() {
	inStr := "ABC"
	used := []bool{false, false, false}
	res := ""
	dfs(inStr, used, 0, res)
	fmt.Println()
	dfs2(inStr, used, res)
}
