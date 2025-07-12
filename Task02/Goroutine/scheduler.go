package main

import (
	"fmt"
	"sync"
	"time"
)

type Task func()

type TaskResult struct {
	Index    int
	Duration time.Duration
}

type TaskDemo struct{}

// 奇偶打印
func (td *TaskDemo) PrintOddEven() {
	chOdd := make(chan struct{})
	chEven := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i += 2 {
			<-chOdd
			fmt.Println(i)
			chEven <- struct{}{}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			<-chEven
			fmt.Println(i)
			if i != 10 {
				chOdd <- struct{}{}
			}
		}
	}()
	chOdd <- struct{}{} // 先让奇数先打印
	wg.Wait()
}

// 任务调度器
func (td *TaskDemo) Scheduler(tasks []Task) []TaskResult {
	var wg sync.WaitGroup
	results := make([]TaskResult, len(tasks))
	for i, task := range tasks {
		wg.Add(1)
		go func(idx int, t Task) {
			defer wg.Done()
			start := time.Now()
			t()
			elapsed := time.Since(start)
			results[idx] = TaskResult{Index: idx, Duration: elapsed}
		}(i, task)
	}
	wg.Wait()
	return results
}

func main() {
	demo := &TaskDemo{}
	fmt.Println("奇偶打印：")
	demo.PrintOddEven()

	tasks := []Task{
		func() { time.Sleep(500 * time.Millisecond); fmt.Println("任务1完成") },
		func() { time.Sleep(300 * time.Millisecond); fmt.Println("任务2完成") },
		func() { time.Sleep(700 * time.Millisecond); fmt.Println("任务3完成") },
	}
	fmt.Println("\n任务调度：")
	results := demo.Scheduler(tasks)
	for _, r := range results {
		fmt.Printf("任务%d耗时: %v\n", r.Index+1, r.Duration)
	}
}
