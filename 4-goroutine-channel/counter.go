package main

import (
	`sync`
	`sync/atomic`
)

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
