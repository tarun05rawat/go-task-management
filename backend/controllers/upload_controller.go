package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tarun05rawat/go-task-management/model"

	"github.com/gin-gonic/gin"
	"github.com/tarun05rawat/go-task-management/database"
	"github.com/tarun05rawat/go-task-management/services"
)

func UploadFiles(c *gin.Context) {
	userID := c.GetString("userID")
	taskID := c.Param("id")

	// Query the task and check ownership
	var task model.Task
	if err := database.DB.First(&task, "id = ? AND user_id = ?", taskID, userID).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this task"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form-data"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files provided"})
		return
	}

	var urls []string

	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer f.Close()

		// Construct S3 key (file path inside bucket)
		key := fmt.Sprintf("tasks/%s/%s", taskID, filepath.Base(file.Filename))

		// Upload to S3 using services.S3Client
		err = services.UploadToS3(f, file.Header.Get("Content-Type"), key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "S3 upload failed"})
			return
		}

		s3Url := fmt.Sprintf("s3://%s/%s", services.BucketName, key)
		urls = append(urls, s3Url)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded successfully",
		"files":   urls,
	})
}

func ListAttachments(c *gin.Context) {
	taskID := c.Param("id")

	// Optional: Check task ownership here like we did for uploads
	userID := c.GetString("userID")
	var task model.Task
	if err := database.DB.First(&task, "id = ? AND user_id = ?", taskID, userID).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this task"})
		return
	}

	// List all files under this task folder in S3
	prefix := fmt.Sprintf("tasks/%s/", taskID)
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(services.BucketName),
		Prefix: aws.String(prefix),
	}

	result, err := services.S3Client.ListObjectsV2(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list objects"})
		return
	}

	var signedUrls []string
	for _, item := range result.Contents {
		url, err := services.GeneratePreSignedURL(*item.Key)
		if err == nil {
			signedUrls = append(signedUrls, url)
		}
	}

	c.JSON(http.StatusOK, gin.H{"attachments": signedUrls})
}
