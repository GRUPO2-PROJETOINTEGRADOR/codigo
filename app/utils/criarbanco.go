package utils

import (
	"log"
	"os"
)

func Criar_banco() error {

	arqbytes, err := os.ReadFile("database/tabelas.sql")
	if err != nil {
		log.Fatalf("Erro na leitura arquivo, err: %e", err)
		return err
	}

	_, err = DB.Exec(string(arqbytes))
	if err != nil {
		log.Fatalln("Erro na crição de tabelas SQL")
		return err
	}

	return nil
}
