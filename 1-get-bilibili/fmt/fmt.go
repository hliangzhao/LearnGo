package fmt

import (
	`io`
	`log`
	`os`
)

var Logger *log.Logger

// init在main之前被调用
func init() {
	file, err := os.OpenFile("trace.txt", os.O_CREATE | os.O_RDWR | os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	// 创建一个Logger实例，同时写内容到标准输出和文件
	Logger = log.New(io.MultiWriter(os.Stdout, file), "Log: ", log.LstdFlags)
}
