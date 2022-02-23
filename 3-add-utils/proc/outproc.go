package main

import (
	"fmt"
	"os"
)

func main() {
	env := os.Environ()
	// 通过os.StartProcess启动外部程序
	pid, err := os.StartProcess("/bin/ls", []string{"-al"}, &os.ProcAttr{
		Env:   env,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	if err != nil {
		fmt.Printf("Error %v starting process", err)
		os.Exit(1)
	}
	fmt.Printf("The process id is %v", pid.Pid)
}
