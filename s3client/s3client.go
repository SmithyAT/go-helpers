package s3client

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"path/filepath"
	"strings"
)

// S3Upload is a function that uploads a file from source to destination to an S3 bucket.
// It requires the source and destination file path, endpoint, accessKeyID, secretAccessKey, and bucketName as parameters.
// It determines content type based on source file extension and uses MinIO client for handling S3 operations.
// In case of a successful upload, it returns the uploaded file size and nil as error.
// In case of an error, it returns zero as file size and the error that occurred.
func S3Upload(sFile, dFile, endpoint, accessKeyID, secretAccessKey, bucketName string) (fileSize int64, retErr error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return 0, err
	}

	// Set the content type
	var contentType string
	switch strings.ToLower(filepath.Ext(sFile)) {
	case ".gz":
		contentType = "application/x-gzip"
	case ".tar":
		contentType = "application/x-tar"
	case ".txt", ".dat", ".log":
		contentType = "text/plain"
	case ".csv":
		contentType = "text/csv"
	default:
		contentType = "application/octet-stream" // default case if none of the extensions match
	}

	// Upload file
	info, err := minioClient.FPutObject(context.Background(), bucketName, dFile, sFile, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return 0, err
	}

	return info.Size, nil
}
