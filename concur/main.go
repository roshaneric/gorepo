package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	gopher1()
}

func gopher1() {
	var mux sync.Mutex
	done := make(chan bool, 1)
	workers := make(chan bool, 10)
	results := make(map[int64]int64)
	page := int64(1)
	var wg sync.WaitGroup

	for {
		select {
		case <-done:
			count := 0
			wg.Wait()
			for k, r := range results {
				if k <= 100 {
					fmt.Println(r)
					count++
				}
			}
			fmt.Printf("Count %v\n", count)
			return
		default:
			workers <- true
			wg.Add(1)
			go func() {
				defer func() {
					<-workers
					wg.Done()
				}()
				pg := atomic.AddInt64(&page, 1)
				res := process(pg)
				mux.Lock()
				results[pg] = res
				mux.Unlock()
				if pg == 100 {
					done <- true
				}
			}()
		}
	}
}

func gopher2() {
	done := make(chan struct{})
	defer close(done)

	requests := make(chan int64)
	results := make(chan int64)
	go func() {
		defer close(requests)
		for page:= int64(1); page <= 1000; page++ {
			select {
			case requests <-page:
			case <-done:
				return
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(10) // 10 workers
	for i := 0; i < 10; i++ {
		go func() {
			for req := range requests {
				res := process(req)
				select {
				case results <- res:
				case <-done:
					return
				}
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	count := 0
	for r := range results {
		fmt.Println(r)
		count++
	}
	fmt.Printf("Count %v\n", count)

}

func process(pg int64) int64{
	return 2 * pg
}
