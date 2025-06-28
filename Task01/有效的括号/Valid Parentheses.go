package main

import "fmt"

func isValid(s string) bool {
	//为空直接返回
	if s == "" {
		return false
	}
	if len(s)%2 != 0 {
		return false // 奇数长度直接返回false
	}

	stack := make([]rune, len(s)/2)
	top := 0

	for _, char := range s {
		switch char {
		case '(', '{', '[':
			if top >= len(stack) {
				//左括号超过一半的长度，直接返回false
				return false
			}
			// 将左括号加入入栈中
			stack[top] = char
			top++
		// 遇到右括号时，检查栈最上面的元素是否匹配，不匹配直接返回false
		case ')':
			if top == 0 || stack[top-1] != '(' {
				return false
			}
			top--
		case '}':
			if top == 0 || stack[top-1] != '{' {
				return false
			}
			top--
		case ']':
			if top == 0 || stack[top-1] != '[' {
				return false
			}
			top--
		}
	}

	return top == 0
}

func main() {
	testCases := []string{
		"()",     // true
		"()[]{}", // true
		"(]",     // false
		"([)]",   // false
		"{[]}",   // true
		"(((",    // false
		")))",    // false
		"",       // true
		"({[]})", // true
		"({[}])", // false
	}

	fmt.Println("=== 括号匹配测试 ===")
	for _, test := range testCases {
		result := isValid(test)
		fmt.Printf("输入: %-10s -> %t\n", test, result)
	}

}
