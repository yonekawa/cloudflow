package aws

import (
	"testing"

	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
)

func TestBatchJobTask_Execute(t *testing.T) {
	t.Parallel()

	sess, err := session.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	testBatchJobTaskSucceeded(t, sess)
	testBatchJobTaskTimeout(t, sess)
}

func testBatchJobTaskSucceeded(t *testing.T, sess *session.Session) {
	jobID := "TESTING"
	submitJob = func(b *batch.Batch, input *batch.SubmitJobInput) (*batch.SubmitJobOutput, error) {
		return &batch.SubmitJobOutput{JobId: &jobID}, nil
	}
	status := "RUNNING"
	describeJobs = func(b *batch.Batch, input *batch.DescribeJobsInput) (*batch.DescribeJobsOutput, error) {
		st := status
		detail := &batch.JobDetail{
			JobId:        &jobID,
			Status:       &st,
			StatusReason: nil,
		}
		out := &batch.DescribeJobsOutput{Jobs: []*batch.JobDetail{detail}}
		status = "SUCCEEDED"
		return out, nil
	}

	bjt := NewBatchJobTask(sess, &batch.SubmitJobInput{})
	if err := bjt.Execute(); err != nil {
		t.Error(err)
	}
}

func testBatchJobTaskTimeout(t *testing.T, sess *session.Session) {
	jobID := "TESTING2"
	submitJob = func(b *batch.Batch, input *batch.SubmitJobInput) (*batch.SubmitJobOutput, error) {
		return &batch.SubmitJobOutput{JobId: &jobID}, nil
	}
	status := "RUNNING"
	describeJobs = func(b *batch.Batch, input *batch.DescribeJobsInput) (*batch.DescribeJobsOutput, error) {
		detail := &batch.JobDetail{
			JobId:        &jobID,
			Status:       &status,
			StatusReason: nil,
		}
		out := &batch.DescribeJobsOutput{Jobs: []*batch.JobDetail{detail}}
		return out, nil
	}

	bjt := NewBatchJobTask(sess, &batch.SubmitJobInput{})
	bjt.PollingTime = 10 * time.Microsecond
	bjt.Timeout = time.Second
	if err := bjt.Execute(); err == nil {
		t.Error("expect to occur timeout but it succeeded")
	}
}
