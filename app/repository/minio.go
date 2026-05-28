package repo

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func ConnectMinio() {
	client, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})

	if err != nil {
		log.Fatal(err)
	}

	MinioClient = client

	exists, err := client.BucketExists(context.Background(), "auditorias")
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		log.Println("MinIO conectado com sucesso! Bucket auditorias encontrado.")
	}
}
