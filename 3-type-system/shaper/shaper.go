package main

import (
	"fmt"
	"math"
)

type Square struct {
	side float64
}

type Circle struct {
	radius float64
}

// Shaper 接口命令的规范：以"er"结尾
type Shaper interface {
	Area() float64
}

func (s *Square) Area() float64 {
	return s.side * s.side
}

func (c *Circle) Area() float64 {
	return math.Pi * math.Pow(c.radius, 2)
}

func main() {
	var areaIntf Shaper
	// 分配/创建一个结构体指针
	s := new(Square)
	s.side = 5

	// use type assertion
	// 使用 var.(Intf_name) 来测试某个变量 var 是否实现了 Intf_name 接口
	areaIntf = s
	if t, ok := areaIntf.(*Square); ok {
		fmt.Printf("Square: %f\n", t.Area())
	}
	if t, ok := areaIntf.(*Circle); ok {
		fmt.Printf("Circle: %f\n", t.Area())
	}

	c := Circle{radius: 5}
	areaIntf = &c
	if t, ok := areaIntf.(*Square); ok {
		fmt.Printf("Square: %f\n", t.Area())
	}
	if t, ok := areaIntf.(*Circle); ok {
		fmt.Printf("Circle: %f\n", t.Area())
	}
}
