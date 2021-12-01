package main

import (
	`fmt`
	`strconv`
)

/* 实现golang版本的工厂模式 */

type Car struct {
	Model string
	Manufacture string
	BuildYear int
}

// Cars Car实例指针组成的切片
type Cars []*Car

func (c *Car) String() string {
	return "(" + c.Model + " -- " + c.Manufacture + " -- " + strconv.Itoa(c.BuildYear) + ")"
}

func (cs Cars) processAll(f func(c *Car)) {
	for _, c := range cs {
		f(c)
	}
}

// FindAll 遍历全部车，将所有使得g(c)为true的车组成切片返回
func (cs Cars) FindAll(g func(c *Car) bool) Cars {
	ret := make(Cars, 0)
	// 这一段和下面的一样
	// for _, c := range cs {
	// 	if g(c) {
	// 		ret = append(ret, c)
	// 	}
	// }
	cs.processAll(func(c *Car) {
		if g(c) {
			ret = append(ret, c)
		}
	})
	return ret
}

var cars Cars = []*Car{
	{"Fiesta", "Ford", 2008},
	{"XL 450", "BMW", 2011},
	{"D600", "Mercedes", 2009},
	{"X 800", "BMW", 2008},
}

func TestFindAll() {
	fmt.Println("All new BMWs:", cars.FindAll(func(c *Car) bool {
		return c.Manufacture == "BMW" && c.BuildYear > 2000
	}))
}

func MakeSortedAppender(manufacturers []string) (func(c *Car), map[string]Cars) {
	sortedCars := make(map[string]Cars)
	for _, m := range manufacturers {
		sortedCars[m] = make([]*Car, 0)
		// sortedCars[m] = make(Cars, 0)
	}
	sortedCars["other"] = make([]*Car, 0)

	appender := func(c *Car) {
		if _, ok := sortedCars[c.Manufacture]; ok {
			sortedCars[c.Manufacture] = append(sortedCars[c.Manufacture], c)
		} else {
			sortedCars["other"] = append(sortedCars["other"], c)
		}
	}
	return appender, sortedCars
}

func TestMakeSortedAppender() {
	manufacturers := []string{"Ford", "BMW", "Mercedes", "Jaguar"}
	// TODO：调用appender的时候外部的sortedCars也会被修改
	appender, sortedCars := MakeSortedAppender(manufacturers)
	cars.processAll(appender)
	fmt.Println(sortedCars)
}

func main() {
	TestFindAll()
	TestMakeSortedAppender()
}