package main

import (
	`fmt`
	`math/rand`
	`sync`
	`sync/atomic`
	`time`
)

/* 协程是Go语言中实现并发的一种方式。协程之间传递信息和同步的渠道是channel。*/


// 一个计算任务
func say(id string) {
	time.Sleep(time.Second)
	fmt.Println("I am done! id:", id)
}

// 默认的方式执行。第一个say的调用花费1秒，结束后再进行第二个say的调用。
// 这是顺序执行的模式。
func defaultRun() {
	say("worker #1")
	say("worker #2")
}


// wg有三个方法：add、wait、done
var wg sync.WaitGroup

func say2(id string) {
	time.Sleep(time.Second)
	fmt.Println("I am done! id:", id)
	wg.Done()           // 任务完成
}


// 消费通道数据的协程案例
func player(name string, ch chan int) {
	defer wg.Done()             // 函数退出前自动触发

	for {
		ball, ok := <- ch       // 从通道里面拿值。当通道被关闭的时候，ok为false
		if !ok {
			fmt.Printf("channel is closed. %s wins\n", name)
			return
		}
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(100)
		if n % 10 == 0 {
			// 把球打飞
			close(ch)
			fmt.Printf("%s missed the ball\n", name)
			return
		}
		ball++
		fmt.Printf("%s receives the ball in round %d\n", name, ball)
		ch <- ball              // 将数据写入通道
	}
}


// risk condition示例（多个线程对同一个共享资源的写入和读取产生了重叠）
// 这是因为++这个操作并非原子操作，即是编码为多行机器指令执行的
var counter int32

func UnsafeCounter() {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		counter++
	}
}

// 解决方法1：使用互斥锁
var mtx sync.Mutex

func SafeCounter() {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		mtx.Lock()
		counter++
		mtx.Unlock()
	}
}

// 解决方法2：将其申明为原子操作

func AtomicIncCounter() {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		// sync/atomic
		atomic.AddInt32(&counter, 1)
	}
}


// 解决方法3：将公共数据以channel的形式共享
// buffer设置为0会导致死锁。这是因为无缓冲的通道将像是一个没有厚度的门，这里有协程放，哪里就必须有协程取。
// 但是全部inc操作结束后没有协程来取。有缓冲的通道则不需要立即有协程来取，数据是最后由main来取的
var intCh = make(chan int, 1)

func ChannelIncCounter() {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		c := <- intCh
		c++
		intCh <- c
	}
}

func main() {
	defaultRun()

	/* 下面将开启一个协程来执行worker 1，即worker 1和主线程是并行执行的，所以整体花费1秒 */
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

	/* 调用wg实现上面说的"主线程等待"。此时两个协程并行执行，因此总共花费1秒。 */
	wg.Add(2)       // 主线程要等待两个协程完成执行
	go say2("worker #1")
	go say2("worker #2")
	wg.Wait()


	/* 使用channel在协程之间传递信息 */
	wg.Add(2)
	ch := make(chan int, 0)      // string是channel内存储的数据的类型，0是指定的buffer大小
	go player("Julia", ch)
	go player("Mike", ch)

	ch <- 0                      // 主线程扮演裁判的角色，创建球并传给某个协程
	// 尤其要注意，通道要在协程被创建之后再创建，否则会因为没有协程消费这个通道数据而导致dead lock

	wg.Wait()


	/* 共享变量如果不加锁，则多个协程一起修改此数据时会出现问题 */
	wg.Add(2)
	go UnsafeCounter()
	go UnsafeCounter()
	wg.Wait()
	fmt.Println(counter)


	/* 解决方法1：使用互斥锁 */
	wg.Add(2)
	go SafeCounter()
	go SafeCounter()
	wg.Wait()
	fmt.Println(counter)

	/* 解决方法2：将其申明为原子操作 */
	wg.Add(2)
	go AtomicIncCounter()
	go AtomicIncCounter()
	wg.Wait()
	fmt.Println(counter)

	/* 解决方法3：将公共数据以channel的形式共享 */
	wg.Add(2)
	go ChannelIncCounter()
	go ChannelIncCounter()
	intCh <- 0
	wg.Wait()
	fmt.Println(<- intCh)
}
