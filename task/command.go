package task

import "os/exec"

type CommandTask struct {
	name string
	args []string
}

func NewCommandTask(name string, args ...string) *CommandTask {
	return &CommandTask{name: name, args: args}
}

func (cmd *CommandTask) Execute() error {
	return exec.Command(cmd.name, cmd.args...).Run()
}
