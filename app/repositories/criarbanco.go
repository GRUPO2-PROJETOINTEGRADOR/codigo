package utils

import (
	"log"
	"os"
)

func Criar_banco() error {

	arqbytes, err := os.ReadFile("database/schema.sql")
	if err != nil {
		log.Fatalf("Erro na leitura arquivo, err: %e", err)
		return err
	}

	_, err = DB.Exec(string(arqbytes))
	if err != nil {
		log.Fatalln("Erro na crição de tabelas SQL")
		return err
	}

	seed, err := os.ReadFile("database/seed.sql")
	if err != nil {
		log.Fatalf("Erro na leitura arquivo, err: %e", err)
		return err
	}
	_, err = DB.Exec(string(seed))
	if err != nil {
		log.Println("Erro na seed")
		return nil
	}

	return nil
}
