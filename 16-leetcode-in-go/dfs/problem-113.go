package main

import `fmt`

type treeNode struct {
	val   int
	left  *treeNode
	right *treeNode
}

func GetSumPath(t *treeNode, sum, val int, res []int) {
	if t != nil {
		val += t.val
		res = append(res, t.val)
		if t.left != nil {
			GetSumPath(t.left, sum, val, res)
		}
		if t.right != nil {
			GetSumPath(t.right, sum, val, res)
		}
		// 截止条件
		if val == sum {
			fmt.Println(res)
		} else {
			val -= t.val
			res = res[:len(res)-1]
			return
		}
	}
}

// GetSumPath2 大神的写法：严格按照三个步骤走
func GetSumPath2(t *treeNode, sum int, res []int) {
	// 是叶子结点，应该停止递归了
	if t.left == nil && t.right == nil {
		if sum == 0 {
			fmt.Println(res)
		}
		return
	}
	// 候选
	if t.left != nil {
		res = append(res, t.left.val)
		GetSumPath2(t.left, sum-t.left.val, res)
		res = res[:len(res)-1]
	}
	if t.right != nil {
		res = append(res, t.right.val)
		GetSumPath2(t.right, sum-t.right.val, res)
		res = res[:len(res)-1]
	}
}

func main() {
	t1 := &treeNode{val: 5}
	t21 := &treeNode{val: 4}
	t22 := &treeNode{val: 8}
	t1.left = t21
	t1.right = t22
	t31 := &treeNode{val: 11}
	t32 := &treeNode{val: 13}
	t33 := &treeNode{val: 4}
	t21.left = t31
	t22.left = t32
	t22.right = t33
	t41 := &treeNode{val: 7}
	t42 := &treeNode{val: 2}
	t43 := &treeNode{val: 5}
	t44 := &treeNode{val: 1}
	t31.left = t41
	t31.right = t42
	t33.left = t43
	t33.right = t44

	val := 0
	var res []int
	GetSumPath(t1, 22, val, res)

	var res2 []int
	res2 = append(res2, t1.val)
	GetSumPath2(t1, 22-t1.val, res2)
}
