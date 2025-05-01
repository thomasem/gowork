package main

import (
	"fmt"

	"github.com/thomasem/gowork"
)

func main() {
	resultsCh, pool := gowork.NewResultingPool[string]()
	pool.Start()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for result := range resultsCh {
			fmt.Println(result)
		}
	}()

	for i := range 20 {
		pool.Submit(func() string { return fmt.Sprintf("Job %d ran", i) })
	}

	pool.Stop()
	<-done
}
