package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ： sync.Mutex 的使用、并发数据安全。
type Conter struct {
	mu    sync.Mutex
	count int
}

func (c *Conter) Add() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}
func (c *Conter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// 题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// 考察点 ：原子操作、并发数据安全。

type AtomicCounter struct {
	count int32
}

func (c *AtomicCounter) Add() {
	atomic.AddInt32(&c.count, 1)
}
func (c *AtomicCounter) Get() int32 {
	return atomic.LoadInt32(&c.count)
}
func main() {
	var wg sync.WaitGroup
	atomic_counter := &AtomicCounter{}
	c := &Conter{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				//互斥锁自增
				c.Add()
				//原子自增
				atomic_counter.Add()
			}
		}()
	}
	wg.Wait()

	fmt.Printf("最终 MutexConter count: %d\n", c.Get())

	fmt.Printf("最终 AtomicCounter count: %d\n", atomic_counter.Get())

}
