package cloudflow

import (
	"sync"

	multierror "github.com/hashicorp/go-multierror"
)

// Task represents task interface of workflow.
type Task interface {
	Execute() error
}

type namedTask struct {
	name string
	task Task
}

// ParallelTask represents parallel task on workflow.
type ParallelTask struct {
	tasks []*namedTask
}

// NewParallelTask creates a parallel task by task list.
func NewParallelTask() *ParallelTask {
	return &ParallelTask{tasks: make([]*namedTask, 0)}
}

func (pt *ParallelTask) AddTask(name string, task Task) {
	pt.tasks = append(pt.tasks, &namedTask{name: name, task: task})
}

// Execute implement Task.Execute.
func (pt *ParallelTask) Execute() error {
	errChan := make(chan error)
	var wg sync.WaitGroup

	for _, nt := range pt.tasks {
		wg.Add(1)
		go func(t Task) {
			if err := t.Execute(); err != nil {
				errChan <- err
			}
			wg.Done()
		}(nt.task)
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
