package main

import (
	"fmt"      // Importa o pacote "fmt" para saída formatada, como imprimir mensagens no console.
	"log"      // Importa o pacote "log" para registrar mensagens de erro.
	"net/http" // Importa o pacote "net/http" para criar um servidor HTTP.
)

func main() {
	// Cria um servidor de arquivos que serve os arquivos da pasta "./static".
	fileserver := http.FileServer(http.Dir("./static"))

	// Associa o servidor de arquivos à rota raiz ("/").
	http.Handle("/", fileserver)

	// Imprime no console a mensagem indicando que o servidor está rodando na porta 8081.
	fmt.Printf("port running on http://localhost:8081/\n")

	// Inicia o servidor HTTP na porta 8081. Se ocorrer um erro, ele será registrado e o programa será encerrado.
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err) // Registra o erro e encerra o programa.
	}
}
