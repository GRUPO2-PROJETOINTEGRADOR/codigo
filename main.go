package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"codigo/app/controllers"
	utils "codigo/app/repository"
	"codigo/app/routes"
	"codigo/app/services"
)

func main() {
	// Initialize static routes (if any)
	routes.Rotas()

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

	// API routes (only orientacoes needed for now)
	// http.Handle("/api/inspecoes", &inspecoesHandler{})
	// http.HandleFunc("/api/inspecoes/criar", func(w http.ResponseWriter, r *http.Request) {
	//     if r.Method == http.MethodPost {
	//         h := &inspecoesHandler{}
	//         h.ServeHTTP(w, r)
	//     }
	// })
	// http.HandleFunc("/api/inspecoes/deletar", inspecaoDeletarHandler)
	// http.HandleFunc("/api/eco/participantes", participantesHandler)
	// http.HandleFunc("/api/eco/participantes/criar", func(w http.ResponseWriter, r *http.Request) {
	//     if r.Method == http.MethodPost {
	//         participantesHandler(w, r)
	//     }
	// })
	// http.HandleFunc("/api/eco/participantes/toggle-status", func(w http.ResponseWriter, r *http.Request) {
	//     if r.Method == http.MethodPatch {
	//         participantesHandler(w, r)
	//     }
	// })
	// http.HandleFunc("/api/eco/residuos", residuosHandler)
	// http.HandleFunc("/api/eco/residuos/lancar", func(w http.ResponseWriter, r *http.Request) {
	//     if r.Method == http.MethodPost {
	//         residuosHandler(w, r)
	//     }
	// })
	// http.HandleFunc("/api/eco/kits", kitsHandler)
	// http.HandleFunc("/api/eco/kits/registrar", func(w http.ResponseWriter, r *http.Request) {
	//     if r.Method == http.MethodPost {
	//         kitsHandler(w, r)
	//     }
	// })

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

	http.HandleFunc("/api/orientacoes", orientacaoController.ListarJSONHandler)
	http.HandleFunc("/api/orientacoes/criar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			orientacaoController.SalvarHandler(w, r)
		}
	})

	// http.HandleFunc("/api/relatorios/seguranca", relatoriosSegurancaHandler)

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
