package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
	"sync"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var mu sync.Mutex

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				mu.Lock()
				atomic.AddInt32(&runTasksCount, 1)
				mu.Unlock()
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		mu.Lock()
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
		mu.Unlock()
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration
		var mu sync.Mutex

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				mu.Lock()
				atomic.AddInt32(&runTasksCount, 1)
				mu.Unlock()
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("negative error count", func(t *testing.T) {
			tasksCount := 50
    		tasks := make([]Task, 0, tasksCount)

    		var runTasksCount int32
    		var mu sync.Mutex

    		for i := 0; i < tasksCount; i++ {
    			err := fmt.Errorf("error from task %d", i)
    			tasks = append(tasks, func() error {
    				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
    				mu.Lock()
    				atomic.AddInt32(&runTasksCount, 1)
    				mu.Unlock()
    				return err
    			})
    		}

    		workersCount := 10
    		maxErrorsCount := -23
    		res := Run(tasks, workersCount, maxErrorsCount)

    		require.Equal(t, nil, res, "errors were thrown")
    		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})
}
