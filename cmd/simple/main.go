package main

import (
	"fmt"

	"github.com/thomasem/gowork"
)

func main() {
	pool := gowork.NewPool()
	pool.Start()
	defer pool.Stop()

	for i := range 100 {
		pool.Submit(func() {
			fmt.Printf("Job %d is running...\n", i)
		})
	}
}
