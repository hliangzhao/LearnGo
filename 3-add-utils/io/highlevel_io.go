package main

import (
	`bufio`
	`compress/gzip`
	`fmt`
	`io`
	`io/ioutil`
	`os`
	`strconv`
	`strings`
)

/* 调用ioutil提供的高级API接口 */

func CopyFile(srcFilename, dstFilename string) {
	buf, readerErr := ioutil.ReadFile(srcFilename)
	if readerErr != nil {
		fmt.Printf("error %v happened when reading file: %v", readerErr, srcFilename)
		os.Exit(1)
	}
	writerErr := ioutil.WriteFile(dstFilename, buf, 0644)
	if writerErr != nil {
		panic(writerErr.Error())
	}
}

type Page struct {
	Title string
	Body []byte
}

func (p *Page) save() (err error) {
	return ioutil.WriteFile(p.Title, p.Body, 0644)
}

func (p *Page) load(title string) (err error) {
	p.Title = title
	p.Body, err = ioutil.ReadFile(p.Title)
	return
}

func CopyFile2(srcFilename, dstFilename string) (written int64, err error) {
	src, srcErr := os.Open(srcFilename)
	if srcErr != nil {
		fmt.Printf("error %v happened when opening src file: %v", srcErr, srcFilename)
		os.Exit(1)
	}
	defer func() {
		err = src.Close()
		if err != nil {
			fmt.Printf("error %v happened when closing file: %v", err, srcFilename)
		}
	}()

	dst, dstErr := os.Create(dstFilename)
	if dstErr != nil {
		fmt.Printf("error %v happened when opening dst file: %v", dstErr, dstFilename)
		os.Exit(1)
	}
	defer func() {
		err = dst.Close()
		if err != nil {
			fmt.Printf("error %v happened when closing file: %v", err, dstFilename)
		}
	}()

	return io.Copy(src, dst)
}

func ReadZipFile(filename string) {
	inputFileHandle, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error %v happened when opening file: %v", err, filename)
	}
	defer func() {
		err = inputFileHandle.Close()
		if err != nil {
			fmt.Printf("error %v happened when closing file: %v", err, filename)
		}
	}()

	var r *bufio.Reader
	inputReader, err := gzip.NewReader(inputFileHandle)
	if err != nil {
		r = bufio.NewReader(inputFileHandle)
	} else {
		r = bufio.NewReader(inputReader)
	}

	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			fmt.Println("Read file done")
		}
		fmt.Print(line)
	}
}


type Product struct {
	Title string
	Price float64
	Quantity int
}

func ParseProducts(filename string) []Product {
	data := make([]Product, 0)
	inputFileHandle, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error %v happened when opening file: %v", err, filename)
	}
	defer func() {
		err = inputFileHandle.Close()
		if err != nil {
			fmt.Printf("error %v happened when closing file: %v", err, filename)
		}
	}()

	inputReader := bufio.NewReader(inputFileHandle)
	for {
		inStr, readerErr := inputReader.ReadString('\n')
		if readerErr == io.EOF {
			return data
		}
		inStr = inStr[:len(inStr) - 1]
		inArr := strings.Split(inStr, ";")
		price, _ := strconv.ParseFloat(inArr[1], 64)
		quantity, _ := strconv.Atoi(inArr[2])
		data = append(data, Product{inArr[0], price, quantity})
	}
}