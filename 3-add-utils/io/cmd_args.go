package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

/* 读取命令行参数 */

// ReadFromArgs 读取命令行参数并打印出来
func ReadFromArgs() {
	who := "Alice "
	// TODO：os.Args是所有在shell中打入的内容，第一个是程序的名字，紧接着的是参数
	if len(os.Args) > 1 {
		who += strings.Join(os.Args[1:], " ")
	}
	fmt.Println(who)
}

/* 使用并识别flag */

var NewLine = flag.Bool("n", false, "print newline")

const (
	Space   = " "
	Newline = "\n"
)

// UseFlag 解析读取的换行flag
// 示例：
// 输入：./cmd_args -n abc asd qwe
// 输出：
// abc
// asd
// qwe%
func UseFlag() {
	flag.PrintDefaults()
	// 解析定义的flags
	flag.Parse()
	var s = ""
	// flag.NArg()返回去除了flag的参数的个数（注意，是包含程序名这个"第一参数"的！）
	// TODO：注意区分flag和arg！
	for i := 0; i < flag.NArg(); i++ {
		if i > 0 {
			s += Space
			// 如果NewLine这个flag出现在了我参数中，则执行相应的操作
			if *NewLine {
				s += Newline
			}
		}
		s += flag.Arg(i)
	}
	// 将字符串写到标准输出
	_, err := os.Stdout.WriteString(s)
	if err != nil {
		fmt.Printf("error %v happened when writing to stdout", err)
		os.Exit(1)
	}
}

// CatWithFlag 如果没有参数，则echo输入的内容；如果有参数，则解析为文件名，将文件中的内容打印到标准输出
func CatWithFlag() {
	flag.Parse()
	if flag.NArg() == 0 {
		Cat1(bufio.NewReader(os.Stdin))
		// Cat2(os.Stdin)
	}
	for i := 0; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s:error reading from %s: %s\n", os.Args[0], flag.Arg(i), err.Error())
			continue
		}
		Cat1(bufio.NewReader(f))
		// Cat2(f)
		func() {
			err = f.Close()
			if err != nil {
				fmt.Printf("error %v happened when closing the file: %v\n", err, flag.Arg(i))
			}
		}()
	}
}

var PrintLineNum = flag.Bool("l", false, "print line number")

// Cat1 使用buffer存放数据
func Cat1(r *bufio.Reader) {
	lineNum := 1
	for {
		buf, err := r.ReadBytes('\n')
		if *PrintLineNum {
			_, _ = fmt.Fprintf(os.Stdout, "%d: %s", lineNum, buf)
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "%s", buf)
		}
		if err == io.EOF {
			break
		}
		lineNum++
	}
}

const SliceSize = 512

// Cat2 使用字节切片存放数据
func Cat2(f *os.File) {
	var buf [SliceSize]byte
	for {
		switch readBytesNum, readErr := f.Read(buf[:]); true {
		case readBytesNum < 0:
			_, _ = fmt.Fprintf(os.Stderr, "cat: error reading: %s\n", readErr.Error())
			os.Exit(1)
		case readBytesNum == 0: // EOF
			return
		case readBytesNum > 0:
			// write to os.Stdout
			if writtenBytesNUm, writeErr := os.Stdout.Write(buf[0:readBytesNum]); writtenBytesNUm != readBytesNum {
				_, _ = fmt.Fprintf(os.Stderr, "cat: error writing: %s\n", writeErr.Error())
				os.Exit(1)
			}
		}
	}
}

func main() {
	// ReadFromArgs()
	UseFlag()
}
