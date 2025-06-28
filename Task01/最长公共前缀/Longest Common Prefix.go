package main

import "fmt"

func main() {
	//测试代码
	strs := []string{"flower", "flow", "flight"}
	fmt.Println(longestCommonPrefix(strs))
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	//找出最短字符串的长度
	minLength := len(strs[0])
	for i := 1; i < len(strs); i++ {
		if len(strs[i]) < minLength {
			minLength = len(strs[i])
		}
	}
	low, hight := 0, minLength
	for low < hight {

		//先最短字符长度的一半
		mid := (hight-low+1)/2 + low

		if !ispreifx(mid, strs) {
			hight = mid - 1
		} else {
			low = mid
		}
	}
	return strs[0][0:low]
}

func ispreifx(index int, strs []string) bool {
	str0, count := strs[0][:index], len(strs)
	for i := 0; i < count; i++ {
		if str0 != strs[i][:index] {
			return false
		}
	}
	return true
}
