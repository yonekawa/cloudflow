package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// LambdaInvokeTask invokes lambda function.
type LambdaInvokeTask struct {
	Session     *session.Session
	InvokeInput *lambda.InvokeInput
}

// NewLambdaInvokeTask creates a lambda invoke task.
func NewLambdaInvokeTask(sess *session.Session, input *lambda.InvokeInput) *LambdaInvokeTask {
	return &LambdaInvokeTask{
		Session:     sess,
		InvokeInput: input,
	}
}

// Execute implement Task.Execute.
func (li *LambdaInvokeTask) Execute() error {
	f := lambda.New(li.Session)

	_, err := invoke(f, li.InvokeInput)
	if err != nil {
		return err
	}

	// TODO: Notify invoke result

	return nil
}

// for mock testing
var invoke = func(f *lambda.Lambda, input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
	return f.Invoke(input)
}
