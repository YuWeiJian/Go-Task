package main

import (
	"fmt"
	"sort"
)

func main() {
	testCases := [][][]int{
		{{1, 3}, {2, 6}, {8, 10}, {15, 18}},       // [[1,6], [8,10], [15,18]]
		{{1, 4}, {4, 5}},                          // [[1,5]]
		{{1, 4}, {2, 3}},                          // [[1,4]]
		{{1, 3}, {2, 6}, {8, 10}, {15, 18}},       // [[1,6], [8,10], [15,18]]
		{{1, 4}, {0, 4}},                          // [[0,4]]
		{{1, 4}, {0, 0}},                          // [[0,0], [1,4]]
		{{2, 3}, {4, 5}, {6, 7}, {8, 9}, {1, 10}}, // [[1,10]]
	}

	fmt.Println("=== 区间合并测试 ===")
	for i, test := range testCases {
		result := merge(test)
		fmt.Printf("测试 %d: %v -> %v\n", i+1, test, result)
	}
}

func merge(intervals [][]int) [][]int {
	if len(intervals) <= 1 {
		return intervals
	}

	// 排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	var result [][]int
	start := intervals[0][0]
	end := intervals[0][1]

	for i := 1; i < len(intervals); i++ {
		if intervals[i][0] <= end {
			// 重叠，更新结束位置
			end = max(end, intervals[i][1])
		} else {
			// 不重叠，添加当前区间，开始新区间
			result = append(result, []int{start, end})
			start = intervals[i][0]
			end = intervals[i][1]
		}
	}

	// 添加最后一个区间
	result = append(result, []int{start, end})

	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
