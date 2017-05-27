# cloudflow
cloudflow is a workflow engine written in Go.
Designed to running with cloud computing platform.

# Installation

Install depends libraries.

```golang
go get github.com/hashicorp/go-multierror
```

Install cloudflow.

```console
go get github.com/yonekawa/cloudflow
```

# Usage

### Define workflow

```go
wf := NewWorkflow()
wf.RegisterTask("download", &DownloadTask{...})
wf.RegisterTask("process", &ProcessTask{...})
wf.RegisterParallelTask("process_parallel", NewParallelTask([]Task{&ProcessTask{}, &ProcessTask{}}))
wf.RegisterTask("output", &OutputTask{...})
```

### Run workflow

```go
wf.Run()
wf.RunFrom("process")
wf.RunOnly("output")
```

# License
This library is distributed under the MIT license found in the [LICENSE](https://github.com/yonekawa/cloudflow/blob/master/LICENSE) file.
