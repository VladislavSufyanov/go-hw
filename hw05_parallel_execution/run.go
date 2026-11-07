package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return nil
	}
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}
	amountOfTasks := len(tasks)
	var wg sync.WaitGroup
	chTasks := make(chan Task, n)
	chResults := make(chan error, amountOfTasks)
	for range min(n, amountOfTasks) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range chTasks {
				err := t()
				chResults <- err
			}
		}()
	}
	var taskIndex int
	for _, t := range tasks[:min(n, amountOfTasks)] {
		chTasks <- t
		taskIndex++
	}
	var numError int
	for range amountOfTasks {
		err := <-chResults
		if err != nil {
			numError++
			if numError >= m {
				close(chTasks)
				wg.Wait()
				close(chResults)
				return ErrErrorsLimitExceeded
			}
		}
		if taskIndex < amountOfTasks {
			chTasks <- tasks[taskIndex]
			taskIndex++
		}
	}
	close(chTasks)
	wg.Wait()
	close(chResults)
	return nil
}
