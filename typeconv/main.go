package main

import (
	"fmt"
)

func main() {
	var x interface{} = "zoom"
	i, ok := x.(string)

	if ok {
		fmt.Println(i)
	}

	fmt.Println("EOF")
}
