package main

import `fmt`

func main() {
	/* 数组 */
	arr := [5]string{"ok", "asd", "123", "f", "12321"}          // 指定数组大小
	arr2 := [...]string{"12321", "yes", "dsf"}                  // 编译器自行计算长度

	for idx, value := range arr {
		fmt.Println(idx, value)
	}
	ModifyArray(arr) // "传值"，arr本身并未被改变
	fmt.Println(arr)
	// 数组所申请的空间cap等于数组的长度len
	fmt.Println(arr2, len(arr2), cap(arr2))

	// deep copy，开辟了新的内存空间，将数据复制进去
	arr3 := arr
	arr3[0] = "abc"
	fmt.Println("arr:", arr, ", arr3:", arr3)




	/* 切片 */
	// 切片由一个首地址指针和两个表示长度和容量的整数组成
	// 切片是不指定大小的，是一种"动态数组"，可按需自动增长和缩小
	s := []string{"abc", "xyz"}
	fmt.Println(len(s), cap(s))

	// 使用make来创建数据类型（slice, map, channel）
	s2 := make([]int, 2, 2)
	fmt.Println(len(s2), cap(s2))
	s2 = append(s2, 1)                  // slice可以append
	fmt.Println(s2, len(s2), cap(s2))           // 容量自动翻倍，无需用户操心

	ModifySlice(s2) // "传址"
	fmt.Println(s2)

	// 下面的操作输出表明通过截取操作拿到的新切片并未开辟新的存储空间，而只是生成了对局部区域的新引用。
	// 修改这个引用所指向的空间内所存储的数据是会改变原切片的
	s3 := s2[:3]
	fmt.Println(s3, len(s3), cap(s3))
	s3[0] = -100
	fmt.Println(s2)

	// 多维切片（每个维度的大小可以不一样）
	multiSlice := make([][]int, 0)
	multiSlice = append(multiSlice, []int{1, 2, 34})
	multiSlice = append(multiSlice, []int{12, 21, 3, 12})
	fmt.Println(multiSlice)
	multiSlice[0][1] = 20
	fmt.Println(multiSlice)




	/* 映射 */
	m := map[string]int {
		"Julia":90,
		"Mike": 100,
	}
	// 若键不存在，则值未value类型的默认值（此处是0）
	fmt.Println(m["abc"], m["dfd"], len(m))

	// value，exist
	score, exist := m["abc"]
	fmt.Println(score, exist)

	// map作为函数的参数，和slice类似，是传址
	ModifyMap(m)
	fmt.Println(m)

	// 遍历
	for k, v := range m {
		fmt.Println("Key:", k, ", Value:", v)
	}

	// 删除
	delete(m, "Julia")
	fmt.Println(m)
}

// ModifyArray 参数传递形式：传值，即arr被copy了一份。因此原参不会被改变
func ModifyArray(arr [5]string) {
	// O(n)
	arr[0] = "abc"
}

// ModifySlice 同样的写法，对于slice而言，则是传址
func ModifySlice(s []int) {
	// O(1)，传进来的只是slice首地址
	s[0] = 100
}

// ModifyMap 同样的写法，对于map而言，也是传址
func ModifyMap(m map[string]int) {
	m["Narcissus"] = 98
}
