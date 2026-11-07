package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("m <= 0", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)
		for range tasksCount {
			tasks = append(tasks, func() error {
				time.Sleep(100 * time.Millisecond)
				return nil
			})
		}
		workersCount := 3
		err := Run(tasks, workersCount, 0)
		require.ErrorIs(t, err, ErrErrorsLimitExceeded)
		err = Run(tasks, workersCount, -10)
		require.ErrorIs(t, err, ErrErrorsLimitExceeded)
	})

	t.Run("tasks without errors and without time.Sleep", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for range tasksCount {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100)+1)
			sumTime += taskSleep
			tasks = append(tasks, func() error {
				require.Eventually(t, func() bool {
					return true
				}, taskSleep*2, taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("only one worker", func(t *testing.T) {
		tasksCount := 100
		tasks := make([]Task, 0, tasksCount)
		for range tasksCount {
			tasks = append(tasks, func() error {
				<-time.Tick(time.Millisecond * time.Duration(rand.Intn(5)+1))
				return nil
			})
		}
		workersCount := 1
		err := Run(tasks, workersCount, 1)
		require.Nil(t, err)
	})

	t.Run("one error", func(t *testing.T) {
		tasksCount := 100
		tasks := make([]Task, 0, tasksCount)
		for range tasksCount {
			tasks = append(tasks, func() error {
				<-time.After(time.Millisecond * time.Duration(rand.Intn(5)+1))
				return nil
			})
		}
		tasks[49] = func() error {
			<-time.After(time.Millisecond * time.Duration(rand.Intn(5)+1))
			return fmt.Errorf("task error")
		}
		workersCount := 30
		err := Run(tasks, workersCount, 1)
		require.ErrorIs(t, err, ErrErrorsLimitExceeded)
	})

	t.Run("amount of tasks less or equal than n", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)
		for range tasksCount {
			tasks = append(tasks, func() error {
				<-time.After(time.Millisecond * time.Duration(rand.Intn(50)+1))
				return nil
			})
		}
		workersCount := 15
		err := Run(tasks, workersCount, 1)
		require.Nil(t, err)
		workersCount = 10
		err = Run(tasks, workersCount, 1)
		require.Nil(t, err)
	})
}
