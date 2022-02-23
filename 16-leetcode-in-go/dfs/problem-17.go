package main

import (
	`fmt`
	`strconv`
)

var words []string

func AlphabetArrange(numbers string, level int, res string) {
	// 截止
	if level == len(numbers) {
		words = append(words, res)
		return
	}

	// 候选
	num := numbers[level]
	for _, word := range NumToWords(string(num)) {
		res += word
		AlphabetArrange(numbers, level+1, res)
		res = res[:len(res)-1]
	}
}

func NumToWords(num string) []string {
	if numInt, err := strconv.ParseInt(num, 10, 0); err == nil {
		if numInt < 7 {
			var res []string
			startVal := 97 + 3*(numInt-2)
			res = append(res, string(rune(startVal)), string(rune(startVal+1)), string(rune(startVal+2)))
			return res
		}
		if numInt == 7 {
			return []string{"p", "q", "r", "s"}
		}
		if numInt == 8 {
			return []string{"t", "u", "v"}
		}
		if numInt == 9 {
			return []string{"w", "x", "y", "z"}
		}
	}
	return nil
}

func main() {
	numbers := "234"
	var res string
	AlphabetArrange(numbers, 0, res)
	fmt.Print("[")
	for idx, str := range words {
		if idx != len(words)-1 {
			fmt.Print("\"", str, "\", ")
		} else {
			fmt.Print("\"", str, "\"")
		}
	}
	fmt.Println("].")
}
