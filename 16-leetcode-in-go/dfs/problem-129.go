package main

import `fmt`

type Node struct {
	val   int
	left  *Node
	right *Node
}

func getSum(t *Node, curSum int, res *int) {
	if t.left == nil && t.right == nil {
		*res += curSum
		return
	}
	if t.left != nil {
		getSum(t.left, curSum*10+t.left.val, res)
	}
	if t.right != nil {
		getSum(t.right, curSum*10+t.right.val, res)
	}
}

func main() {
	t1 := &Node{val: 4}
	t21 := &Node{val: 9}
	t22 := &Node{val: 0}
	t1.left = t21
	t1.right = t22
	t31 := &Node{val: 5}
	t32 := &Node{val: 1}
	t21.left = t31
	t21.right = t32

	res := new(int)
	getSum(t1, t1.val, res)
	fmt.Println(*res)
}
