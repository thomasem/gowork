# GoWork: Simple Worker Pool Implementation

Minimal implementation of a worker pool in Golang that can be used to execute work concurrently.

## Basic Usage (concurrent `func()`)

```go
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

```

## Usage With Results (`func() T`)

**Note**: The results channel will be closed when the pool is stopped, so you have to stop the pool *before* waiting for
the results to finish.

```go
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
```
