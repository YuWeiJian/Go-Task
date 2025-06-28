package main

import "fmt"

func singleNumber(nums []int) int {
	result := 0
	for _, num := range nums {
		result ^= num
	}
	return result
}
func singleNumberWithMap(nums []int) int {
	count := make(map[int]int)

	for _, num := range nums {
		count[num]++
	}

	for num, freq := range count {
		if freq == 1 {
			return num
		}
	}
	return -1
}

func main() {
	// 测试用例
	testCases := [][]int{
		{2, 2, 1},       // 期望输出: 1
		{4, 1, 2, 1, 2}, // 期望输出: 4
		{1},             // 期望输出: 1
		{1, 1, 2, 2, 3}, // 期望输出: 3
		{5, 3, 5, 3, 7}, // 期望输出: 7
	}

	for i, nums := range testCases {
		result := singleNumber(nums)
		fmt.Printf("测试用例 %d: %v -> %d\n", i+1, nums, result)
	}

	for i, nums := range testCases {
		result := singleNumberWithMap(nums)
		fmt.Printf("测试用例 %d: %v -> %d\n", i+1, nums, result)
	}

}
