package main

import (
	`fmt`
	`strconv`
	`strings`
)

/* 自定义Error类型，捕获到panic之后执行相应处理 */

type ParseErr struct {
	Idx  int
	Word string
	Err  error
}

func (e *ParseErr) String() string {
	return fmt.Sprintf("error parsing %q as int", e.Word)
}

func Parse(in string) (numbers []int, err error) {
	// 捕获error要以defer func(){}()的形式编写！
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
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
			fmt.Println(err)
			continue
		}
		fmt.Printf("Successfully parsing as : %v\n", nums)
	}
}

func main() {
	TestPanicRecover()
}