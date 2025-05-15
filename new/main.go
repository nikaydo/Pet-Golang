package main

import (
	"fmt"
	"time"
)

func main() {

	ch1, ch2 := make(chan int, 10), make(chan int, 10)

	for i := range 10 {
		ch1 <- i
		ch2 <- i + 10
	}

	close(ch1)
	close(ch2)

	go func() {
		for range ch1 {
			fmt.Println(<-ch1)
		}
	}()

	go func() {
		for range ch2 {
			fmt.Println(<-ch2)
		}
	}()

	time.Sleep(time.Second * 2)
}
