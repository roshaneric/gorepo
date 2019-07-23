package main

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

var theLimit int64 = 50

func main() {
	fmt.Println("**** Control workers through channel, unknown end ****")
	channelGopher()
	fmt.Println("**** Control workers through loop ****")
	loopGopher()
	fmt.Println("**** Control workers through channel, known end ****")
	loopChannelGopher()
}

func channelGopher() {
	result, e := subChannelGopher()
	if e != nil {
		fmt.Printf("%v\n", e)
		printResults(result)
		return
	}

	printResults(result)
}

func subChannelGopher() (map[int64]int64, error) {
	done := make(chan bool, 1)
	errCh := make(chan error, 1)
	workers := make(chan bool, 10)
	results := make(map[int64]int64)
	page := int64(0)
	var mux sync.Mutex
	var wg sync.WaitGroup

	for {
		select {
		case e := <-errCh:
			return results, e
		case <-done:
			wg.Wait()
			return results, nil
		default:
			workers <- true
			wg.Add(1)
			go func() {
				defer func() {
					<-workers
					wg.Done()
				}()
				pg := atomic.AddInt64(&page, 1)
				res, err := process(pg, theLimit)
				if err != nil {
					errCh <- err
					return
				}
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

type theResult struct {
	Page   int64
	Result int64
	Err    error
}

func loopGopher() {
	done := make(chan struct{})
	defer close(done)

	requests := make(chan int64)
	results := make(chan *theResult)
	go func() {
		defer close(requests)
		for page := int64(1); page <= 100; page++ {
			select {
			case requests <- page:
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
				processWrapper(req, theLimit, results, done)
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
		if r.Err != nil {
			fmt.Println(r.Err)
			return
		}
		count++
	}
	fmt.Printf("Count %v\n", count)
}

func loopChannelGopher() {
	result, e := subLoopChannelGopher()
	if e != nil {
		fmt.Printf("%v\n", e)
		printResults(result)
		return
	}

	printResults(result)
}

func subLoopChannelGopher() (map[int64]int64, error) {
	maxWorkers := 10
	errCh := make(chan error)
	done := make(chan bool, 1)
	workers := make(chan bool, maxWorkers)
	wg := &sync.WaitGroup{}
	results := make(map[int64]int64)
	var mux sync.Mutex

	for pg := int64(1); pg <= 100; pg++ {
		workers <- true
		wg.Add(1)
		go func(pg int64, errCh chan error, wg *sync.WaitGroup, workers chan bool) {
			defer func() {
				<-workers
				wg.Done()
			}()
			res, err := process(pg, theLimit)
			if err != nil {
				select {
				case errCh <- err: //since we are returning on first err, we need to make this non-clocking for multiple errs
				default:
				}
				return
			}
			mux.Lock()
			results[pg] = res
			mux.Unlock()
		}(pg, errCh, wg, workers)
	}
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case err := <-errCh:
		return results, err
	}
	close(errCh)
	return results, nil
}

func process(pg int64, limit int64) (int64, error) {
	result := 2 * pg
	if result > limit {
		return result, fmt.Errorf("%v is above limit", result)
	}
	return result, nil
}

func processWrapper(pg int64, limit int64, resultsCh chan<- *theResult, done <-chan struct{}) {
	result, err := process(pg, limit)
	select {
	case resultsCh <- &theResult{Page: pg, Result: result, Err: err}:
	case <-done:
		return
	}
}

func printResults(results map[int64]int64) {
	if results == nil {
		fmt.Println("nil results")
		return
	}
	var keys []int64
	for k := range results {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, k := range keys {
		fmt.Println("Key:", k, "Value:", results[k])
	}
}
