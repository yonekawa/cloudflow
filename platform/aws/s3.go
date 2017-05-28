package aws

import (
	"io/ioutil"
	"os"
	"path"

	"path/filepath"
	"sync"

	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/go-multierror"
)

// S3BulkUploadTask uploads local files in src dir into s3 dst folder.
type S3BulkUploadTask struct {
	Session     *session.Session
	SrcDir      string
	S3DstFolder string
	Bucket      string
}

// NewS3BulkUploadTask creates a s3 bulk upload task.
func NewS3BulkUploadTask(sess *session.Session, srcDir, s3DstDir, bucket string) *S3BulkUploadTask {
	return &S3BulkUploadTask{
		Session:     sess,
		SrcDir:      srcDir,
		S3DstFolder: s3DstDir,
		Bucket:      bucket,
	}
}

// Execute implement Task.Execute
func (up *S3BulkUploadTask) Execute() error {
	files, err := readFilesInDir(up.SrcDir)
	if err != nil {
		return err
	}

	svc := s3.New(up.Session)

	wg := sync.WaitGroup{}
	errChan := make(chan error)
	for _, info := range files {
		wg.Add(1)
		go func(info os.FileInfo) {
			defer wg.Done()

			srcFile := filepath.Join(up.SrcDir, info.Name())
			dstS3Key := path.Join(up.S3DstFolder, info.Name())
			file, err := os.Open(srcFile)
			if err != nil {
				errChan <- err
				return
			}
			defer file.Close()

			_, err = putObject(svc, &s3.PutObjectInput{
				Key:    aws.String(dstS3Key),
				Bucket: aws.String(up.Bucket),
				Body:   file,
			})
			if err != nil {
				errChan <- err
			}
		}(info)
	}

	resultChan := make(chan error)
	go func() {
		var result *multierror.Error
		for err := range errChan {
			result = multierror.Append(result, err)
		}
		resultChan <- result.ErrorOrNil()
	}()

	wg.Wait()
	close(errChan)

	return <-resultChan
}

func readFilesInDir(dir string) ([]os.FileInfo, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]os.FileInfo, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		files = append(files, e)
	}

	return files, nil
}

// for mock testing
var putObject = func(svc *s3.S3, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return svc.PutObject(input)
}

// S3BulkDownloadTask downloads files in s3 folder into local dst dir.
type S3BulkDownloadTask struct {
	Session     *session.Session
	S3SrcFolder string
	DstDir      string
	Bucket      string
}

// NewS3BulkDownloadTask creates a s3 bulk download task.
func NewS3BulkDownloadTask(sess *session.Session, s3SrcFolder, dstDir, bucket string) *S3BulkDownloadTask {
	return &S3BulkDownloadTask{
		Session:     sess,
		S3SrcFolder: s3SrcFolder,
		DstDir:      dstDir,
		Bucket:      bucket,
	}
}

// Execute implement Task.Execute.
func (down *S3BulkDownloadTask) Execute() error {
	svc := s3.New(down.Session)

	list, err := listObjectsV2(svc, &s3.ListObjectsV2Input{
		Bucket:    aws.String(down.Bucket),
		Delimiter: aws.String("/"),
		Prefix:    aws.String(down.S3SrcFolder + "/"),
	})
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	errChan := make(chan error)
	for _, c := range list.Contents {
		wg.Add(1)
		go func(c *s3.Object) {
			defer wg.Done()

			out, err := getObject(svc, &s3.GetObjectInput{
				Key:    c.Key,
				Bucket: aws.String(down.Bucket),
			})
			if err != nil {
				errChan <- err
				return
			}

			dstPath := filepath.Join(down.DstDir, path.Base(*c.Key))
			file, err := os.Create(dstPath)
			if err != nil {
				errChan <- err
				return
			}
			defer file.Close()

			_, err = io.Copy(file, out.Body)
			if err != nil {
				errChan <- err
			}
		}(c)
	}

	resultChan := make(chan error)
	go func() {
		var result *multierror.Error
		for err := range errChan {
			result = multierror.Append(result, err)
		}
		resultChan <- result.ErrorOrNil()
	}()

	wg.Wait()
	close(errChan)

	return <-resultChan
}

// for mock testing
var getObject = func(svc *s3.S3, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	return svc.GetObject(input)
}
var listObjectsV2 = func(svc *s3.S3, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return svc.ListObjectsV2(input)
}
