package cloudflow

import (
	"sync"

	multierror "github.com/hashicorp/go-multierror"
)

// Task represents task interface of workflow.
type Task interface {
	Execute() error
}

// ParallelTask represents parallel task on workflow.
type ParallelTask struct {
	tasks []Task
}

// NewParallelTask creates a parallel task by task list.
func NewParallelTask(tasks []Task) *ParallelTask {
	return &ParallelTask{tasks: tasks}
}

// Execute implement Task.Execute.
func (pt *ParallelTask) Execute() error {
	errChan := make(chan error)
	var wg sync.WaitGroup

	for _, t := range pt.tasks {
		wg.Add(1)
		go func(t Task) {
			if err := t.Execute(); err != nil {
				errChan <- err
			}
			wg.Done()
		}(t)
	}

	resultChan := make(chan error)
	go func() {
		var result *multierror.Error
		for err := range errChan {
			result = multierror.Append(result, err)
		}
		resultChan <- result.ErrorOrNil()
	}()

	wg.Wait()
	close(errChan)

	return <-resultChan
}
