package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3session *s3.S3
)

const (
	BUCKET_NAME  = "firsttry-bucket"
	BUCKET_NAME2 = "secondtryunique123-bucket"
	REGION       = "ap-southeast-1"
)

func init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})))
}

//func for see listbucket
func listBuckets() (resp *s3.ListBucketsOutput) {
	resp, err := s3session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		panic(err)
	}

	return resp
}

func listObjects(bucketName string) (resp *s3.ListObjectsV2Output) {
	resp, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		panic(err)
	}

	return resp
}

func createBucket(bucketName string) (resp *s3.CreateBucketOutput) {
	resp, err := s3session.CreateBucket(&s3.CreateBucketInput{
		// ACL: aws.String(s3.BucketCannedACLPrivate),
		// ACL: aws.String(s3.BucketCannedACLPublicRead),
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(REGION),
		},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				fmt.Println("Bucket name already in use!")
				panic(err)
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				fmt.Println("Bucket exists and is owned by you!")
			default:
				panic(err)
			}
		}
	}

	return resp
}

func downloadObject(filename, bucketName string) {
	fmt.Println("Downloading: ", filename)

	resp, err := s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})

	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile("downloaded_file/"+filename, body, 0644)
	if err != nil {
		panic(err)
	}
}

func uploadObject(filename string, bucketName string) (resp *s3.PutObjectOutput) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println("Uploading:", filename)
	resp, err = s3session.PutObject(&s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String(bucketName),
		Key:    aws.String(strings.Split(filename, "/")[1]),
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
	})

	if err != nil {
		panic(err)
	}

	return resp
}

func deleteObject(filename string, bucketName string) (resp *s3.DeleteObjectOutput) {
	fmt.Println("Deleting: ", filename)
	resp, err := s3session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})

	if err != nil {
		panic(err)
	}

	return resp
}

func main() {
	// folder := "upload_file"
	//  print list our bucket
	fmt.Println(listBuckets())

	// print list object  in specific bucket
	fmt.Println(listObjects(BUCKET_NAME))

	//createbucket
	// fmt.Println(createBucket(BUCKET_NAME2))

	//download object data
	// downloadObject("Task New WOB.txt", BUCKET_NAME)

	//upload file to aws object
	// uploadObject("upload_file/logo.png", BUCKET_NAME)

	//delete object in bucket aws
	// deleteObject("logo.png", BUCKET_NAME)

	//try upload all file in one folder
	// files, _ := ioutil.ReadDir(folder)
	// fmt.Println(files)
	// for _, file := range files {
	// 	if file.IsDir() {
	// 		continue
	// 	} else {
	// 		uploadObject(folder+"/"+file.Name(), BUCKET_NAME)
	// 	}
	// }

	// fmt.Println(listObjects(BUCKET_NAME))

	// for _, object := range listObjects(BUCKET_NAME).Contents {
	// 	downloadObject(*object.Key, BUCKET_NAME)
	// 	deleteObject(*object.Key, BUCKET_NAME)
	// }

	// fmt.Println(listObjects(BUCKET_NAME))

}
