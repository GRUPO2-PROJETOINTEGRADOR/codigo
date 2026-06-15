package repo

import (
	"context"
	"mime/multipart"
	"path"
	"strings"

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

func RemoverPDF(anexo string) error {
	if anexo == "" {
		return nil
	}

	// Extrai o nome do arquivo da URL completa
	fileName := path.Base(anexo)
	if fileName == "." || fileName == "/" || strings.TrimSpace(fileName) == "" {
		return nil
	}

	err := MinioClient.RemoveObject(
		context.Background(),
		"auditorias",
		fileName,
		minio.RemoveObjectOptions{},
	)
	return err
}
