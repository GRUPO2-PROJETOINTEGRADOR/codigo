package utils

import (
    "log"
    "os"
)

// Criar_banco creates the database schema. It now only executes the schema SQL file.
// Seed data execution is optional and any errors are logged but not returned, to allow the server to start.
func Criar_banco() error {
    // Read schema SQL
    arqbytes, err := os.ReadFile("database/schema.sql")
    if err != nil {
        log.Fatalf("Erro na leitura do schema.sql: %v", err)
        return err
    }
    // Execute schema
    _, err = DB.Exec(string(arqbytes))
    if err != nil {
        log.Fatalf("Erro ao criar tabelas: %v", err)
        return err
    }
    // Attempt to load seed data (optional)
    seedBytes, err := os.ReadFile("database/seed.sql")
    if err != nil {
        log.Printf("Seed file not found or unreadable: %v", err)
        return nil
    }
    if _, err = DB.Exec(string(seedBytes)); err != nil {
        log.Printf("Erro ao executar seed data (ignored): %v", err)
        // Continue without failing
    }
    return nil
}
