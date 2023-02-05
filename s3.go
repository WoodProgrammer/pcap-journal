package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func uploadS3(bucketName string, filePrefix string, fileName string) {

	fullName := filePrefix + "/" + fileName
	file, err := os.Open(fileName)

	if err != nil {
		exitErrorf("Unable to open file %q, %v", fileName, err)
	}

	defer file.Close()

	sess, err := session.NewSession()

	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fullName),

		Body: file,
	})

	if err != nil {
		exitErrorf("Unable to upload %q to %q, %v", fileName, bucketName, err)
	}

	fmt.Printf("Successfully uploaded %q to %q under %q directory \n", fileName, bucketName, filePrefix)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n")
	os.Exit(1)
}
