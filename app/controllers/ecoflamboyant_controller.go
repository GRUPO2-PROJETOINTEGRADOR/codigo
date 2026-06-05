package controllers

import (
	"codigo/app/models"
	utils "codigo/app/repository"
	s "codigo/app/services"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type EcoflamboyantController struct{}

func (c *EcoflamboyantController) ListarEcoFlamboyantHandler(w http.ResponseWriter, r *http.Request) {
	participantes, err := s.ListarParticipantes(utils.DB)
	if err != nil {
		log.Printf("Erro ao listar participantes: %v", err)
		http.Error(w, "Erro ao carregar participantes", http.StatusInternalServerError)
		return
	}

	lojas, err := s.ListarLojas(utils.DB)
	if err != nil {
		log.Printf("Erro ao listar lojas: %v", err)
		http.Error(w, "Erro ao carregar lojas", http.StatusInternalServerError)
		return
	}

	residuos, err := s.ObterResiduos(utils.DB)
	if err != nil {
		log.Printf("Erro ao listar resíduos: %v", err)
		http.Error(w, "Erro ao carregar resíduos", http.StatusInternalServerError)
		return
	}

	kits, err := s.ObterKits(utils.DB)
	if err != nil {
		log.Printf("Erro ao listar kits: %v", err)
		http.Error(w, "Erro ao carregar kits", http.StatusInternalServerError)
		return
	}

	totalKits, err := s.ObterTotalKits(utils.DB)
	if err != nil {
		log.Printf("Erro ao somar kits: %v", err)
		http.Error(w, "Erro ao carregar total de kits", http.StatusInternalServerError)
		return
	}

	totalLojasParticipantes, crescimentoLojas, err := s.ObterDadosLojas(utils.DB)
	if err != nil {
		log.Printf("Erro ao obter dados de lojas: %v", err)
		http.Error(w, "Erro ao carregar dados de lojas", http.StatusInternalServerError)
		return
	}

	fluxoKits, err := s.ObterFluxoKits(utils.DB)
	if err != nil {
		log.Printf("Erro ao obter fluxo de kits: %v", err)
		http.Error(w, "Erro ao carregar fluxo de kits", http.StatusInternalServerError)
		return
	}

	volumeTotalGeral, totalAdubo, totalDescartado, taxaAproveitamento, fluxoResiduos, err := s.ObterResumoResiduos(utils.DB)
	if err != nil {
		log.Printf("Erro ao obter resumo de resíduos: %v", err)
	}

	data := models.EcoFlamboyantPageData{
		Participantes:          participantes,
		Lojas:                  lojas,
		Residuos:               residuos,
		TotalResiduos:          len(residuos),
		Kits:                   kits,
		TotalKits:              totalKits,
		TotalLojasParticipantes: totalLojasParticipantes,
		CrescimentoLojas:        crescimentoLojas,
		FluxoKits:              fluxoKits,
		VolumeTotalGeral:       volumeTotalGeral,
		TotalAdubo:             totalAdubo,
		TotalDescartado:        totalDescartado,
		TaxaAproveitamento:     taxaAproveitamento,
		FluxoResiduos:          fluxoResiduos,
	}

	tmpl := template.Must(template.ParseFiles("templates/conservacao/eco-flamboyant.html"))
	tmpl.ExecuteTemplate(w, "eco-flamboyant", data)
}

func (c *EcoflamboyantController) CriarParticipanteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	lojaID := r.FormValue("loja_id")
	dataEntradaStr := r.FormValue("data_entrada")
	dataSaidaStr := r.FormValue("data_saida")

	dataEntrada, err := time.Parse("2006-01-02", dataEntradaStr)
	if err != nil {
		http.Error(w, "Data de entrada inválida", http.StatusBadRequest)
		return
	}

	var dataSaida *time.Time
	if dataSaidaStr != "" {
		parsed, err := time.Parse("2006-01-02", dataSaidaStr)
		if err != nil {
			http.Error(w, "Data de saída inválida", http.StatusBadRequest)
			return
		}
		dataSaida = &parsed
	}

	file, header, err := r.FormFile("anexo_eco")
	if err != nil {
		http.Error(w, "Termo de aceite obrigatório", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(header.Filename[strings.LastIndex(header.Filename, ".")+1:])
	if ext != "pdf" {
		http.Error(w, "Apenas arquivos PDF são aceitos", http.StatusBadRequest)
		return
	}

	dados, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Erro ao ler arquivo", http.StatusInternalServerError)
		return
	}

	err = s.CriarParticipante(utils.DB, lojaID, dataEntrada, dataSaida, header.Filename, dados)
	if err != nil {
		log.Printf("Erro ao criar participante: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/conservacao/eco-flamboyant", http.StatusSeeOther)
}

func (c *EcoflamboyantController) CriarResiduoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	err := s.CriarResiduo(utils.DB,
		r.FormValue("loja_id"),
		r.FormValue("data_coleta"),
		r.FormValue("peso_kg"),
		r.FormValue("aproveitado"))
	if err != nil {
		log.Printf("Erro ao criar resíduo: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/conservacao/eco-flamboyant?aba=residuos", http.StatusSeeOther)
}

func (c *EcoflamboyantController) CriarKitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	err := s.CriarKit(utils.DB,
		r.FormValue("loja_id"),
		r.FormValue("data_entrega_kit"),
		r.FormValue("qnt_kit"))
	if err != nil {
		log.Printf("Erro ao criar kit: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/conservacao/eco-flamboyant?aba=kits", http.StatusSeeOther)
}

func (c *EcoflamboyantController) AlterarStatusLoja(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	lojaID := r.FormValue("loja_id")
	acao := r.FormValue("acao")
	switch acao {
	case "inativar":
		s.InativarLoja(utils.DB, lojaID)
	case "ativar":
		s.AtivarLoja(utils.DB, lojaID)
	}
	http.Redirect(w, r, "/conservacao/eco-flamboyant?aba=lojas", http.StatusSeeOther)
}

func (c *EcoflamboyantController) DownloadTermo(w http.ResponseWriter, r *http.Request) {
	lojaID := strings.TrimPrefix(r.URL.Path, "/conservacao/eco-flamboyant/termo/")
	nome, dados, err := utils.BuscarTermoPorLoja(utils.DB, lojaID)
	if err != nil || dados == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=\""+nome+"\"")
	w.Write(dados)
}
