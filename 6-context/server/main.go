package main

import (
	`fmt`
	`log`
	`net/http`
	`time`
)

func handler(w http.ResponseWriter, r *http.Request) {
	// 从request中拿出ctx
	// 从抓包上来看，跨进程的ctx，本质上是客户端提前向server端（也就是server端还在执行过程中）发送了fin，
	// 然后server端接到这个fin的时候就会通知所有监听了r.Context的地方。
	// 直接在服务器上打印r.Context的内容也可以看出客户端的ctx的timeOut属性是没有传递给server端的
	ctx := r.Context()

	// 真的在处理请求...那么消息来自于服务端自己创建的通道
	complete := make(chan struct{})

	go func() {
		// really do handling...
		time.Sleep(time.Second * 4)
		complete <- struct{}{}              // 创建一个空的struct的写法
	}()

	select {
	case <- complete:
		// Fprintf是将字符串写入一个writer中
		_, err := fmt.Fprintf(w, "server handling finished.\n")
		if err != nil {
			log.Fatalln(err)
		}
	// 用time.After()模拟本服务端执行任务的时间开销
	case <- time.After(time.Second * 3):
		_, err := fmt.Fprintf(w, "server sleeping finished.\n")
		if err != nil {
			log.Fatalln(err)
		}
	case <- ctx.Done():
		// ctx.Done()返回一个只读的通道，当ctx被取消的时候，这个通道会被客户传入一个值表示他已经被关闭了
		// 因此，对于本服务端，我们需要使用select语句来坚挺context.Done()是否有值
		err := ctx.Err()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatalln(http.ListenAndServe(":9000", nil))
}
