package main

import (
	"fmt"
	"time"
)

// https://blog.golang.org/concurrency-timeouts
func main() {
	runTimeout()
}

func runTimeout() {
	ch := make(chan int)
	timeout := timeout(1 * time.Second)
	go timeoutSend(ch)
	select {
	case x := <-ch:
		// a read from ch has occurred
		fmt.Println("data", x)
	case <-timeout:
		// the read from ch has timed out
		fmt.Println("timeout")
	}
}

func timeoutSend(c chan<- int) {
	time.Sleep(3 * time.Second)
	c <- 1
}

func timeout(t time.Duration) <-chan bool {
	// return chan go return model
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(t)
		timeout <- true
	}()
	return timeout
}
