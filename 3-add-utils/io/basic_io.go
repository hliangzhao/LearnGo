package main

import (
	`bufio`
	`fmt`
	`io`
	`log`
	`os`
	`strconv`
	`strings`
)

/* 从标准输入中读取数据到缓冲区：示例1 */

// TODO：同样地，os.Stdin之类的属于IO流，buffer reader或buffer writer是套接在流之上方便存取数据的"奶嘴"，
//  ReadString或WriteString则可以理解为奶嘴提供的具体的存储数据方法，如"吸"、"咬"等。显然，读取或写入的直接就是字符串了
//  总结一下：
//  流来自package os（标准流os.Stdin, os.Stdout等；文本流通过os.OpenFile, os.Open, os.Create等方式创建）；
//  有了流之后，还需要套上一个奶嘴才可以进行读写的操作，而奶嘴通常来自bufio（NewReader或NewWriter）；
//  然后，就可以调用奶嘴的Read、ReadString、Write等方法正式开始读写
//  当然，还可以通过os.Read等方法直接读取（内部实现是先open后read）

// Calculator 从标准输入中读取逆波兰表达式并计算
func Calculator() {
	// TODO：bufio reader是一个套在标准输入流上的奶嘴
	reader := bufio.NewReader(os.Stdin)
	calc := new(Stack)
	fmt.Println("Input numbers (0 ~ 9999) and operators (+, -, *, /), type q at the end to quit:")
	// 调用ReadString读取一行字符串
	// TODO：注意，这一行字符串是带有换行符的！
	tokenStr, err := reader.ReadString('\n')
	if err != nil {
		os.Exit(1)
	}

	// TODO：手动去除换行符
	tokenStr = tokenStr[:len(tokenStr)-1]
	tokenArr := strings.Fields(tokenStr)
	for _, tk := range tokenArr {
		if tk == "q" {
			fmt.Println(calc.Pop())
			return
		} else if tk >= "0" && tk <= "9999" {
			i, _ := strconv.ParseFloat(tk, 64)
			calc.Push(i)
		} else if tk == "+" {
			// p is the former, q is the later
			q := calc.Pop()
			p := calc.Pop()
			calc.Push(p + q)
		} else if tk == "-" {
			q := calc.Pop()
			p := calc.Pop()
			calc.Push(p - q)
		} else if tk == "*" {
			q := calc.Pop()
			p := calc.Pop()
			calc.Push(p * q)
		} else if tk == "/" {
			q := calc.Pop()
			p := calc.Pop()
			calc.Push(p / q)
		} else {
			fmt.Println("Found not supported char!")
			os.Exit(1)
		}
	}
}

const StackSize = 10000

type Stack struct {
	idx  int // 第一个空位置的地址
	data [StackSize]float64
}

func (s *Stack) Push(value float64) {
	if s.idx == StackSize {
		// TODO：log.Fatalln很好用！是println跟上一个os.Exit
		log.Fatalln("stack is full now")
	}
	s.data[s.idx] = value
	s.idx++
}
func (s *Stack) Pop() float64 {
	s.idx--
	return s.data[s.idx]
}
func (s *Stack) Top() float64 {
	return s.data[s.idx-1]
}
func (s *Stack) Peek() float64 {
	return s.data[s.idx-1]
}
func (s *Stack) String() string {
	str := ""
	for i := 0; i < s.idx; i++ {
		str += "(" + strconv.Itoa(i) + ":" + strconv.FormatFloat(s.data[i], 'e', 1, 64) + ")\n"
	}
	return str
}

/* 从标准输入中读取数据到缓冲区：示例2 */

// ReadCount 从标准输入中读取多行文字并返回统计结果
func ReadCount() {
	nChars, nWords, nLines := 0, 0, 0
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter, type S to exit: ")
	for {
		// keep '\n'
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("error %v happened when read from stdin\n", err)
			os.Exit(1)
		}
		if input == "S\n" || input == "s\n" {
			break
		}
		// TODO：因为input是带有换行符的，因此这里统计的时候需要减1
		nChars += len(input) - 1
		nWords += len(strings.Fields(input))
		nLines++
	}
	fmt.Printf("Chars: %v\nWords: %v\nLines: %v\n", nChars, nWords, nLines)
}

/* 从文本中读取数据到缓冲区 */

// ReadFromFile 从filename标识的文件中读取内容
func ReadFromFile(filename string) {
	fileHandler, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error %v happened when opening the file: %v\n", err, filename)
		os.Exit(1)
	}
	defer func() {
		if err := fileHandler.Close(); err != nil {
			fmt.Printf("error %v happened when closing the file: %v\n", err, filename)
			os.Exit(1)
		}
	}()

	// TODO：bufio reader是一个套在文本输入流上的奶嘴
	inputReader := bufio.NewReader(fileHandler)
	for {
		inStr, err := inputReader.ReadString('\n')
		fmt.Println(inStr)
		if err == io.EOF {
			return
		}
	}
}

/* 将数据写入文本：创建文件句柄，将数据写入用该句柄创建的缓冲区，调用缓冲区的flush方法 */

// Write2File 将字符串写入文本
func Write2File(filename string) {
	// TODO：os.OpenFile是一个基本方法，提供的选项比较多
	//  基于OpenFile的方法有os.Open、os.Create等，前者是打开一个只读的文本流
	fileHandler, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("error %v happened when opening / creating the file: %v\n", err, filename)
		os.Exit(1)
	}
	defer func() {
		fileCloseErr := fileHandler.Close()
		if fileCloseErr != nil {
			fmt.Printf("error %v happened when closing the file: %v\n", fileCloseErr, filename)
			os.Exit(1)
		}
	}()

	// TODO：bufio writer是一个套在文本输入流上的奶嘴
	writer := bufio.NewWriter(fileHandler)
	_, err = writer.WriteString("Some sentences")
	if err != nil {
		fmt.Printf("error %v happened when write a string to the file handle for file: %v\n", err, filename)
		os.Exit(1)
	}
	err = writer.Flush()
	if err != nil {
		fmt.Printf("error %v happened when flushing\n", err)
		os.Exit(1)
	}
}

/* 从文本中读取数据，将读到的数据写入另一个文本中 */

func ReadThenWrite(srcFilename, dstFilename string) {
	inFileHandler, errIn := os.Open(srcFilename) // os.Open()以只读的形式打开文本
	outFileHandler, errOut := os.OpenFile(dstFilename, os.O_WRONLY|os.O_CREATE, 0644)
	if errIn != nil || errOut != nil {
		fmt.Printf("error %v or %v happened when open files\n", errIn, errOut)
		os.Exit(1)
	}
	defer func() {
		inFileCloseErr := inFileHandler.Close()
		if inFileCloseErr != nil {
			fmt.Printf("error %v happened when closing the file: %v\n", inFileCloseErr, srcFilename)
			os.Exit(1)
		}
		outFileCloseErr := outFileHandler.Close()
		if outFileCloseErr != nil {
			fmt.Printf("error %v happened when closing the file: %v\n", outFileCloseErr, dstFilename)
			os.Exit(1)
		}
	}()

	reader := bufio.NewReader(inFileHandler)
	writer := bufio.NewWriter(outFileHandler)
	for {
		inStr, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		}
		outStr := string(inStr) + "\n"
		_, err = writer.WriteString(outStr)
		if err != nil {
			fmt.Printf("error %v happened when writing\n", err)
			os.Exit(1)
		}
		err = writer.Flush()
		if err != nil {
			fmt.Printf("error %v happened when flushing\n", err)
			os.Exit(1)
		}
	}
}

func ReadThenWrite2(file1, file2 string) {
	inStream, err := os.Open(file1)
	if err != nil {
		log.Fatalln("failed to open file")
	}
	outStream, err := os.OpenFile(file2, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalln("failed to create file")
	}
	defer func() {
		if err = inStream.Close(); err != nil {
			log.Fatalln("failed to close file")
		}
		if err = outStream.Close(); err != nil {
			log.Fatalln("failed to close file")
		}
	}()

	reader := bufio.NewReader(inStream)
	writer := bufio.NewWriter(outStream)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalln("failed to read string")
		}
		if _, err = writer.WriteString(str); err != nil {
			log.Fatalln("failed to write string")
		}
		if err = writer.Flush(); err != nil {
			log.Fatalln("failed to flush")
		}
	}
}
