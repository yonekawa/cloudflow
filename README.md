# cloudflow
[![Build Status](https://travis-ci.org/yonekawa/cloudflow.svg?branch=master)](https://travis-ci.org/yonekawa/cloudflow)

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

### Define workflow

```go
wf := cloudflow.NewWorkflow()
wf.AddTask("download", &DownloadTask{...})
wf.AddTask("process", &ProcessTask{...})
wf.AddTask("parallel", NewParallelTask([]Task{&ProcessTask{}, &ProcessTask{}}))
wf.AddTask("output", &OutputTask{...})
```

### Run workflow

```go
wf.Run()
wf.RunFrom("process")
wf.RunOnly("output")
```

# License
This library is distributed under the MIT license found in the [LICENSE](https://github.com/yonekawa/cloudflow/blob/master/LICENSE) file.
