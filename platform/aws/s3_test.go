package aws

import (
	"testing"

	"errors"
	"io/ioutil"
	"strings"

	"path/filepath"

	"path"

	"strconv"

	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestS3BulkUploadTask_Execute(t *testing.T) {
	t.Parallel()

	uploadFiles := make([]string, 0)
	putObject = func(svc *s3.S3, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
		uploadFiles = append(uploadFiles, *input.Key)
		if strings.HasSuffix(*input.Key, "_error") {
			return nil, errors.New("error")
		}
		return &s3.PutObjectOutput{}, nil
	}

	srcDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	srcFiles := [3]string{"file1", "file2", "file3"}
	for i, f := range srcFiles {
		ioutil.WriteFile(filepath.Join(srcDir, f), []byte(strconv.Itoa(i)), 0666)
	}

	sess, err := session.NewSession()
	if err != nil {
		t.Fatal(err)
	}

	task := NewS3BulkUploadTask(sess, srcDir, "/dst", "file-bucket")
	if err := task.Execute(); err != nil {
		t.Error(err)
	}

	for _, f := range srcFiles {
		found := false
		for _, u := range uploadFiles {
			if u == path.Join(task.S3DstFolder, f) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("src file:%v does not uploaded", f)
		}
	}

	ioutil.WriteFile(filepath.Join(srcDir, "file_error"), []byte("error"), 0666)
	if err := task.Execute(); err == nil {
		t.Error("expect to fail upload but it succeeded")
	}
}

func TestS3BulkDownloadTask_Execute(t *testing.T) {
	t.Parallel()

	dstDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	srcDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	readClosers := make(map[string]io.ReadCloser)
	srcFiles := [3]string{"file1", "file2", "file3"}
	for i, f := range srcFiles {
		p := filepath.Join(srcDir, f)
		ioutil.WriteFile(p, []byte(strconv.Itoa(i)), 0666)
		readCloser, err := os.Open(p)
		if err != nil {
			t.Fatal(err)
		}
		readClosers[f] = readCloser
	}
	defer func() {
		for _, f := range readClosers {
			f.Close()
		}
	}()

	sess, err := session.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	task := NewS3BulkDownloadTask(sess, "/s3src", dstDir, "file-bucket")

	listObjectsV2 = func(svc *s3.S3, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
		entries, err := ioutil.ReadDir(srcDir)
		if err != nil {
			return nil, err
		}

		contents := make([]*s3.Object, 0, len(entries))
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			contents = append(contents, &s3.Object{Key: aws.String(path.Join(task.S3SrcFolder, e.Name()))})
		}
		return &s3.ListObjectsV2Output{Contents: contents}, nil
	}
	getObject = func(svc *s3.S3, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
		if strings.HasSuffix(*input.Key, "_error") {
			return nil, errors.New("error")
		}

		file := path.Base(*input.Key)
		if body := readClosers[file]; body != nil {
			return &s3.GetObjectOutput{Body: body}, nil
		}
		return nil, errors.New("file not found")
	}

	if err := task.Execute(); err != nil {
		t.Error(err)
	}

	ioutil.WriteFile(filepath.Join(srcDir, "file_error"), []byte("error"), 0666)
	if err := task.Execute(); err == nil {
		t.Error("expect to fail upload but it succeeded")
	}
}
