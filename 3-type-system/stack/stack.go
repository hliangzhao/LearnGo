package main

import (
	`fmt`
	`log`
	`reflect`
	`strconv`
)

/* 运用struct和数组实现一个栈 */

const StackSize = 10000

type Stack struct {
	idx int                     // 第一个空位置的地址
	data [StackSize]float64
}

func (s *Stack) Push(value float64) {
	if s.idx == StackSize {
		log.Fatalln("stack is full now")
	}
	s.data[s.idx] = value
	s.idx++
}


func (s *Stack) Pop() float64 {
	s.idx--
	return s.data[s.idx]
}

func (s *Stack) Top() float64 {
	return s.data[s.idx - 1]
}

func (s *Stack) Peek() float64 {
	return s.data[s.idx - 1]
}

func (s *Stack) String() string {
	str := ""
	for i := 0; i < s.idx; i++ {
		str += "(" + strconv.Itoa(i) + ":" + strconv.FormatFloat(s.data[i], 'e', 1, 64) + ")\n"
	}
	return str
}

func TestStack() {
	s := Stack{}                 // 结构体变量实例
	for i := 0; i < 10; i++ {
		s.Push(float64(i))
	}
	fmt.Println(reflect.TypeOf(s), &s)

	sAddr := new(Stack)          // 结构体变量指针
	sAddr2 := &Stack{}           // 底层仍然会调用new()
	for i := 0; i < 10; i++ {
		sAddr.Push(float64(i))
	}
	fmt.Println(reflect.TypeOf(sAddr), sAddr, reflect.TypeOf(sAddr2))
}

func main() {
	TestStack()
}