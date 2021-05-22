package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	closeRead()
	closeWrite()
	readClose()
	writeClose()
	closeRange()
	rangeClose()
}

func closeRead() {
	ch := make(chan int)

	var s sync.WaitGroup

	s.Add(1)
	go func() {
		defer s.Done()

		x, ok := <-ch
		fmt.Println("closeRead", "recv", x, ok) // recv 0 false
	}()

	time.Sleep(2 * time.Second)
	close(ch)

	s.Wait()
}

func closeWrite() {
	ch := make(chan int)

	var s sync.WaitGroup

	s.Add(1)
	go func() {
		defer s.Done()
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("closeWrite", "panic", err)
			}
		}()

		ch <- 1 // panic
		fmt.Println("closeWrite", "exit")
	}()

	time.Sleep(2 * time.Second)
	close(ch)

	s.Wait()
}

func readClose() {
	ch := make(chan int)
	close(ch)

	var s sync.WaitGroup

	s.Add(1)
	go func() {
		defer s.Done()

		time.Sleep(2 * time.Second)
		x, ok := <-ch
		fmt.Println("readClose", "recv", x, ok) // recv 0 false
	}()

	s.Wait()
}

func writeClose() {
	ch := make(chan int)
	close(ch)

	var s sync.WaitGroup

	s.Add(1)
	go func() {
		defer s.Done()
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("writeClose", "panic", err)
			}
		}()

		time.Sleep(2 * time.Second)
		ch <- 1
		fmt.Println("writeClose", "exit") // panic
	}()

	s.Wait()
}

func closeRange() {
	ch := make(chan int)

	var s sync.WaitGroup

	s.Add(1)
	go func() {
		defer s.Done()

		for x := range ch {
			fmt.Println("closeRange", "recv", x) // recv 0 false
		}
		fmt.Println("closeRange", "exit") // exit
	}()

	time.Sleep(2 * time.Second)
	close(ch)
	s.Wait()
}

func rangeClose() {
	ch := make(chan int)
	close(ch)

	var s sync.WaitGroup

	s.Add(1)
	go func() {
		defer s.Done()

		time.Sleep(2 * time.Second)
		for x := range ch {
			fmt.Println("rangeClose", "recv", x) // recv 0 false
		}
		fmt.Println("rangeClose", "exit") // exit
	}()

	s.Wait()
}
