package main

import (
	`fmt`
	`math/rand`
	`time`
)

// 消费通道数据的协程案例
func player(name string, ch chan int) {
	defer wg.Done() // 函数退出前自动触发

	for {
		ball, ok := <-ch // 从通道里面拿值。当通道被关闭的时候，ok为false
		if !ok {
			fmt.Printf("channel is closed. %s wins\n", name)
			return
		}
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(100)
		if n%10 == 0 {
			// 把球打飞
			close(ch)
			fmt.Printf("%s missed the ball\n", name)
			return
		}
		ball++
		fmt.Printf("%s receives the ball in round %d\n", name, ball)
		ch <- ball // 将数据写入通道
	}
}
