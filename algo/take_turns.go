package main

import (
	"fmt"
	"sync"
)

func printOddNumbers(oddCh chan int, evenCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 10; i += 2 {
		fmt.Println("Odd:", i)
		oddCh <- i // 将奇数发送到oddCh通道
		<-evenCh   // 等待从evenCh通道接收信号
	}
}

func printEvenNumbers(oddCh chan int, evenCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 2; i <= 10; i += 2 {
		<-oddCh // 等待从oddCh通道接收信号
		fmt.Println("Even:", i)
		evenCh <- i // 将偶数发送到evenCh通道
	}
}

func main() {
	oddCh := make(chan int)
	evenCh := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)
	go printOddNumbers(oddCh, evenCh, &wg)
	go printEvenNumbers(oddCh, evenCh, &wg)

	wg.Wait()
	close(oddCh)
	close(evenCh)
}
