package main

import (
    "fmt"
    "sync"
    "time"
)

func dosomething(millisecs time.Duration, wg *sync.WaitGroup) {
    duration := millisecs * time.Millisecond
    time.Sleep(duration)
    fmt.Println("Function in background, duration:", duration)
    wg.Done()
}

func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go dosomething(200, &wg)
    wg.Add(1)
    go dosomething(400, &wg)
    wg.Add(1)
    go dosomething(150, &wg)
    wg.Add(1)
    go dosomething(600, &wg)

    wg.Wait()
    fmt.Println("Done")
}