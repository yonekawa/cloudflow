package cloudflow

import (
	"errors"
	"testing"
)

var runResult = make([]string, 0)

type testTask struct {
	name string
}

func (t *testTask) Execute() error {
	runResult = append(runResult, t.name+" execute")
	return nil
}

type errorTask struct {
	name string
}

func (t *errorTask) Execute() error {
	runResult = append(runResult, t.name+" execute")
	return errors.New("fail")
}

func TestWorkflow_Run(t *testing.T) {
	t.Parallel()

	runResult = make([]string, 0)

	wf := NewWorkflow()
	wf.RegisterTask("a", &testTask{name: "a"})
	wf.RegisterTask("b", &testTask{name: "b"})
	err := wf.Run()

	if err != nil {
		t.Error(err)
	}
	tests := []string{"a execute", "b execute"}
	if len(runResult) != len(tests) {
		t.Errorf("workflow: incorrect task result length expect:%v got:%v ", len(tests), len(runResult))
	}
	for i, test := range tests {
		if runResult[i] != test {
			t.Errorf("workflow: invalid task result expect:%v got: %v", test, runResult[i])
		}
	}

	runResult = make([]string, 0)

	wf = NewWorkflow()
	wf.RegisterTask("a", &testTask{name: "a"})
	wf.RegisterTask("b", &testTask{name: "b"})
	wf.RegisterTask("c", &errorTask{name: "c"})
	wf.RegisterTask("d", &testTask{name: "d"})
	err = wf.Run()

	if err == nil {
		t.Error("workflow: workflow not raises error when task failed")
	}

	tests = []string{"a execute", "b execute", "c execute"}
	if len(runResult) != len(tests) {
		t.Errorf("workflow: incorrect task result length expect:%v got:%v ", len(tests), len(runResult))
	}

	for i, test := range tests {
		if runResult[i] != test {
			t.Errorf("workflow: invalid task result expect:%v got: %v", test, runResult[i])
		}
	}

	runResult = make([]string, 0)

	wf = NewWorkflow()
	wf.RegisterTask("a", &testTask{name: "a"})
	wf2 := NewWorkflow()
	wf2.RegisterTask("b", &testTask{name: "b"})
	wf2.RegisterTask("c", &testTask{name: "c"})
	wf.RegisterTask("bc", wf2)
	wf.RegisterTask("d", &testTask{name: "d"})
	err = wf.Run()

	if err != nil {
		t.Error(err)
	}
	tests = []string{"a execute", "b execute", "c execute", "d execute"}
	if len(runResult) != len(tests) {
		t.Errorf("workflow: incorrect task result length expect:%v got:%v ", len(tests), len(runResult))
	}
	for i, test := range tests {
		if runResult[i] != test {
			t.Errorf("workflow: invalid task result expect:%v got: %v", test, runResult[i])
		}
	}

	runResult = make([]string, 0)

	subFlow := NewWorkflow()
	subFlow.RegisterTask("b", &testTask{name: "b"})
	subFlow.RegisterTask("c", &errorTask{name: "c"})

	wf = NewWorkflow()
	wf.RegisterTask("a", &testTask{name: "a"})
	wf.RegisterTask("bc", subFlow)
	wf.RegisterTask("d", &testTask{name: "d"})
	err = wf.Run()

	if err == nil {
		t.Error("workflow: workflow not raises error when task failed")
	}
	tests = []string{"a execute", "b execute", "c execute"}
	if len(runResult) != len(tests) {
		t.Errorf("workflow: incorrect task result length expect:%v got:%v ", len(tests), len(runResult))
	}
	for i, test := range tests {
		if runResult[i] != test {
			t.Errorf("workflow: invalid task result expect:%v got: %v", test, runResult[i])
		}
	}
}

var runFromResult = make([]string, 0)

type testFromTask struct {
	name string
}

func (t *testFromTask) Execute() error {
	runFromResult = append(runFromResult, t.name+" execute")
	return nil
}

func TestWorkflow_RunFrom(t *testing.T) {
	t.Parallel()

	runFromResult = make([]string, 0)

	w := NewWorkflow()
	w.RegisterTask("a", &testFromTask{name: "a"})
	w.RegisterTask("b", &testFromTask{name: "b"})
	w.RegisterTask("c", &testFromTask{name: "c"})
	w.RegisterTask("d", &testFromTask{name: "d"})
	err := w.RunFrom("b")

	if err != nil {
		t.Error(err)
	}
	tests := []string{"b execute", "c execute", "d execute"}
	if len(runFromResult) != len(tests) {
		t.Errorf("workflow: incorrect task result length expect:%v got:%v ", len(tests), len(runFromResult))
	}
	for i, test := range tests {
		if runFromResult[i] != test {
			t.Errorf("workflow: invalid task result expect:%v got: %v", test, runFromResult[i])
		}
	}

	err = w.RunFrom("unknown")
	if err == nil {
		t.Error("workflow: workflow not raises error when task not found")
	}
}

var runOnlyResult = make([]string, 0)

type testOnlyTask struct {
	name string
}

func (t *testOnlyTask) Execute() error {
	runOnlyResult = append(runOnlyResult, t.name+" execute")
	return nil
}

func TestWorkflow_RunOnly(t *testing.T) {
	t.Parallel()

	runOnlyResult = make([]string, 0)

	w := NewWorkflow()
	w.RegisterTask("a", &testOnlyTask{name: "a"})
	w.RegisterTask("b", &testOnlyTask{name: "b"})
	w.RegisterTask("c", &testOnlyTask{name: "c"})
	w.RegisterTask("d", &testOnlyTask{name: "d"})
	err := w.RunOnly("b")

	if err != nil {
		t.Error(err)
	}
	tests := []string{"b execute"}
	if len(runOnlyResult) != len(tests) {
		t.Errorf("workflow: incorrect task result length expect:%v got:%v ", len(tests), len(runOnlyResult))
	}
	for i, test := range tests {
		if runOnlyResult[i] != test {
			t.Errorf("workflow: invalid task result expect:%v got: %v", test, runOnlyResult[i])
		}
	}

	err = w.RunOnly("unknown")
	if err == nil {
		t.Error("workflow: workflow not raises error when task not found")
	}
}