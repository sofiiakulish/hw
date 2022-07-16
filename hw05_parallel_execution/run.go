package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var wg sync.WaitGroup

	goroutines := make(chan int, n)
	errors := make(chan error, len(tasks))

	errorCount := 0

	for k, t := range tasks {
		goroutines <- k
		wg.Add(1)
		go func(t Task, goroutines <-chan int) {
			defer wg.Done()
			result := t()
			if result != nil {
				errors <- result
			}
			<-goroutines
		}(t, goroutines)
		select {
		case <-errors:
			errorCount++
			if m > 0 && errorCount >= m {
				return ErrErrorsLimitExceeded
			}
			default:
		}
	}

	wg.Wait()

	return nil
}
