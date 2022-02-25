package main

import `fmt`

// TODO：在这个例子中，dfs是作为一个子程序的，外部还有别的逻辑需要实现

func GetIslandNum(arr [][]int) int {
	if len(arr) == 0 {
		return 0
	}
	if len(arr[0]) == 0 {
		return 0
	}
	res := 0
	used := make([][]bool, len(arr))
	for i := 0; i < len(arr); i++ {
		used[i] = make([]bool, len(arr[i]))
	}
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			if used[i][j] == false && arr[i][j] == 1 {
				dfs200(arr, used, i, j)
				res++
			}
		}
	}
	return res
}

func dfs200(arr [][]int, used [][]bool, i, j int) {
	// 截止

	// 候选
	near := getNear(len(arr), len(arr[0]), i, j)
	for _, pos := range near {
		if arr[pos[0]][pos[1]] == 1 && used[pos[0]][pos[1]] == false {
			used[pos[0]][pos[1]] = true
			dfs200(arr, used, pos[0], pos[1])
		}
	}
}

func getNear(rowLen, colLen, i, j int) [][]int {
	var near [][]int
	if i-1 >= 0 {
		near = append(near, []int{i - 1, j})
	}
	if i+1 < rowLen {
		near = append(near, []int{i + 1, j})
	}
	if j-1 >= 0 {
		near = append(near, []int{i, j - 1})
	}
	if j+1 < colLen {
		near = append(near, []int{i, j + 1})
	}
	return near
}

func main() {
	arr := make([][]int, 4)
	arr[0] = []int{1, 1, 0, 0, 0}
	arr[1] = []int{1, 1, 0, 0, 0}
	arr[2] = []int{0, 0, 1, 0, 0}
	arr[3] = []int{0, 0, 0, 1, 1}
	fmt.Println(GetIslandNum(arr))
}
