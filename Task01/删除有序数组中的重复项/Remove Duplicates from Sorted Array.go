package main

import "fmt"

func main() {

	// 测试代码
	fmt.Println(removeDuplicates([]int{1, 1, 2}))
}

func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	i := 0
	for j := 1; j < len(nums); j++ {
		if nums[i] != nums[j] {
			i++
			nums[i] = nums[j]
		}
	}
	return i + 1
}
