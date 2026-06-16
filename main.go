package main

//Comentários com auxílio do CHAT GPT
import (
	"fmt"
	"log"
	"net/http"
	"os"

	// Repository / camada de acesso ao banco
	utils "codigo/app/repository"

	// Arquivo responsável pelas rotas adicionais
	"codigo/app/routes"
)

func main() {

	// Aqui são carregadas as rotas definidas no package routes.
	routes.Rotas()

	// CONEXÃO COM O BANCO POSTGRESQL
	// utils.Connect():
	//   - abre conexão com PostgreSQL
	//   - inicializa utils.DB
	if err := utils.Connect(); err != nil {
		log.Fatalln("Erro na conexão com o servidor!")
	}

	// Fecha conexão com banco quando aplicação encerrar
	defer utils.DB.Close()

	// Garante que as tabelas necessárias existam.
	if err := utils.Criar_banco(); err != nil {
		log.Fatalln("Erro na criação de tabelas SQL")

	}

	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	http.Handle("/", http.FileServer(http.Dir("templates")))

	// =========================================================
	// PORTA DO SERVIDOR
	// =========================================================
	// Tenta ler variável de ambiente:
	//   SERVER_PORT
	//
	// Caso não exista:
	//   usa porta 8080
	// =========================================================
	port := os.Getenv("SERVER_PORT")

	if port == "" {
		port = "8080"
	}

	// =========================================================
	// LOG VISUAL DE INICIALIZAÇÃO
	// =========================================================
	log.Println("╔══════════════════════════════════════════════════╗")
	log.Println("║   JP Mall — Servidor Go rodando com sucesso!     ║")

	log.Printf(
				"║   Acesse: http://localhost:%s		       ║\n",
		port,
	)

	log.Println("║   API pronta para desenvolvimento.               ║")
	log.Println("╚══════════════════════════════════════════════════╝")

	// =========================================================
	// INICIALIZAÇÃO DO SERVIDOR HTTP
	// =========================================================
	// ListenAndServe:
	//   inicia servidor HTTP
	//
	// ":" + port:
	//   exemplo -> :8080
	//
	// nil:
	//   usa DefaultServeMux do net/http
	// =========================================================
	if err := http.ListenAndServe(":"+port, nil); err != nil {

		// Escreve erro diretamente no stderr
		fmt.Fprintf(
			os.Stderr,
			"Erro ao iniciar servidor: %v\n",
			err,
		)

		os.Exit(1)
	}
}


