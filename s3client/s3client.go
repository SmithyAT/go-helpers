package s3client

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"path/filepath"
	"strings"
)

// func S3Upload(sFile, dFile, endpoint, accessKeyID, secretAccessKey, bucketName, logid string) error {
// 	client, err := NewS3Client(endpoint, accessKeyID, secretAccessKey)
// 	if err != nil {
// 		return err
// 	}
//
// 	info, err := client.UploadFile(sFile, dFile, bucketName)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

type S3Client struct {
	MinioClient *minio.Client
	BucketName  string
}

func NewS3Client(endpoint, accessKeyID, secretAccessKey string) (*S3Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}
	return &S3Client{MinioClient: minioClient}, nil
}

func (s *S3Client) UploadFile(sFile, dFile, bucketName string) (minio.UploadInfo, error) {
	contentType := getContentType(sFile)
	info, err := s.MinioClient.FPutObject(context.Background(), bucketName, dFile, sFile, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return info, err
	}
	return info, nil
}

func getContentType(sFile string) string {
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
	return contentType
}
