package main

import (
	"errors"
	"fmt"
	"strconv"
)

/* 结合切片和interface实现可装载通用数据类型的stack */

type GeneralStack []interface{}

func (s GeneralStack) Len() int {
	return cap(s)
}

func (s GeneralStack) Cap() int {
	return cap(s)
}

func (s GeneralStack) IsEmpty() bool {
	return len(s) == 0
}

func (s *GeneralStack) Push(e interface{}) {
	*s = append(*s, e)
}

func (s *GeneralStack) Pop() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("empty stack")
	}

	s1 := *s
	*s = s1[:len(s1)-1]
	return s1[len(s1)-1], nil
}

func TestGeneralStack() {
	s := &GeneralStack{}
	for i := 0; i < 10; i++ {
		s.Push(i)
		s.Push(float64(i) + 0.01)
		s.Push("'" + strconv.Itoa(i) + "'")
	}
	fmt.Println(*s)
}
