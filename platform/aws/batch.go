package aws

import (
	"time"

	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/batch"
)

var defaultTimeout = time.Minute
var defaultPollingTime = time.Second

// BatchJobTask execute AWS Batch Job.
type BatchJobTask struct {
	Session        *session.Session
	SubmitJobInput *batch.SubmitJobInput
	PollingTime    time.Duration
	Timeout        time.Duration
}

// NewBatchJobTask creates a AWS Batch Job task.
func NewBatchJobTask(session *session.Session, input *batch.SubmitJobInput) *BatchJobTask {
	return &BatchJobTask{
		Session:        session,
		SubmitJobInput: input,
		PollingTime:    defaultPollingTime,
		Timeout:        defaultTimeout,
	}
}

// Execute implement Task.Execute
func (bjt *BatchJobTask) Execute() error {
	completeChan := make(chan error)

	go func() {
		b := batch.New(bjt.Session)
		submit, err := submitJob(b, bjt.SubmitJobInput)
		if err != nil {
			completeChan <- err
		}

		elapsed := 0 * time.Millisecond

		for {
			describe, err := describeJobs(b, &batch.DescribeJobsInput{Jobs: []*string{submit.JobId}})
			if err != nil {
				completeChan <- err
			}

			if len(describe.Jobs) == 0 {
				completeChan <- fmt.Errorf("cloudflow: aws batch job:%v not found", submit.JobId)
				break
			}

			job := describe.Jobs[0]
			switch *job.Status {
			case "SUCCEEDED":
				completeChan <- nil
				break
			case "FAILED":
				completeChan <- fmt.Errorf("cloudflow: aws batch job id:%v failed by reason:%v", job.JobId, job.StatusReason)
				break
			default:
				time.Sleep(bjt.PollingTime)
				elapsed += bjt.PollingTime
			}

			if elapsed >= bjt.Timeout {
				completeChan <- fmt.Errorf("cloudflow: aws batch job id:%v timed out", job.JobId)
				break
			}
		}
	}()

	return <-completeChan
}

// for mock testing
var submitJob = func(b *batch.Batch, input *batch.SubmitJobInput) (*batch.SubmitJobOutput, error) {
	return b.SubmitJob(input)
}
var describeJobs = func(b *batch.Batch, input *batch.DescribeJobsInput) (*batch.DescribeJobsOutput, error) {
	return b.DescribeJobs(input)
}
