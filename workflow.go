package cloudflow

import (
	"fmt"

	"strings"

	"log"
	"os"
)

// Workflow contains tasks list of workflow definition.
type Workflow struct {
	tasks  []*namedTask
	logger *log.Logger
}

type namedTask struct {
	name string
	task Task
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
	return fmt.Errorf("workflow: task %v not found in: %v", name, wf.buildTaskFlowSummary(wf.tasks))
}

// RunOnly runs workflow only task specified.
func (wf *Workflow) RunOnly(name string) error {
	for i, t := range wf.tasks {
		if name == t.name {
			return wf.run(wf.tasks[i : i+1])
		}
	}
	return fmt.Errorf("workflow: task %v not found in: %v", name, wf.buildTaskFlowSummary(wf.tasks))
}

func (wf *Workflow) run(tasks []*namedTask) error {
	wf.logger.Print(fmt.Sprintf("workflow: Invoke workflow %v tasks: %v", len(tasks), wf.buildTaskFlowSummary(tasks)))

	for i, t := range tasks {
		wf.logger.Print(fmt.Sprintf("workflow: Start task: %v", tasks[i].name))
		if err := t.task.Execute(); err != nil {
			return err
		}
		wf.logger.Print(fmt.Sprintf("workflow: Complete task: %v", tasks[i].name))
	}
	return nil
}

func (wf *Workflow) buildTaskFlowSummary(tasks []*namedTask) string {
	names := make([]string, len(tasks))
	for i, t := range tasks {
		names[i] = t.name
	}
	return strings.Join(names, "->")
}

// Execute defined workflow that interface of Task.
// Workflow can be Task on other workflow definition.
func (wf *Workflow) Execute() error {
	return wf.Run()
}
