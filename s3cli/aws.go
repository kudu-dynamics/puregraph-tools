package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

func getS3() (*s3.S3, *session.Session) {
	endpoint, ok := os.LookupEnv("AWS_ENDPOINT_URL")
	if !ok {
        # XXX change endpoint
		endpoint = "https://<endpoint>"
	}

	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		region = "us-east-1"
	}

	_, insecure := os.LookupEnv("INSECURE_SKIP_VERIFY")

	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
		DisableSSL: aws.Bool(false),
		Endpoint:   aws.String(endpoint),
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecure,
				},
			},
		},
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession := session.Must(session.NewSession(s3Config))
	s3Client := s3.New(newSession)
	return s3Client, newSession
}

func downloadObject(bucket, object, filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filepath": filepath,
		}).Fatal("unable to open file")
	}
	defer file.Close()

	_, sess := getS3()
	downloader := s3manager.NewDownloader(sess)

	log.WithFields(log.Fields{
		"bucket": bucket,
		"object": object,
	}).Debug("download started...")

	_, err = downloader.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(object),
		},
	)
	if err != nil {
		log.WithFields(log.Fields{
			"bucket": bucket,
			"error":  err,
			"object": object,
		}).Error("download failed...")

		os.Remove(filepath)

		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				// Set a specific error code when a file is not found.
				os.Exit(int(syscall.ENOENT))
			}
		}

		os.Exit(1)
	}

	log.WithFields(log.Fields{
		"bucket": bucket,
		"object": object,
	}).Debug("download complete...")
}

func uploadObject(bucket, object, filepath string) {
	_, sess := getS3()
	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(filepath)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"filepath": filepath,
		}).Fatal("unable to open file")
	}

	log.WithFields(log.Fields{
		"bucket": bucket,
		"object": object,
	}).Debug("upload started...")

	_, err = uploader.Upload(
		&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(object),
			Body:   file,
		},
	)
	if err != nil {
		log.WithFields(log.Fields{
			"bucket": bucket,
			"error":  err,
			"object": object,
		}).Fatal("upload failed...")
	}

	log.WithFields(log.Fields{
		"bucket": bucket,
		"object": object,
	}).Debug("upload complete...")
}
