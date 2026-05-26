package utils

import (
	"database/sql"
	f "fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var DB *sql.DB

func Connect() error {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
		return err
	}

	server := f.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	DB, err = sql.Open("postgres", server)

	err = DB.Ping()
	if err != nil {
		log.Println("Erro real da conexão:", err)
		return err
	}

	log.Println("Conexão com o banco realizada com sucesso!")
	f.Println("Conexão com o banco realizada com sucesso!")

	return nil
}
