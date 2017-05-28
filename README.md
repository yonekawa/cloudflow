# cloudflow
[![Build Status](https://travis-ci.org/yonekawa/cloudflow.svg?branch=master)](https://travis-ci.org/yonekawa/cloudflow)
[![GoDoc](https://godoc.org/github.com/yonekawa/cloudflow?status.svg)](http://godoc.org/github.com/yonekawa/cloudflow)
[![Go Report Card](https://goreportcard.com/badge/github.com/yonekawa/cloudflow)](https://goreportcard.com/report/github.com/yonekawa/cloudflow)

# Description
cloudflow is a workflow engine written in Go.
Designed to running with cloud computing platform.

# Installation

Install depends libraries.

```golang
go get github.com/aws/aws-sdk-go
go get github.com/hashicorp/go-multierror
```

Install cloudflow.

```console
go get github.com/yonekawa/cloudflow
```

# Usage

### Define & Run workflow

```go
import "github.com/yonekawa/cloudflow"

wf := cloudflow.NewWorkflow()
wf.AddTask("download", &DownloadTask{...})
wf.AddTask("process", &ProcessTask{...})

// parallel execution
pt := NewParallelTask()
pt.AddTask("parallel-1", &ProcessTask{})
pt.AddTask("parallel-2", &ProcessTask{})
wf.AddTask("parallel", pt)

wf.AddTask("output", &OutputTask{...})

// Show task summary
fmt.Print(wf.Summary())

wf.Run()
wf.RunFrom("process")
wf.RunOnly("output")
```

# Builtin tasks

### task.CommandTask

`CommandTask` executes local command by `exec.Command`.

```go
import "github.com/yonekawa/cloudflow/task"

cmd := task.NewCommandTask("go", "help", "build")
err := cmd.Execute()
```

### aws.S3BulkUploadTask & aws.S3BulkDownloadTask

`aws.S3BulkUploadTask` uploads local files in src dir into S3 dst folder.

```go
import "github.com/aws/aws-sdk-go/session"
import "github.com/yonekawa/cloudflow/platform/aws"

sess := session.Must(session.NewSession())
// Upload ./src files into s3:/s3-bucket/s3dst/
task := aws.NewS3BulkUploadTask(sess, "./src", "/s3dst", "s3-bucket")
err := task.Execute()
```

`aws.S3BulkDownloadTask` downloads files in S3 folder into local dst dir.

```go
import "github.com/aws/aws-sdk-go/session"
import "github.com/yonekawa/cloudflow/platform/aws"

sess := session.Must(session.NewSession())
// Download s3:/s3-bucket/s3dst/ files into ./dst
task := aws.NewS3BulkUploadTask(sess, "/s3src", "./dst", "s3-bucket")
err := task.Execute()
```

### aws.BatchJobTask

`aws.BatchJobTask` submit [AWS Batch](https://aws.amazon.com/jp/documentation/batch/) Job and wait to complete a job.

```go
import "github.com/aws/aws-sdk-go/session"
import "github.com/aws/aws-sdk-go/service/batch"
import "github.com/yonekawa/cloudflow/platform/aws"

sess := session.Must(session.NewSession())
task := aws.NewBatchJobTask(sess, &batch.SubmitJobInput{
  JobDefinition: aws.String("job definition ARN"),
  JobQueue:      aws.String("job queue ARN"),
  JobName:       aws.String("job name"),
})
err := task.Execute()
```

You can change polling interval and timeout.

```go
task := aws.NewBatchJobTask(sess, &batch.SubmitJobInput{
  JobDefinition: aws.String("job definition ARN"),
  JobQueue:      aws.String("job queue ARN"),
  JobName:       aws.String("job name"),
})
task.PollingTime = 5 * time.Minute
task.Timeout = 10 * time.Hour
```

### aws.LambdaInvokeTask

`aws.LambdaInvokeTask` invokes lambda function.

```go
import awssdk "github.com/aws/aws-sdk-go/aws"
import "github.com/aws/aws-sdk-go/session"
import "github.com/aws/aws-sdk-go/service/lambda"
import "github.com/yonekawa/cloudflow/platform/aws"

sess := session.Must(session.NewSession())
task := aws.NewLambdaInvokeTask(sess, &batch.SubmitJobInput{
  FunctionName: awssdk.String("function ARN"),
})
err := task.Execute()
```

# License
This library is distributed under the MIT license found in the [LICENSE](https://github.com/yonekawa/cloudflow/blob/master/LICENSE) file.
