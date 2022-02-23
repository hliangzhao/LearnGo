package main

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO：自定义Error类型，捕获到panic之后执行相应处理：首先通过recover()获得error，然后解析它

// ParseErr 自定义error应当包含一个error类型的字段
type ParseErr struct {
	Idx  int
	Word string
	Err  error
}

// String 是自定义error必须要实现的方法
func (e *ParseErr) String() string {
	return fmt.Sprintf("error parsing %q as int", e.Word)
}

func Parse(in string) (numbers []int, err error) {
	// 捕获error的匿名函数总是会在Parse结束的时候被执行。如果没有错误，则正常推出，否则解析error
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				// TODO：这里要format的是r！而非err。err是执行type assertion时发生的错误，我们并不关心！
				err = fmt.Errorf("package: %v", r)
			}
		}
	}()

	fields := strings.Fields(in)
	if len(fields) == 0 {
		panic("no words to parse")
	}
	for idx, value := range fields {
		num, err := strconv.Atoi(value)
		if err != nil {
			// TODO：对于自定义的error，应当传入的，是实例的地址！
			//  否则自定义String()方法无法被执行
			panic(&ParseErr{idx, value, err})
		}
		numbers = append(numbers, num)
	}

	return
}

func TestPanicRecover() {
	var examples = []string{
		"1 2 3 4 5",
		"100 50 25 12.5 6.25",
		"2 + 2 = 4",
		"1st class",
		"",
	}
	for _, ex := range examples {
		fmt.Printf("Parsing %v:\n", ex)
		nums, err := Parse(ex)
		if err != nil {
			// TODO；如果产生了error，因为panic已经负责处理，因此这里直接continue即可
			fmt.Println(err)
			continue
		}
		fmt.Printf("Successfully parsing as : %v\n", nums)
	}
}

func main() {
	TestPanicRecover()
}
