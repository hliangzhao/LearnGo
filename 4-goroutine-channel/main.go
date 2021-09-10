package main

import (
	`fmt`
	`sync`
)

/* 协程是Go语言中实现并发的一种方式。协程之间传递信息和同步的渠道是channel。*/

// wg有三个方法：add、wait、done
var wg sync.WaitGroup

func main() {
	// ****************************************************************
	/*
		顺序执行的代码
	*/
	sequentialRun()

	/*
		下面将开启一个协程来执行worker #1，即worker #1和主线程是并行执行的，所以整体花费1秒
	*/
	// 发起一个goroutine，只需要在这个函数的调用前面加上关键字"go"即可。
	// 对匿名函数发起goroutine也是可行的
	go say("worker #1")
	// go func(id string) {
	// 	time.Sleep(time.Second)
	// 	fmt.Println("I am done! id:", id)
	// }("worker #1")
	say("worker #2")

	/*
		此时main函数将很快退出，两个协程不会执行完就直接退出。这是因为主线程退出时，所有其附属的协程会被中断。
		如果主线程希望等待全体协程执行完再退出，那么，需要引入"wait group"进行同步。
	*/
	go say("worker #1")
	go say("worker #2")

	/*
		调用wg实现上面说的"主线程等待"。此时两个协程并行执行，因此总共花费1秒。
	*/
	wg.Add(2) // 主线程要等待两个协程完成执行
	go say2("worker #1")
	go say2("worker #2")
	wg.Wait()



	// ****************************************************************
	/*
		使用channel在协程之间传递信息
	*/
	wg.Add(2)
	ch := make(chan int, 0) // string是channel内存储的数据的类型，0是指定的buffer大小
	go player("Julia", ch)
	go player("Mike", ch)

	ch <- 0 // 主线程扮演裁判的角色，创建球并传给某个协程
	// 尤其要注意，通道要在协程被创建之后再创建，否则会因为没有协程消费这个通道数据而导致dead lock

	wg.Wait()



	// ****************************************************************
	/*
		共享变量如果不加锁，则多个协程一起修改此数据时会出现问题
	*/
	wg.Add(2)
	go UnsafeCounter()
	go UnsafeCounter()
	wg.Wait()
	fmt.Println(counter)

	/*
		解决方法1：使用互斥锁
	*/
	wg.Add(2)
	go SafeCounter()
	go SafeCounter()
	wg.Wait()
	fmt.Println(counter)

	/*
		解决方法2：将其申明为原子操作
	*/
	wg.Add(2)
	go AtomicIncCounter()
	go AtomicIncCounter()
	wg.Wait()
	fmt.Println(counter)

	/*
		解决方法3：将公共数据以channel的形式共享
	*/
	wg.Add(2)
	go ChannelIncCounter()
	go ChannelIncCounter()
	intCh <- 0
	wg.Wait()
	fmt.Println(<-intCh)
}
