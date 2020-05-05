package services

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"mime/multipart"
	"os"
	"time"
)

type CloudInterface interface {
	UploadLocalFile(filename string, key string) error
	UploadMultipartFile(file multipart.File, key string) error
}

type Cloud struct {
	AwsSession *session.Session
	AwsBucket  string
}

// Idea here is to add a facory, so we can return/support multiple cloud solutions (AWS, GCP, etc)
var NewCloud = func() (*Cloud, error) {
	cloud := &Cloud{}

	return cloud, nil
}

//
func (c *Cloud) InitAws() error {
	_, err := c.GetAwsBucket()

	if err != nil {
		return err
	}

	_, err = c.GetAwsSession()

	return err
}

//
func (c *Cloud) GetAwsBucket() (string, error) {
	var err error

	if c.AwsBucket != "" {
		return c.AwsBucket, nil
	}

	c.AwsBucket = os.Getenv("UXT_BUCKET")

	if c.AwsBucket == "" {
		err := errors.New(fmt.Sprint("UXT_BUCKET environment varriable not set"))
		return "", err
	}

	return c.AwsBucket, err
}

// retrieve aws session
func (c *Cloud) GetAwsSession() (*session.Session, error) {
	var err error

	if c.AwsSession != nil {
		return c.AwsSession, nil
	}

	region := os.Getenv("AWS_REGION")

	if region == "" {
		err = errors.New(fmt.Sprint("AWS_REGION environment varriable not set"))
		return nil, err
	}

	c.AwsSession, err = session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	return c.AwsSession, err
}

//
func (c *Cloud) UploadLocalFile(filename string, key string) error {
	err := c.InitAws()

	if err != nil {
		return err
	}

	file, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	uploader := s3manager.NewUploader(c.AwsSession)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.AwsBucket),
		Key:    aws.String(key),
		Body:   file,
	})

	return err
}

//
func (c *Cloud) UploadReader(reader io.Reader, key string) error {
	err := c.InitAws()

	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(c.AwsSession)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.AwsBucket),
		Key:    aws.String(key),
		Body:   reader,
	})

	return err
}

//
func (c *Cloud) UploadMultipartFile(file multipart.File, key string) error {
	err := c.InitAws()

	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(c.AwsSession)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.AwsBucket),
		Key:    aws.String(key),
		Body:   file,
	})

	return err
}

//
func (c *Cloud) GetResourceUrl(key string, minutes time.Duration) (string, error) {
	err := c.InitAws()

	if err != nil {
		return "", err
	}

	svc := s3.New(c.AwsSession)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.AwsBucket),
		Key:    aws.String(key),
	})

	urlStr, err := req.Presign(minutes * time.Minute)

	if err != nil {
		return "", err
	}

	return urlStr, nil
}

//
func (c *Cloud) FileExists(key string) bool {
	err = c.InitAws()

	if err != nil {
		return err
	}

	svc := s3.New(c.AwsSession)

	input := &s3.HeadObjectInput{
		Bucket: aws.String(c.AwsBucket),
		Key:    aws.String(key),
	}

	result, err := svc.HeadObject(input)

	if err != nil {
		return false
	}

	return true
}
