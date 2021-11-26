package main

import `fmt`

// 内置类型（basic types）
// bool, string, int (int8, int16, int32, int64), uint, uintptr, uint8, ...
// byte (alias for uint8)
// rune (alias for int32)
// float32, float64, complex64, complex128
// string

// 引用类型（在标头header里存放了指向底层数据结构的指针，总是浅复制）
// TODO：为什么引用类型的行为逻辑是"浅复制"呢？这是因为这些变量是在动态增长的，
//  程序在运行的时候无法获知其占用的内存空间范围。只能通过头指针来获取其访问位置。
// slice, map, channel, interface, func


// 类型、结构体的"大小写"决定了是否包外可见



// 使用接口

type notifier interface {
	// 所有实现了notify函数的结构体，都是一个notifier接口
	notify()
}

// 因为所有实现了notify函数的结构体，都是一个notifier接口
// 所以不同的、满足上述条件的结构体调用本方法都是合法的，因此实现了"多态机制"
func sendNotification(n notifier) {
	n.notify()
}



// 使用结构体

type User struct {
	// fields
	name string
	email string
}

// 使用"值接受者"，调用时使用这个值的"副本"来执行。显然，这不如指针接收者高效！
func (u User) notify()  {
	fmt.Printf("user %s notified\n", u.name)
}

func (u User) printName() {
	fmt.Println(u.name)
}

// 使用"指针接收者"，共享调用方法时接收者所指向的值
// 因此，如果想要一个方法修改调用者本身的字段，需要使用指针接收者
func (u *User) changeName(name string) {
	u.name = name
}

// 下面这个方法不会改变调用者的name字段
func (u User) changeName2(name string) {
	u.name = name
}


type Admin struct {
	name string
	email string
	admin bool
}

// 注意这里实现了notifier接口的形式是"指针接收者"，因此调用的时候要传递指针
func (a *Admin) notify() {
	fmt.Printf("admin %s notified\n", a.name)
}



// 结构体的嵌套

type Employee struct {
	level int
	a     Admin
}

func main() {
	u := User{name: "Narcissus", email: "narcissus@gmail.com"}
	u.changeName("Julia")
	fmt.Println(u)

	a := Admin{name: "Mike", email: "mike@google.com", admin: true}

	sendNotification(u)
	sendNotification(&a)            // 注意这里实现了notifier接口的形式是"指针接收者"，因此调用的时候要传递指针

	e := Employee{level: 0, a: Admin{
		name:  "Peter",
		email: "peter@abc.com",
		admin: false,
	}}
	fmt.Println(e)
	e.a.notify()
}
