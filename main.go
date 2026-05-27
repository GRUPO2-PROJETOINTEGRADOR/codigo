// ATENÇÃO: Este main.go é apenas para testes locais da interface.
// O backend oficial em produção está em codigo/.
// As rotas, handlers e mock data aqui servem exclusivamente para
// validar o frontend durante o desenvolvimento.
package main

import (
	utils "codigo/app/repository"
	"codigo/app/routes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// ============================================================================
// MAIN - Servidor HTTP para servir as páginas estáticas da pasta NOVO.
//
// Este servidor contém rotas de API mockadas (para testes locais) para desenvolvimento frontend.
// Quando for integrar com PostgreSQL, substitua os dados em memória por
// consultas ao banco usando o padrão:
//   routes → controllers → services → repositories
//
// Referência: codigo/guia.md
// ============================================================================

// --- Mock data structures (substituir por models do banco) ---

type Loja struct {
	ID        string `json:"id"`
	Nome      string `json:"nome"`
	Categoria string `json:"categoria"`
}

type Inspecao struct {
	ID            int64  `json:"id"`
	LojaID        string `json:"loja_id"`
	Tipo          string `json:"tipo"`
	DataAuditoria string `json:"data_auditoria"`
	Nota          int    `json:"nota"`
	Responsavel   string `json:"responsavel"`
	Funcao        string `json:"funcao"`
	Observacoes   string `json:"observacoes"`
	AnexoURL      string `json:"anexo_url"`
}

type Participante struct {
	ID         int64  `json:"id"`
	LojaID     string `json:"lojaId"`
	LojaNome   string `json:"lojaNome"`
	DataInicio string `json:"dataInicio"`
	DataSaida  string `json:"dataSaida,omitempty"`
	Status     string `json:"status"`
}

type Residuo struct {
	ID         int64   `json:"id"`
	LojaID     string  `json:"loja_id"`
	LojaNome   string  `json:"lojaNome"`
	Data       string  `json:"data"`
	KgTotal    float64 `json:"kg_total"`
	Aproveitou bool    `json:"aproveitou"`
}

type Kit struct {
	ID         int64  `json:"id"`
	LojaID     string `json:"loja_id"`
	LojaNome   string `json:"lojaNome"`
	Data       string `json:"data"`
	Quantidade int    `json:"quantidade_kit"`
}

type Orientacao struct {
	ID                  int64  `json:"id"`
	LojaID              string `json:"loja_id"`
	ResponsavelPresente string `json:"responsavel_presente"`
	FuncaoResponsavel   string `json:"funcao_responsavel"`
	DataOrientacao      string `json:"data_orientacao"`
	Observacoes         string `json:"observacoes"`
}

// --- In-memory stores (substituir por PostgreSQL) ---

var (
	lojas = []Loja{
		{ID: "QS-03", Nome: "CASA BAUDUCCO", Categoria: "Gastronomia"},
		{ID: "AE-02", Nome: "OUTBACK STEAKHOUSE", Categoria: "Gastronomia"},
		{ID: "AE-05", Nome: "COCO BAMBU", Categoria: "Gastronomia"},
		{ID: "T-104", Nome: "BURGER KING", Categoria: "Gastronomia"},
		{ID: "T-125", Nome: "SUBWAY", Categoria: "Gastronomia"},
		{ID: "101", Nome: "Restaurante Alpha", Categoria: "Gastronomia"},
		{ID: "102", Nome: "Lanchonete Beta", Categoria: "Gastronomia"},
		{ID: "103", Nome: "Pizzaria Gamma", Categoria: "Gastronomia"},
		{ID: "104", Nome: "Sushi Express", Categoria: "Gastronomia"},
		{ID: "105", Nome: "Churrascaria do Mall", Categoria: "Gastronomia"},
		{ID: "201", Nome: "Loja Moda Zeta", Categoria: "Moda"},
	}
	inspecoes = []Inspecao{
		{ID: 1, LojaID: "101", Tipo: "Rotina", DataAuditoria: "2026-05-15", Nota: 85, Responsavel: "Auditor João", Funcao: "Veterinário", AnexoURL: "relatorio_maio_alpha.pdf"},
		{ID: 2, LojaID: "102", Tipo: "Controle", DataAuditoria: "2026-05-10", Nota: 62, Responsavel: "Auditora Maria", Funcao: "Fiscal", AnexoURL: "relatorio_maio_beta.pdf"},
		{ID: 3, LojaID: "103", Tipo: "Controle", DataAuditoria: "2026-05-08", Nota: 45, Responsavel: "Auditor Carlos", Funcao: "Técnico", AnexoURL: "relatorio_maio_gamma.pdf"},
		{ID: 4, LojaID: "104", Tipo: "Rotina", DataAuditoria: "2026-05-20", Nota: 91, Responsavel: "Auditora Ana", Funcao: "Veterinário", AnexoURL: "relatorio_maio_delta.pdf"},
		{ID: 5, LojaID: "105", Tipo: "Extraordinária", DataAuditoria: "2026-05-22", Nota: 28, Responsavel: "Auditor Pedro", Funcao: "Fiscal", AnexoURL: "relatorio_maio_epsilon.pdf"},
	}
	participantes = []Participante{
		{ID: 1, LojaID: "101", LojaNome: "Restaurante Alpha", DataInicio: "2026-01-10", Status: "Ativo"},
		{ID: 2, LojaID: "102", LojaNome: "Lanchonete Beta", DataInicio: "2026-02-15", Status: "Ativo"},
		{ID: 3, LojaID: "103", LojaNome: "Pizzaria Gamma", DataInicio: "2026-03-01", Status: "Ativo"},
	}
	residuos = []Residuo{
		{ID: 1, LojaID: "101", LojaNome: "Restaurante Alpha", Data: "2026-05-20", KgTotal: 15.5, Aproveitou: true},
		{ID: 2, LojaID: "102", LojaNome: "Lanchonete Beta", Data: "2026-05-21", KgTotal: 8.2, Aproveitou: true},
		{ID: 3, LojaID: "103", LojaNome: "Pizzaria Gamma", Data: "2026-05-22", KgTotal: 12.0, Aproveitou: false},
	}
	kits = []Kit{
		{ID: 1, LojaID: "101", LojaNome: "Restaurante Alpha", Data: "2026-05-22", Quantidade: 2},
		{ID: 2, LojaID: "102", LojaNome: "Lanchonete Beta", Data: "2026-05-23", Quantidade: 1},
	}
	orientacoes = []Orientacao{
		{ID: 1, LojaID: "101", ResponsavelPresente: "Maria Souza", FuncaoResponsavel: "Gerente", DataOrientacao: "2026-05-18", Observacoes: "Instruções repassadas sobre horário de descarte do lixo orgânico na doca principal."},
	}
	nextID   int64 = 100
	storesMu sync.RWMutex
)

func getNextID() int64 {
	storesMu.Lock()
	defer storesMu.Unlock()
	nextID++
	return nextID
}

// --- Helpers ---

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// --- Handlers (controllers) ---

func lojasHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, lojas)
}

type inspecoesHandler struct{}

func (h *inspecoesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		jsonResponse(w, inspecoes)
	case "POST":
		var i Inspecao
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			jsonError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		i.ID = getNextID()
		storesMu.Lock()
		inspecoes = append(inspecoes, i)
		storesMu.Unlock()
		jsonResponse(w, i)
	default:
		jsonError(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func inspecaoDeletarHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		jsonError(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		jsonError(w, "ID inválido", http.StatusBadRequest)
		return
	}
	storesMu.Lock()
	defer storesMu.Unlock()
	for i, insp := range inspecoes {
		if insp.ID == id {
			inspecoes = append(inspecoes[:i], inspecoes[i+1:]...)
			jsonResponse(w, map[string]string{"status": "ok"})
			return
		}
	}
	jsonError(w, "Inspeção não encontrada", http.StatusNotFound)
}

func participantesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		jsonResponse(w, participantes)
	case "POST":
		var p Participante
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			jsonError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		p.ID = getNextID()
		if p.Status == "" {
			p.Status = "Ativo"
		}
		// Lookup loja name
		for _, l := range lojas {
			if l.ID == p.LojaID {
				p.LojaNome = l.Nome
				break
			}
		}
		storesMu.Lock()
		participantes = append(participantes, p)
		storesMu.Unlock()
		jsonResponse(w, p)
	case "PATCH":
		idStr := r.URL.Query().Get("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			jsonError(w, "ID inválido", http.StatusBadRequest)
			return
		}
		storesMu.Lock()
		for i, p := range participantes {
			if p.ID == id {
				if p.Status == "Ativo" {
					participantes[i].Status = "Inativo"
				} else {
					participantes[i].Status = "Ativo"
				}
				jsonResponse(w, participantes[i])
				storesMu.Unlock()
				return
			}
		}
		storesMu.Unlock()
		jsonError(w, "Participante não encontrado", http.StatusNotFound)
	default:
		jsonError(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func residuosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		jsonResponse(w, residuos)
	case "POST":
		var res Residuo
		if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
			jsonError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		res.ID = getNextID()
		for _, l := range lojas {
			if l.ID == res.LojaID {
				res.LojaNome = l.Nome
				break
			}
		}
		storesMu.Lock()
		residuos = append(residuos, res)
		storesMu.Unlock()
		jsonResponse(w, res)
	default:
		jsonError(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func kitsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		jsonResponse(w, kits)
	case "POST":
		var k Kit
		if err := json.NewDecoder(r.Body).Decode(&k); err != nil {
			jsonError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		k.ID = getNextID()
		for _, l := range lojas {
			if l.ID == k.LojaID {
				k.LojaNome = l.Nome
				break
			}
		}
		storesMu.Lock()
		kits = append(kits, k)
		storesMu.Unlock()
		jsonResponse(w, k)
	default:
		jsonError(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func orientacoesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		jsonResponse(w, orientacoes)
	case "POST":
		var o Orientacao
		if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
			jsonError(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		o.ID = getNextID()
		storesMu.Lock()
		orientacoes = append(orientacoes, o)
		storesMu.Unlock()
		jsonResponse(w, o)
	default:
		jsonError(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func relatoriosSegurancaHandler(w http.ResponseWriter, r *http.Request) {
	startMonth := r.URL.Query().Get("start_month")
	startYear := r.URL.Query().Get("start_year")
	endMonth := r.URL.Query().Get("end_month")
	endYear := r.URL.Query().Get("end_year")

	var filtered []Inspecao

	// Map month name to number
	months := map[string]int{
		"Janeiro": 1, "Fevereiro": 2, "Março": 3, "Abril": 4,
		"Maio": 5, "Junho": 6, "Julho": 7, "Agosto": 8,
		"Setembro": 9, "Outubro": 10, "Novembro": 11, "Dezembro": 12,
	}

	startM := months[startMonth]
	startY, _ := strconv.Atoi(startYear)
	endM := months[endMonth]
	endY, _ := strconv.Atoi(endYear)

	startVal := startY*12 + startM
	endVal := endY*12 + endM

	for _, insp := range inspecoes {
		t, err := time.Parse("2006-01-02", insp.DataAuditoria)
		if err != nil {
			continue
		}
		curVal := t.Year()*12 + int(t.Month())
		if curVal >= startVal && curVal <= endVal {
			filtered = append(filtered, insp)
		}
	}

	jsonResponse(w, filtered)
}

func main() {

	routes.Rotas()

	err := utils.Connect()
	if err != nil {
		log.Fatalln("Erro na conexão com o servidor!")
	}
	defer utils.DB.Close()

	err = utils.Criar_banco()
	if err != nil {
		log.Fatalln("Erro na crição de tabelas SQL")
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// -----------------------------------------------------------------------
	// INTEGRAÇÃO: Descomente o bloco abaixo ao conectar com o Postgres
	// -----------------------------------------------------------------------
	// err := utils.Connect()
	// if err != nil {
	// 	log.Fatalln("Erro na conexão com o servidor!")
	// }
	// defer utils.DB.Close()
	//
	// err = utils.Criar_banco()
	// if err != nil {
	// 	log.Fatalln("Erro na criação de tabelas SQL")
	// }
	// -----------------------------------------------------------------------

	// --- API ROUTES ---
	// INTEGRAÇÃO: Substitua os handlers abaixo por controllers reais
	// que consultam o banco PostgreSQL via repositories.

	http.HandleFunc("/api/lojas", lojasHandler)

	http.HandleFunc("/api/eco/participantes", participantesHandler)
	http.HandleFunc("/api/eco/participantes/criar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			participantesHandler(w, r)
		}
	})
	http.HandleFunc("/api/eco/participantes/toggle-status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			participantesHandler(w, r)
		}
	})

	http.HandleFunc("/api/eco/residuos", residuosHandler)
	http.HandleFunc("/api/eco/residuos/lancar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			residuosHandler(w, r)
		}
	})

	http.HandleFunc("/api/eco/kits", kitsHandler)
	http.HandleFunc("/api/eco/kits/registrar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			kitsHandler(w, r)
		}
	})

	http.HandleFunc("/api/orientacoes", orientacoesHandler)
	http.HandleFunc("/api/orientacoes/criar", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			orientacoesHandler(w, r)
		}
	})

	http.HandleFunc("/api/relatorios/seguranca", relatoriosSegurancaHandler)

	// --- SERVIDOR DE ARQUIVOS ESTÁTICOS ---
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("╔══════════════════════════════════════════════════╗")
	log.Println("║   JP Mall — Servidor Go rodando com sucesso!    ║")
	log.Printf("║   Acesse: http://localhost:%s/templates/index.html       ║\n", port)
	log.Println("║   API mockada pronta para desenvolvimento.      ║")
	log.Println("╚══════════════════════════════════════════════════╝")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao iniciar servidor: %v\n", err)
		os.Exit(1)
	}
}
