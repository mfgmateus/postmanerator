package themes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// TODO fill these in!
const (
	S3Region = "sa-east-1"
	S3Bucket = "developer-docs.monnos.com"
)

func AddFileToS3(s *session.Session, key string, filename string) string {

	// Open the file for use
	file, _ := os.Open(filename)
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(S3Bucket),
		Key:                  aws.String(key),
		ACL:                  aws.String("public-read"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	print(err)

	url := "https://s3-%s.amazonaws.com/%s/%s"
	url = fmt.Sprintf(url, S3Region, S3Bucket, key)
	return url
}

func helperUploadToS3(input string, key string) string {

	key = strings.Trim(key, " ")
	key = strings.ToLower(key)
	key = strings.Replace(key, " ", "-", -1)
	key = strings.Replace(key, "(", "", -1)
	key = strings.Replace(key, ")", "", -1)

	key += ".json"

	// Create a single AWS session (we can re use this if we're uploading many files)
	s, err := session.NewSession(&aws.Config{Region: aws.String(S3Region)})
	if err != nil {
		log.Fatal(err)
	}

	d1 := []byte(input)
	filename := "/tmp/" + key
	err = ioutil.WriteFile(filename, d1, 0644)

	// Upload
	return AddFileToS3(s, key, filename)
}
