package main

import (
	"fmt"
	"sync"
)

// 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
// 考察点 ：通道的基本使用、协程间通信。

// 获取
func getCh(ch <-chan int) {
	for i := range ch {
		fmt.Println("get 10:", i)
	}
}

func getCh2(ch <-chan int) {
	for i := range ch {
		fmt.Println("get 100:", i)
	}
}

// 发送
func sendCh(ch chan<- int) {
	for i := 1; i <= 10; i++ {
		ch <- i
		fmt.Println("Sent 10:", i)
	}
	close(ch) // 关闭通道，通知接收方没有更多数据
}

// 发送
func sendCh2(ch chan<- int) {
	for i := 1; i <= 100; i++ {
		ch <- i
		fmt.Println("Sent 100:", i)
	}
	close(ch) // 关闭通道，通知接收方没有更多数据
}

// 题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
// 考察点 ：通道的缓冲机制。

func main() {

	ch := make(chan int)

	//带有缓冲的通道
	ch2 := make(chan int, 10)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		sendCh(ch)
		sendCh2(ch2)
	}()
	go func() {
		wg.Done()
		getCh(ch)
		getCh2(ch2)
	}()

	wg.Wait()

}
