package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// RootHandler 作为URL的处理函数，其参数固定为ResponseWriter和Request
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// TODO：第一件事，是从request中拿出ctx
	// 从抓包上来看，跨进程的ctx，本质上是客户端提前向server端（也就是server端还在执行过程中）发送了fin，
	// 然后server端接到这个fin的时候就会通知所有监听了r.Context的地方。
	// 直接在服务器上打印r.Context的内容也可以看出客户端的ctx的timeOut属性是没有传递给server端的
	ctx := r.Context()

	complete := make(chan struct{})

	go func() {
		fmt.Printf("handling %v...\n", r)
		time.Sleep(time.Second * 3)
		// TODO：处理完毕，就往complete这个context里面扔数据
		complete <- struct{}{}
	}()

	select {
	// TODO：监听这两个context的内容
	case <-complete:
		// Fprintf是将字符串写入一个writer中
		if _, err := fmt.Fprintf(w, "server handling finished.\n"); err != nil {
			log.Fatalln(err)
		}
	// 用time.After()模拟本服务端执行任务的时间开销
	// case <-time.After(time.Second * 4):
	// 	_, err := fmt.Fprintf(w, "server sleeping finished.\n")
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	case <-ctx.Done():
		// ctx.Done()返回一个只读的通道，当ctx被取消的时候，这个通道会被客户传入一个值表示他已经被关闭了
		// 因此，对于本服务端，我们需要使用select语句来监听context.Done()是否有值
		if err := ctx.Err(); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func main() {
	http.HandleFunc("/", RootHandler)
	log.Fatalln(http.ListenAndServe(":9000", nil))
}
