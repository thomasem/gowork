package gowork

import (
	"testing"
)

func TestSimplePool(t *testing.T) {
	expected := 100
	pool := NewPool(WithLogger(NewNullLogger()))
	pool.Start()

	counter := 0
	for range expected {
		pool.Submit(func() {
			counter++
		})
	}

	pool.Stop()

	if counter != expected {
		t.Errorf("Expected counter to be %d, got %d", expected, counter)
	}
}

func TestResultingPool(t *testing.T) {
	expected := 100
	resultsCh, pool := NewResultingPool[int](WithLogger(NewNullLogger()))
	pool.Start()

	counter := 0
	go func() {
		for range resultsCh {
			counter++
		}
	}()

	for range expected {
		pool.Submit(func() int {
			return 0
		})
	}

	pool.Stop()

	if counter != expected {
		t.Errorf("Expected counter to be %d, got %d", expected, counter)
	}
}
