package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		fmt.Println("concurrency!")
		for t := range ticker.C {
			fmt.Printf("timer: %v\n", t)
		}
	}()
	fmt.Println("service waiting for signal")
	sysSig := make(chan os.Signal, 1)
	signal.Notify(sysSig, os.Interrupt, os.Kill)
	<-sysSig
	fmt.Println("service got quit signal")
	ticker.Stop()
	fmt.Println("timer stopped")
}
