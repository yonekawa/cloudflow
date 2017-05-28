package cloudflow

import (
	"fmt"
	"reflect"
	"strings"

	"log"
	"os"
)

// Workflow contains tasks list of workflow definition.
type Workflow struct {
	tasks  []*namedTask
	logger *log.Logger
}

// NewWorkflow creates a new workflow definition.
func NewWorkflow() *Workflow {
	return &Workflow{
		tasks:  make([]*namedTask, 0),
		logger: log.New(os.Stdout, "[cloudflow] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// SetLogger sets log writer.
func (wf *Workflow) SetLogger(logger *log.Logger) {
	wf.logger = logger
}

// AddTask add task with name.
func (wf *Workflow) AddTask(name string, task Task) {
	wf.tasks = append(wf.tasks, &namedTask{name: name, task: task})
}

// Execute implement Task.Execute.
// Workflow can be Task on other workflow definition.
func (wf *Workflow) Execute() error {
	return wf.Run()
}

// Run defined workflow tasks.
func (wf *Workflow) Run() error {
	return wf.run(wf.tasks)
}

// RunFrom runs workflow from task specified.
func (wf *Workflow) RunFrom(name string) error {
	for i, t := range wf.tasks {
		if name == t.name {
			return wf.run(wf.tasks[i:])
		}
	}
	return fmt.Errorf("workflow: task %v not found in: %v", name, wf.Summary())
}

// RunOnly runs workflow only task specified.
func (wf *Workflow) RunOnly(name string) error {
	for i, t := range wf.tasks {
		if name == t.name {
			return wf.run(wf.tasks[i : i+1])
		}
	}
	return fmt.Errorf("workflow: task %v not found in: %v", name, wf.Summary())
}

func (wf *Workflow) run(tasks []*namedTask) error {
	for i, t := range tasks {
		wf.logger.Print(fmt.Sprintf("workflow: Start task: %v", tasks[i].name))
		if err := t.task.Execute(); err != nil {
			return err
		}
		wf.logger.Print(fmt.Sprintf("workflow: Complete task: %v", tasks[i].name))
	}
	return nil
}

// Summary returns task flow summary.
func (wf *Workflow) Summary() string {
	return buildTaskSummary(wf.tasks, " -> ", true)
}

func buildTaskSummary(tasks []*namedTask, delimiter string, showNumber bool) string {
	names := make([]string, len(tasks))
	for i, t := range tasks {
		var number string
		if showNumber {
			number = fmt.Sprintf("%d.", i+1)
		}
		if w, ok := t.task.(*Workflow); ok {
			names[i] = fmt.Sprintf("%s%s<Workflow>(%s)", number, t.name, w.Summary())
		} else if pt, ok := t.task.(*ParallelTask); ok {
			names[i] = fmt.Sprintf("%s%s<ParallelTask>(%s)", number, t.name, pt.Summary())
		} else {

			names[i] = fmt.Sprintf("%s%s<%s>", number, t.name, nameOfTask(t.task))
		}
	}
	return strings.Join(names, delimiter)
}

func nameOfTask(task Task) string {
	t := reflect.TypeOf(task)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}
