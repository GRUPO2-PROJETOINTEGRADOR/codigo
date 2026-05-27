package utils

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
    // Load environment variables from .env file (if present)
    if err := godotenv.Load(); err != nil {
        log.Printf("Erro ao carregar .env: %v", err)
    }

    // Retrieve required DB connection parameters
    host := os.Getenv("DB_HOST")
    port := os.Getenv("PORT")
    user := os.Getenv("DB_USER")
    pass := os.Getenv("DB_PASS")
    dbname := os.Getenv("DATABASE_NAME")

    // Validate required parameters
    if host == "" || port == "" || user == "" || pass == "" || dbname == "" {
        return fmt.Errorf("missing required DB environment variables")
    }

    // Build connection string
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)

    // Open the database connection
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Printf("Erro ao abrir conexão com o banco: %v", err)
        return err
    }

    // Verify the connection is alive
    if err = db.Ping(); err != nil {
        log.Printf("Erro ao conectar ao banco: %v", err)
        return err
    }

    // Assign to package-level variable after successful connection
    DB = db
    log.Println("Conexão com o banco realizada com sucesso!")
    fmt.Println("Conexão com o banco realizada com sucesso!")
    return nil
}
