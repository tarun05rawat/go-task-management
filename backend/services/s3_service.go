package services

import (
	"mime/multipart"
	"os"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var S3Client *s3.S3

func InitS3() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	}))
	S3Client = s3.New(sess)
}

var BucketName = os.Getenv("S3_BUCKET_NAME")

func UploadToS3(file multipart.File, contentType, key string) error {
	_, err := S3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(BucketName),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(contentType),
		ACL:         aws.String("private"), // Keeps file private
	})
	return err
}

func GeneratePreSignedURL(key string) (string, error) {
	req, _ := S3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(key),
	})
	return req.Presign(15 * time.Minute) // Expires in 15 minutes
}
