package main

import (
	"fmt"
	"strconv"
)

func main() {
	// 测试几个不同的数字
	testCases := []int{3001003, 12321, 12345, -121, 10}
	for _, num := range testCases {
		fmt.Printf("%d is palindrome: %v\n", num, isPalindrome_num(num))
		fmt.Printf("%d is palindrome: %v\n", num, isPalindrome_str(num))
	}
}

func isPalindrome_str(x int) bool {
	if x < 0 {
		return false
	}
	str := strconv.Itoa(x)
	left_index := 0
	right_index := len(str) - 1
	for left_index < right_index {
		if str[left_index] != str[right_index] {
			return false
		}
		left_index++
		right_index--
	}
	return true
}

func isPalindrome_num(x int) bool {
	if x < 0 {
		return false
	}
	if x != 0 && x%10 == 0 {
		return false
	}
	reversed := 0
	for x > reversed {
		reversed = reversed*10 + x%10
		x /= 10
	}
	return x == reversed || x == reversed/10
}
