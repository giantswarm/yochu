// Package s3 provides a client for accessing S3.
package s3

import (
	"os"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/juju/errgo"

	"github.com/giantswarm/yochu/fetchclient"
)

var vLogger = func(f string, v ...interface{}) {}

// Configure sets the logger for this package.
func Configure(vl func(f string, v ...interface{})) {
	vLogger = vl
}

// S3Client is an S3 client configured with a bucket for all requests to use.
type S3Client struct {
	bucket *s3.Bucket
}

// NewS3Client creates a new configured instance of s3Client.
// You can either pass the AWS credentials as arguments to the function or
// pass empty strings. In this case environment variables will be used for
// the credentials. Supported environment variables are AWS_ACCESS_KEY_ID,
// AWS_SECRET_ACCESS_KEY and S3_ENDPOINT.
func NewS3Client(awsAccessKey, awsSecretKey, s3Endpoint, bucket string) (fetchclient.FetchClient, error) {
	if awsAccessKey == "" {
		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	}

	if awsSecretKey == "" {
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}

	if s3Endpoint == "" {
		s3Endpoint = os.Getenv("S3_ENDPOINT")
	}

	if awsAccessKey == "" {
		return nil, errgo.New("AWS_ACCESS_KEY_ID or flag is empty")
	}

	if awsSecretKey == "" {
		return nil, errgo.New("AWS_SECRET_ACCESS_KEY or flag is empty")
	}

	if s3Endpoint == "" {
		return nil, errgo.New("S3_ENDPOINT not found in environment or flag is empty")
	}

	s3 := s3.New(
		aws.Auth{
			AccessKey: awsAccessKey,
			SecretKey: awsSecretKey,
		},
		aws.Region{
			S3Endpoint: "https://" + s3Endpoint,
		},
	)

	s3c := &S3Client{
		bucket: s3.Bucket(bucket),
	}

	return s3c, nil
}

// Get fetches the contents of the file described by the key in the S3Client's bucket,
// and returns them.
func (s3c *S3Client) Get(key string) ([]byte, error) {
	vLogger("  call s3Client.Get(bucket, key): %v - %v", s3c.bucket, key)

	raw, err := s3c.bucket.Get(key)
	if err != nil {
		return nil, mask(err)
	}

	return raw, nil
}
