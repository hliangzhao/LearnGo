package main

import `fmt`

type TreeNode struct {
	val   int
	left  *TreeNode
	right *TreeNode
}

func (t *TreeNode) Value() int {
	return t.val
}

func New(val int) *TreeNode {
	return &TreeNode{val: val}
}

func IsSameTree(p, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	}
	if p == nil || q == nil {
		return false
	}

	// 候选
	return p.val == q.val && IsSameTree(p.left, q.left) && IsSameTree(p.right, q.right)
}

func main() {
	t11 := New(1)
	t12 := New(2)
	t13 := New(3)
	t11.left = t12
	t11.right = t13

	t21 := New(1)
	t22 := New(2)
	// t23 := New(3)
	t21.left = t22
	// t21.right = t23

	fmt.Println(IsSameTree(t11, t21))
}
