package main

import (
	"fmt"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func testoss() {
	// Endpoint, Access Key ID, and Access Key Secret
	endpoint := "http://oss-eu-central-1.aliyuncs.com"
	accessKeyId := "LTAI5tCvfHaENG5F2x1vyEjL"
	accessKeySecret := "tFKSQORQSSw2q9qVpkieBX91KJdkHB"

	// Create a client
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Bucket name and object key
	bucketName := "dmo0001"
	objectKey := "01.txt"

	// Get the bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Generate a presigned URL for GET request (downloading)
	expiration := time.Now().Add(15 * time.Minute) // URL will expire in 15 minutes

	expirationUnix := expiration.Unix()

	url, err := bucket.SignURL(objectKey, oss.HTTPGet, expirationUnix)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the URL
	fmt.Println("Presigned URL:", url)
}
