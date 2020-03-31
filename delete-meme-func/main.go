package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var bucket string

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bucket = "a-memeland"
	memeToDelete := request.QueryStringParameters["uri"]
	splitMemeURI := strings.Split(memeToDelete, "/")
	memeFilename := splitMemeURI[len(splitMemeURI)-1]
	memeDirToDelete := strings.ReplaceAll(memeToDelete, memeFilename, "")

	fmt.Printf("Requested to delete meme: %s", memeToDelete)

	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	svc := s3.New(sess)

	fmt.Printf("Deleting directory contents %s...", memeDirToDelete)
	deleteObjectsInBucket(bucket, memeDirToDelete, svc)

	fmt.Printf("Deleting directory %s ...", memeDirToDelete)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(memeDirToDelete),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(memeDirToDelete),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func deleteObjectsInBucket(bucket, dirKey string, svc *s3.S3) error {
	var deleteObject s3.Delete
	objectOutput, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(dirKey),
	})

	if err != nil {
		return err
	}

	for _, item := range objectOutput.Contents {
		identifier := &s3.ObjectIdentifier{
			Key: item.Key,
		}
		deleteObject.Objects = append(deleteObject.Objects, identifier)
	}

	_, err = svc.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &deleteObject,
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
