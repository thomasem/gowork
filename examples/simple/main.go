package main

import (
	"fmt"

	"github.com/thomasem/gowork"
)

func main() {
	pool := gowork.NewPool()
	pool.Start()

	for i := range 20 {
		pool.Submit(func() {
			fmt.Printf("Job %d is running...\n", i)
		})
	}

	pool.Stop()
}
