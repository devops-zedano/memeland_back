package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/lithammer/shortuuid/v3"
)

var bucket string
var prefix string

type response struct {
	JSONURL  string `json:"jsonUrl"`
	ImageURL string `json:"imageUrl"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bucket := "a-memeland"
	shortUUID := shortuuid.New()
	imageName := fmt.Sprintf("public/memes/%s/%s", shortUUID, request.QueryStringParameters["filename"])
	jsonName := fmt.Sprintf("%s.json", imageName)

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{})

	// Create S3 service client
	svc := s3.New(sess)

	// Get Image Presigned URL
	req1, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(imageName),
	})
	presignedImageURL, err := req1.Presign(5 * time.Minute)

	log.Println("The URL is:", presignedImageURL, " err:", err)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Get JSON Metadata Presigned URL
	req2, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(jsonName),
	})
	presignedJSONURL, err := req2.Presign(5 * time.Minute)

	log.Println("The URL is:", presignedJSONURL, " err:", err)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	theResponse := response{JSONURL: presignedJSONURL, ImageURL: presignedImageURL}
	marshalledBody, err := json.Marshal(theResponse)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(marshalledBody),
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
