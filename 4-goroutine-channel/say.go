package main

import (
	`fmt`
	`time`
)

// 一个计算任务
func say(id string) {
	time.Sleep(time.Second)
	fmt.Println("I am done! id:", id)
}


// 第一个say的调用花费1秒，结束后再进行第二个say的调用。
// 这是顺序执行的模式。
func sequentialRun() {
	say("worker #1")
	say("worker #2")
}


func say2(id string) {
	time.Sleep(time.Second)
	fmt.Println("I am done! id:", id)
	wg.Done()           // 任务完成
}