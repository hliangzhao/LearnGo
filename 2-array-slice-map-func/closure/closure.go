package main

import (
	`errors`
	`fmt`
	`math`
)


// Sqrt 在返回值处指定返回的变量名称
func Sqrt(f float64) (ret float64, err error) {
	if f < 0 {
		ret = math.NaN()
		err = errors.New("neg")
	} else {
		ret = math.Sqrt(f)
		err = nil
	}
	return
}


/* 返回一个函数的函数，被称为"闭包" */

func Compose(f, g func(x float64) float64) func (x float64) float64 {
	return func(x float64) float64 {
		return f(g(x))
	}
}

func increment(x float64) float64 {
	return x + 1
}

func scale(x float64) float64 {
	return 2 * x
}

func TestCompose() {
	compose := Compose(increment, scale)
	compose2 := func(x float64) float64 {
		return 2 * x + 1
	}
	fmt.Println(compose(2), compose2(2))
}

func Fib() func() int {
	// a和b的状态会被内部的函数记住
	a, b := 1, 1
	return func() int {
		a, b = b, a + b
		return b
	}
}

func TestFib() {
	fib := Fib()
	for i := 0; i < 5; i++ {
		fmt.Println(fib())
	}
}

func main() {
	TestCompose()
	TestFib()
}