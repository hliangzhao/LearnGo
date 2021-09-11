package main

import (
	`fmt`
	`os`
)

func CallOuterProc() {
	env := os.Environ()
	pid, err := os.StartProcess("/bin/ls", []string{"ls -al"}, &os.ProcAttr{
		Env:   env,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	if err != nil {
		fmt.Printf("Error %v starting process", err)
		os.Exit(1)
	}
	fmt.Printf("The process id is %v", pid)
}

func main() {
	CallOuterProc()
}
