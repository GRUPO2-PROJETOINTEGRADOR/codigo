package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"codigo/app/controllers"
	utils "codigo/app/repository"

	//	"codigo/app/routes"
	"codigo/app/services"
)

func main() {
	// Initialize static routes (if any)
	//routes.Rotas()

	// Connect to PostgreSQL
	if err := utils.Connect(); err != nil {
		log.Fatalln("Erro na conexão com o servidor!")
	}
	defer utils.DB.Close()

	// Ensure tables exist
	if err := utils.Criar_banco(); err != nil {
		log.Fatalln("Erro na criação de tabelas SQL")
	}

	// Instantiate service and controller for orientação educativa
	orientacaoService := services.OrientacaoService{Repo: utils.OrientacaoRepository{}}
	orientacaoController := controllers.OrientacaoController{Service: orientacaoService}

	// New handler to serve lojas list
	http.HandleFunc("/api/lojas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		lojas, err := utils.Read_lojas()
		if err != nil {
			jsonError(w, "Erro ao buscar lojas", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, lojas)
	})

	http.HandleFunc("/conservacao/orientacao-educativa/salvar", orientacaoController.SalvarHandler)

	// Static file server
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Server port
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("╔══════════════════════════════════════════════════╗")
	log.Println("║   JP Mall — Servidor Go rodando com sucesso!    ║")
	log.Printf("║   Acesse: http://localhost:%s/templates/index.html       ║\n", port)
	log.Println("║   API pronta para desenvolvimento.              ║")
	log.Println("╚══════════════════════════════════════════════════╝")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao iniciar servidor: %v\n", err)
		os.Exit(1)
	}
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// Handlers are defined later in the file (inspecoesHandler, participantesHandler, etc.)
