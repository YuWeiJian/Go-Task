package main

import "fmt"

func main() {

	num := 1
	fmt.Printf("初始: %d\n", num)

	add(&num)
	fmt.Printf("增加后的值: %d\n", num)

	int_slice := []int{1, 2, 4, 5}
	fmt.Printf("原始切片: %v\n", int_slice)
	mulSlice(&int_slice)
	fmt.Printf("乘2后切片值: %v\n", int_slice)
}

/*
*题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
考察点 ：指针的使用、值传递与引用传递的区别。
*/
func add(num *int) {
	*num = *num + 10
}

/*
*题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
考察点 ：指针运算、切片操作。
*/
func mulSlice(nums *[]int) {
	for i := range *nums {
		(*nums)[i] = (*nums)[i] * 2
	}
}
