package repo

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

func UploadPDF(file multipart.File, fileName string, fileSize int64) (string, error) {

	_, err := MinioClient.PutObject(
		context.Background(),
		"auditorias",
		fileName,
		file,
		fileSize,
		minio.PutObjectOptions{
			ContentType: "application/pdf",
		},
	)

	if err != nil {
		return "", err
	}

	url := "http://localhost:9000/auditorias/" + fileName

	return url, nil
}
