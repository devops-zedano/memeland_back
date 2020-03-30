package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// This code is to create a file.
	bucket := "a-memeland"
	file, err := ioutil.TempFile("/tmp", "generated_file-*.txt")

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	defer os.Remove(file.Name())

	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession())

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	if err != nil {
		return events.APIGatewayProxyResponse{}, fmt.Errorf("failed to open file %q, %v", file.Name(), err)
	}

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file.Name()),
		Body:   file,
	})

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// var objects []response
	// var objectsJSON []byte
	// bucket := "a-memeland"

	// // List the objects in the S3 bucket
	// sess, err := session.NewSession(&aws.Config{})

	// if err != nil {
	// 	return events.APIGatewayProxyResponse{}, err
	// }

	// svc := s3.New(sess)
	// resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})

	// if err != nil {
	// 	return events.APIGatewayProxyResponse{}, err
	// }

	// for _, item := range resp.Contents {
	// 	fmt.Println("Name:         ", *item.Key)
	// 	fmt.Println("Last modified:", *item.LastModified)
	// 	fmt.Println("Size:         ", *item.Size)
	// 	fmt.Println("Storage class:", *item.StorageClass)
	// 	fmt.Println("")

	// 	objects = append(objects, response{Name: *item.Key})
	// }

	// objectsJSON, err = json.Marshal(objects)

	// if err != nil {
	// 	return events.APIGatewayProxyResponse{}, err
	// }

	return events.APIGatewayProxyResponse{Body: "Success!", StatusCode: 200}, nil
}

// func uploadFile(session *session.Session, uploadFileDir string) error {

// 	upFile, err := os.Open(uploadFileDir)
// 	if err != nil {
// 		return err
// 	}
// 	defer upFile.Close()

// 	upFileInfo, _ := upFile.Stat()
// 	var fileSize int64 = upFileInfo.Size()
// 	fileBuffer := make([]byte, fileSize)
// 	upFile.Read(fileBuffer)

// 	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
// 		Bucket:               aws.String(AWS_S3_BUCKET),
// 		Key:                  aws.String(uploadFileDir),
// 		ACL:                  aws.String("private"),
// 		Body:                 bytes.NewReader(fileBuffer),
// 		ContentLength:        aws.Int64(fileSize),
// 		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
// 		ContentDisposition:   aws.String("attachment"),
// 		ServerSideEncryption: aws.String("AES256"),
// 	})
// 	return err
// }

func main() {
	lambda.Start(handleRequest)
}
