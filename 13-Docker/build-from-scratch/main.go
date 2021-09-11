package main

import (
	`os`
	`os/exec`
)

func myPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// run 实现docker run <cmd> <args>
func run() {
	// fmt.Println("run")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	// 保证cmd的输入输出和错误流是可见的
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// darwin不支持Sethostname方法！
	// syscall.Sethostname([]byte("container-hash"))
	// TODO：为darwin下的"docker"容器修改主机名

	// TODO：建立新的pid namespace实现宿主机和容器进程之间的隔离
	// 具体地，要把/proc目录挂载到ubuntu-fs的目录中，从而使得ps aux只能访问到目录内部的进程，而无法访问到宿主机的进程

	// TODO：在cgroup中创建新的文件夹并明确相应的资源使用，从而实现资源隔离

	myPanic(cmd.Run())

}

func main() {
	if len(os.Args) <= 1 {
		panic("no enough args")
	}

	if len(os.Args) <= 2 {
		panic("please specific the command to run")
	}

	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("unrecognized command")
	}
}
