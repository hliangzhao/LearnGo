package main

import (
	`fmt`
	`reflect`
	`time`
)

/*
反射：在程序动态运行时获得输入值的类型等信息。
最典型的案例是fmt.Println()，该函数动态获得传入的参数信息从而给出正确的格式输出。
*/

// TODO：1、获取对象的reflect元信息

// printMeta 输入参数是everything
func printMeta(obj interface{}) {
	// reflection obj主要是四种：type、kind、name和value
	t := reflect.TypeOf(obj)
	k := t.Kind()
	n := t.Name()
	v := reflect.ValueOf(obj)
	fmt.Printf("Type: %v,   Type.Kind: %v,   Type.Name: %v,   Value: %v\n", t, k, n, v)
}

func IntSub(a, b int) int {
	return a - b
}

type handler func(float64, float64) float64

func testReflection1() {
	var intVar int64 = 10
	stringVar := "hello"
	type book struct {
		name  string
		pages int
	}
	sum := func(a, b int) int {
		return a + b
	}
	var floatSub handler = func(a, b float64) float64 {
		return a - b
	}
	s := make([]int, 5, 10)

	printMeta(intVar)
	printMeta(stringVar)
	printMeta(book{name: "book", pages: 500})
	printMeta(sum)
	printMeta(IntSub)
	printMeta(floatSub)
	printMeta(s)
}

// TODO：2、通过Interface()和类型转换方法从reflect变量中拿到数据

func testReflection2() {
	floatVar := 3.14
	v := reflect.ValueOf(floatVar)
	// TODO：通过v.Interface()方法拿到这个数据（拿到的数据类型为interface{}），
	//  然后再通过cast方法转换为对应的数据类型
	newFloatVar := v.Interface().(float64)
	fmt.Println(newFloatVar)

	sliceVar := make([]int, 5)
	v = reflect.ValueOf(sliceVar)
	newSliceVar := v.Interface().([]int)
	for i := 0; i < 5; i++ {
		newSliceVar[i] = i
	}
	// 通过反射拿到的是引用
	fmt.Println(sliceVar, newSliceVar)

	// 当v还是reflect obj的时候，若v对应的数据类型为slice，则可以直接调用reflect.Append实现slice的Append
	v = reflect.Append(v, reflect.ValueOf(2))
	newSliceVar2 := v.Interface().([]int)
	fmt.Println(sliceVar, newSliceVar, newSliceVar2)
}

// TODO：3、通过reflect修改对象

func testReflection3() {
	floatVar := 3.14
	v := reflect.ValueOf(floatVar)
	fmt.Printf("is float canSet: %v, canAddr: %v\n", v.CanSet(), v.CanAddr())

	vp := reflect.ValueOf(&floatVar)
	// 用ptr.Elem()返回reflect ptr指向的元素
	fmt.Printf("is float canSet: %v, canAddr: %v\n", vp.Elem().CanSet(), vp.Elem().CanAddr())
	vp.Elem().SetFloat(2.3333)
	fmt.Println(floatVar, vp.Elem())
}

// TODO：4、通过reflect调用对象的方法

type Student struct {
	name string
}

func (s *Student) DoHomework() {
	fmt.Printf("%s is doing homework\n", s.name)
}

func (s *Student) DoHomework2(id int) {
	fmt.Printf("%s is doing homework #%d\n", s.name, id)
}

func testReflection4() {
	s := Student{name: "julia"}
	// 注意DoHomework是指针接收者，因此v的参数应该是&，否则后面的v.MethodByName("DoHomework")找不到
	v := reflect.ValueOf(&s)
	methodV := v.MethodByName("DoHomework")
	if methodV.IsValid() {
		methodV.Call([]reflect.Value{})
	}
	methodV = v.MethodByName("DoHomework2")
	in := []reflect.Value{reflect.ValueOf(2)}
	if methodV.IsValid() {
		methodV.Call(in)
	}
}

// TODO：5、通过reflect拿到一个函数并对其进行修改

// makeTimeFunc 将一个输入函数进行装饰。整个过程以reflect的方式执行
func makeTimeFunc(f interface{}) interface{} {
	t := reflect.TypeOf(f)
	v := reflect.ValueOf(f)
	if t.Kind() != reflect.Func {
		panic("expect a func")
	}

	wrapper := reflect.MakeFunc(t, func(args []reflect.Value) (result []reflect.Value) {
		start := time.Now()
		res := v.Call(args)
		end := time.Now()
		fmt.Printf("The func takes %v\n", end.Sub(start))
		return res
	})

	return wrapper.Interface()
}

func testFuncForMake() {
	time.Sleep(time.Second)
}

func testReflection5() {
	timedFunc := makeTimeFunc(testFuncForMake).(func())
	timedFunc()
}

// 6 some examples

type User struct {
	name string
	age  int
}

func (u User) PrintName() {
	fmt.Println(u.name)
}

func (u User) PrintAge() {
	fmt.Println(u.age)
}

type Aggregator func(int, int) int

var (
	user = User{
		name: "Julia",
		age:  24,
	}
	add Aggregator = func(a, b int) int {
		return a + b
	}
	sub Aggregator = func(a, b int) int {
		return a - b
	}
)

func inspect(variable interface{}) {
	t := reflect.TypeOf(variable)
	v := reflect.ValueOf(variable)
	if t.Kind() == reflect.Struct {
		// 获取结构体的全体fields
		fmt.Println("fields:")
		for idx := 0; idx < t.NumField(); idx++ {
			fieldType := t.Field(idx)
			fieldVal := v.Field(idx)
			fmt.Printf("\t %v = %v\n", fieldType.Name, fieldVal)
		}

		// 获取结构体的方法
		fmt.Println("methods:")
		for idx := 0; idx < t.NumMethod(); idx++ {
			methodType := t.Method(idx)
			fmt.Printf("\t input_num = %v, output_num = %v\n", methodType.Type.NumIn(), methodType.Type.NumOut())
		}
	} else if t.Kind() == reflect.Func {
		fmt.Printf("this func has %d inputs, %d outputs\n", t.NumIn(), t.NumOut())
	}
}

func testReflection6() {
	inspect(user)
	inspect(add)
	inspect(sub)
}

func main() {
	testReflection1()
	testReflection2()
	testReflection3()
	testReflection4()
	testReflection5()
	testReflection6()
}
