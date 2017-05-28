package aws

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func TestLambdaInvokeTask_Execute(t *testing.T) {
	t.Parallel()

	invoke = func(f *lambda.Lambda, input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
		return &lambda.InvokeOutput{StatusCode: aws.Int64(200)}, nil
	}

	sess, err := session.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	task := NewLambdaInvokeTask(sess, &lambda.InvokeInput{FunctionName: aws.String(" arn:aws:lambda:us-east-1:000000000000:function:MyFunc")})
	if err := task.Execute(); err != nil {
		t.Error(err)
	}

	invoke = func(f *lambda.Lambda, input *lambda.InvokeInput) (*lambda.InvokeOutput, error) {
		return &lambda.InvokeOutput{StatusCode: aws.Int64(500)}, errors.New("error")
	}

	if err := task.Execute(); err == nil {
		t.Error("expect to fail task but it succeeded")
	}
}
