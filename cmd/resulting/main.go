package main

import (
	"fmt"

	"github.com/thomasem/gowork"
)

func main() {
	resultsCh, rp := gowork.NewResultingPool[string]()
	rp.Start()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for result := range resultsCh {
			fmt.Println(result)
		}
	}()

	for i := range 20 {
		rp.Submit(func() string { return fmt.Sprintf("Job %d ran", i) })
	}

	rp.Stop()
	<-done
}
