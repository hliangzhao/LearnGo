package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

/* 调用ioutil提供的高级API接口 */

// TODO：ioutil提供的读写方法操作的都是字节数组

func CopyFile(srcFilename, dstFilename string) {
	// TODO: 读出来的直接就是字节数组
	bytes, readerErr := ioutil.ReadFile(srcFilename)
	if readerErr != nil {
		fmt.Printf("error %v happened when reading file: %v", readerErr, srcFilename)
		os.Exit(1)
	}
	writerErr := ioutil.WriteFile(dstFilename, bytes, 0644)
	if writerErr != nil {
		panic(writerErr.Error())
	}
}

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() (err error) {
	return ioutil.WriteFile(p.Title, p.Body, 0644)
}

func (p *Page) load(title string) (err error) {
	p.Title = title
	p.Body, err = ioutil.ReadFile(p.Title)
	return
}

func CopyFile2(file1, file2 string) (written int64, err error) {
	inStream, _ := os.Open(file1)
	outStream, _ := os.Create(file2)
	defer func() {
		if err := inStream.Close(); err != nil {
			log.Fatalln("failed to close file")
		}
		if err = outStream.Close(); err != nil {
			log.Fatalln("failed to close file")
		}
	}()

	return io.Copy(outStream, inStream)
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
	// TODO：bufio和gzip提供的NewReader方法，可以理解为提供不同功能的奶嘴
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
			return
		}
		fmt.Print(line)
	}
}

type Product struct {
	Title    string
	Price    float64
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
		inStr = inStr[:len(inStr)-1]
		inArr := strings.Split(inStr, ";")
		price, _ := strconv.ParseFloat(inArr[1], 64)
		quantity, _ := strconv.Atoi(inArr[2])
		data = append(data, Product{inArr[0], price, quantity})

		// TODO：判断EOF要放在最后，因为如果源文件只有一行，那么readErr就会被置为EOF了，如果直接返回则data为空
		if readerErr == io.EOF {
			return data
		}
	}
}
