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

	server := f.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DATABASE_NAME"))

	DB, err = sql.Open("postgres", server)

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Erro na conexão com o BANCO!")
		return err
	}
	log.Println(("Conexão com o banco realizada com sucesso!"))
	f.Println("Conexão com o banco realizada com sucesso!")
	return nil

}
