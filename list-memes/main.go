package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type response struct {
	Name string `json:"name"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// This code is to create a file.
	// file, err := ioutil.TempFile("dir", "prefix")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer os.Remove(file.Name())

	// // The session the S3 Uploader will use
	// sess := session.Must(session.NewSession())

	// // Create an uploader with the session and default options
	// uploader := s3manager.NewUploader(sess)

	// f, err := os.Open(file.Name())
	// if err != nil {
	// 	return fmt.Errorf("failed to open file %q, %v", file.Name(), err)
	// }

	var objects []response
	var objectsJSON []byte
	bucket := "a-memeland"

	// List the objects in the S3 bucket
	sess, err := session.NewSession(&aws.Config{})

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")

		objects = append(objects, response{Name: *item.Key})
	}

	objectsJSON, err = json.Marshal(objects)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{Body: string(objectsJSON), StatusCode: 200}, nil
}

func errorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

func main() {
	lambda.Start(handleRequest)
}
