package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var testBucket string = "mysuperstronktestbuck"

func TestLoadToBucket(t *testing.T) {

}

func TestCreateBucket(t *testing.T) {
	CreateBucket()
}

func TestGetBucketObject(t *testing.T) {

	curr, err := os.Getwd()
	fmt.Println(curr)
	mockFile, err := os.Create("test.txt")
	data := []byte("this is test files content")
	mockFile.Write(data)

	if err != nil {
		log.Println(err)
	}

	defer mockFile.Close()
	defer os.Remove(mockFile.Name())

	err = LoadToBucket(BucketName, mockFile.Name(), data)

	// Call get bucket method
	if err != nil {
		t.Errorf(" expected no error but go %s", err)
	}

}
