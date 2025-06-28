package main

import "fmt"

func main() {

	//测试代码
	fmt.Println(plusOne([]int{1, 2, 3}))
}

func plusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}
	return append([]int{1}, digits...)
}
