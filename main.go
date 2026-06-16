package main

//Comentários com auxílio do CHAT GPT
import (
	"encoding/json"
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

	// Endpoint:
	//   GET /api/lojas
	// Fluxo:
	//   1. Valida método HTTP
	//   2. Busca lojas no banco
	//   3. Retorna JSON
	http.HandleFunc("/api/lojas", func(w http.ResponseWriter, r *http.Request) {

		// Permite apenas GET
		if r.Method != http.MethodGet {

			// Retorna erro 405
			http.Error(
				w,
				"Método não permitido",
				http.StatusMethodNotAllowed,
			)

			return
		}

		// Busca lojas no banco
		lojas, err := utils.Read_lojas()

		// Tratamento de erro
		if err != nil {

			jsonError(
				w,
				"Erro ao buscar lojas",
				http.StatusInternalServerError,
			)

			return
		}

		// Retorna lista em JSO
		jsonResponse(w, lojas)
	})

	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	// SERVIDOR DE ARQUIVOS RAIZ

	// Serve arquivos diretamente da raiz do projeto.
	//
	// Exemplo:
	//   localhost:8080/templates/index.html
	// OBS:
	// Em produção isso geralmente NÃO é recomendado,
	// porque expõe arquivos do projeto.

	http.Handle("/", http.FileServer(http.Dir(".")))

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
	log.Println("║   JP Mall — Servidor Go rodando com sucesso!    ║")

	log.Printf(
		"║   Acesse: http://localhost:%s/templates/index.html       ║\n",
		port,
	)

	log.Println("║   API pronta para desenvolvimento.              ║")
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

// =========================================================
// FUNÇÃO AUXILIAR - RESPOSTA JSON
// =========================================================
// Recebe qualquer estrutura e transforma em JSON.
//
// Exemplo:
//
//	jsonResponse(w, usuario)
//
// =========================================================
func jsonResponse(w http.ResponseWriter, data interface{}) {

	// Define header da resposta
	w.Header().Set("Content-Type", "application/json")

	// Converte struct/map/slice para JSON
	json.NewEncoder(w).Encode(data)
}

// =========================================================
// FUNÇÃO AUXILIAR - RESPOSTA DE ERRO JSON
// =========================================================
// Padroniza respostas de erro da API.
//
// Exemplo de retorno:
//
//	{
//	  "error": "Erro ao buscar lojas"
//	}
//
// =========================================================
func jsonError(w http.ResponseWriter, msg string, code int) {

	// Header JSON
	w.Header().Set("Content-Type", "application/json")

	// Status HTTP
	w.WriteHeader(code)

	// Corpo da resposta
	json.NewEncoder(w).Encode(
		map[string]string{
			"error": msg,
		},
	)
}
