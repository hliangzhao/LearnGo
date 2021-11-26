package main

import (
	`context`
	`fmt`
	`time`
)

/*
context主要是用来在协程之间传递上下文信息，包括取消信号、超时时间、截止时间以及一些键值对等。
典型使用场景：客户端发起取消等待的信号，服务端收到之后就会终止响应本次请求，从而节约服务端资源。

创建一个空的context（没有ddl、timeout等信息）：
（1）ctx := context.Background()
	获得一个根context，该context不能被取消，用在main函数或者顶级请求中来派生其他context
（2）ctx := context.TODO()
	获得一个空context，只能用于高等级或者当您不确定使用什么context类型。

利用空context派生出新的context：
（1）ctx, cancel := context.WithCancel(ctx)
（2）ctx, cancel := context.WithTimeout(ctx, time duration)
（3）ctx, cancel := context.WithDeadline(ctx, time)
*/


// doSth 测试单进程内协程之间的context信号传递
func doSth(ctx context.Context) {
	select {
	case <- time.After(time.Second * 5):
		// 用time.After()模拟本服务端执行任务的时间开销
		fmt.Println("task finished")
	case <- ctx.Done():
		// ctx.Done()返回一个只读的通道，当ctx被取消的时候，这个通道会被客户传入一个值表示他已经被关闭了
		// 因此，对于本服务端，我们需要使用select语句来监听context.Done()是否有值
		// 如果context.Done()有值，则执行本case
		err := ctx.Err()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func main() {
	// 测试1 测试单进程内协程之间的context信号传递
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// main函数（客户）最多只能等待两秒
	go func() {
		time.Sleep(2 * time.Second)
		cancel()
	}()
	// 客户希望doSth在两秒之内可以执行完成。如果做不到，则客户会发送cancel信号到ctx.Done()这个通道中
	// 因为doSth作为服务端一直在监听这个通道是否有值，因此，一旦发现有值，则不再继续执行。
	doSth(ctx)
}