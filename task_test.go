package cloudflow

import "testing"

var results = make(map[string]bool, 0)
var resultChan = make(chan string)

type testParallelTask struct {
	name string
}

func (t *testParallelTask) Execute() error {
	resultChan <- t.name + " execute"
	return nil
}

func TestParallelTask(t *testing.T) {
	t.Parallel()

	completeChan := make(chan bool)
	go func() {
		for key := range resultChan {
			results[key] = true
		}
		completeChan <- true
	}()

	wf := NewWorkflow()
	pt := NewParallelTask()
	pt.AddTask("a", &testParallelTask{name: "a"})
	pt.AddTask("b", &testParallelTask{name: "b"})
	pt.AddTask("c", &testParallelTask{name: "c"})
	wf.AddTask("abc", pt)
	wf.AddTask("d", &testParallelTask{name: "d"})
	if err := wf.Run(); err != nil {
		t.Error(err)
	}

	close(resultChan)
	<-completeChan

	for name, b := range results {
		if !b {
			t.Errorf("workflow: task: %v is not running", name)
		}
	}
}
