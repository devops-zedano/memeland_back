package main

import (
	"encoding/json"
	"fmt"
	"memeland/model"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var bucket, memesPrefix string

type response struct {
	URI         string `json:"uri"`
	Description string `json:"description"`
}

func handleRequest() (events.APIGatewayProxyResponse, error) {
	var objects []response
	var objectsJSON []byte
	bucket = "a-memeland"
	memesPrefix = "public/memes/"

	// List the objects in the S3 bucket
	sess, err := session.NewSession(&aws.Config{})

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	svc := s3.New(sess)
	downloader := s3manager.NewDownloader(sess)

	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(memesPrefix),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")

		// Only add files, which don't have an ending slash.
		if !strings.HasSuffix(*item.Key, "/") {
			if strings.HasSuffix(*item.Key, ".json") {
				imageURI := strings.ReplaceAll(*item.Key, ".json", "")
				metadata, err := getJSONMetadata(*item.Key, downloader)

				if err != nil {
					return events.APIGatewayProxyResponse{}, err
				}

				objects = append(objects, response{
					URI:         imageURI,
					Description: metadata.Description,
				})
			}
		}
	}

	objectsJSON, err = json.Marshal(objects)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(objectsJSON),
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func getJSONMetadata(fileKey string, downloader *s3manager.Downloader) (model.Metadata, error) {
	var buf []byte
	var metadata model.Metadata
	awsBufferWriter := aws.NewWriteAtBuffer(buf)

	_, err := downloader.Download(awsBufferWriter, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	})

	if err != nil {
		return model.Metadata{}, nil
	}

	json.Unmarshal(awsBufferWriter.Bytes(), &metadata)

	return metadata, nil
}

func errorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

func main() {
	lambda.Start(handleRequest)
}
