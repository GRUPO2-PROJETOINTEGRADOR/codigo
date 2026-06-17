package controllers

import (
	"codigo/app/models"
	utils "codigo/app/repository"
	s "codigo/app/services"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type EcoflamboyantController struct{}

func (c *EcoflamboyantController) ListarEcoFlamboyantHandler(w http.ResponseWriter, r *http.Request) {
	pagina := 1
	if p := r.URL.Query().Get("pagina"); p != "" {
		if n, err := strconv.Atoi(p); err == nil && n > 0 {
			pagina = n
		}
	}
	offset := (pagina - 1) * 10

	lojas, err := s.ListarLojas(utils.DB)
	if err != nil {
		log.Printf("Erro ao listar lojas: %v", err)
		http.Error(w, "Erro ao carregar lojas", http.StatusInternalServerError)
		return
	}

	lojasParticipantes, err := s.ObterLojasParticipantes(utils.DB)
	if err != nil {
		log.Printf("Erro ao listar lojas participantes: %v", err)
		http.Error(w, "Erro ao carregar lojas participantes", http.StatusInternalServerError)
		return
	}

	filtroDataInicio := r.URL.Query().Get("filtro_data_inicio")
	filtroDataFim := r.URL.Query().Get("filtro_data_fim")
	filtroLojaID := r.URL.Query().Get("filtro_loja_id")

	nomeLojaFiltrada := ""
	if filtroLojaID != "" {
		for _, l := range lojas {
			if l.ID == filtroLojaID {
				nomeLojaFiltrada = l.Nome + " (LUC " + l.ID + ")"
				break
			}
		}
	}

	aba := r.URL.Query().Get("aba")
	if aba == "" {
		aba = "lojas"
	}

	var (
		participantes               []models.Participante
		residuos                    []models.Residuo
		kits                        []models.Kit
		registros                   []models.RegistroAuditoria
		totalResiduos               int
		totalPaginasLojas           int
		totalPaginasResiduos        int
		totalPaginasKits            int
		totalPaginasRegistros       int
	)

	switch aba {
	case "lojas":
		participantes, _ = s.ListarParticipantes(utils.DB, filtroDataInicio, filtroDataFim, 10, offset)
		residuos, _ = s.ObterResiduos(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID, 0, 0)
		kits, _ = s.ObterKits(utils.DB, "", "", "", 0, 0)
		registros, _ = s.ListarAuditorias(utils.DB, "", "", "", 0, 0)
		total, _ := s.ContarParticipantes(utils.DB, filtroDataInicio, filtroDataFim)
		totalPaginasLojas = int(math.Ceil(float64(total) / 10.0))
	case "residuos":
		participantes, _ = s.ListarParticipantes(utils.DB, "", "", 0, 0)
		residuos, _ = s.ObterResiduos(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID, 10, offset)
		kits, _ = s.ObterKits(utils.DB, "", "", "", 0, 0)
		registros, _ = s.ListarAuditorias(utils.DB, "", "", "", 0, 0)
		total, _ := s.ContarResiduos(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID)
		totalResiduos = total
		totalPaginasResiduos = int(math.Ceil(float64(total) / 10.0))
	case "kits":
		participantes, _ = s.ListarParticipantes(utils.DB, filtroDataInicio, filtroDataFim, 0, 0)
		residuos, _ = s.ObterResiduos(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID, 0, 0)
		kits, _ = s.ObterKits(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID, 10, offset)
		registros, _ = s.ListarAuditorias(utils.DB, "", "", "", 0, 0)
		total, _ := s.ContarKits(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID)
		totalPaginasKits = int(math.Ceil(float64(total) / 10.0))
	case "registros":
		participantes, _ = s.ListarParticipantes(utils.DB, filtroDataInicio, filtroDataFim, 0, 0)
		residuos, _ = s.ObterResiduos(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID, 0, 0)
		kits, _ = s.ObterKits(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID, 0, 0)
		registros, _ = s.ListarAuditorias(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID, 10, offset)
		total, _ := s.ContarAuditorias(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID)
		totalPaginasRegistros = int(math.Ceil(float64(total) / 10.0))
	}

	totalKits, err := s.ObterTotalKits(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID)
	if err != nil {
		log.Printf("Erro ao somar kits: %v", err)
	}

	totalLojasParticipantes, crescimentoLojas, err := s.ObterDadosLojas(utils.DB, filtroDataInicio, filtroDataFim)
	if err != nil {
		log.Printf("Erro ao obter dados de lojas: %v", err)
	}

	fluxoKits, err := s.ObterFluxoKits(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID)
	if err != nil {
		log.Printf("Erro ao obter fluxo de kits: %v", err)
	}

	volumeTotalGeral, totalAdubo, totalDescartado, taxaAproveitamento, fluxoResiduos, err := s.ObterResumoResiduos(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID)
	if err != nil {
		log.Printf("Erro ao obter resumo de resíduos: %v", err)
	}

	data := models.EcoFlamboyantPageData{
		Participantes:           participantes,
		Lojas:                   lojas,
		Residuos:                residuos,
		TotalResiduos:           totalResiduos,
		Kits:                    kits,
		TotalKits:               totalKits,
		TotalLojasParticipantes: totalLojasParticipantes,
		CrescimentoLojas:        crescimentoLojas,
		FluxoKits:               fluxoKits,
		VolumeTotalGeral:        volumeTotalGeral,
		TotalAdubo:              totalAdubo,
		TotalDescartado:         totalDescartado,
		TaxaAproveitamento:      taxaAproveitamento,
		FluxoResiduos:           fluxoResiduos,
		Registros:               registros,
		AbaAtiva:                aba,
		TodasLojas:              lojasParticipantes,
		FiltroDataInicio:        filtroDataInicio,
		FiltroDataFim:           filtroDataFim,
		FiltroLojaID:            filtroLojaID,
		NomeLojaFiltrada:        nomeLojaFiltrada,
		HojeStr:                 time.Now().Format("2006-01-02"),
		PaginaAtual:             pagina,
		TotalPaginasLojas:       totalPaginasLojas,
		TotalPaginasResiduos:    totalPaginasResiduos,
		TotalPaginasKits:        totalPaginasKits,
		TotalPaginasRegistros:   totalPaginasRegistros,
		ItensPorPagina:          10,
	}

	tmpl := tmplEcoFlamboyant()
	tmpl.ExecuteTemplate(w, "eco-flamboyant", data)
}

func tmplEcoFlamboyant() *template.Template {
	return template.Must(template.New("eco-flamboyant.html").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}).ParseFiles("templates/conservacao/eco-flamboyant.html"))
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

	dados, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Erro ao ler arquivo", http.StatusInternalServerError)
		return
	}

	err = s.CriarParticipante(utils.DB, lojaID, dataEntrada, dataSaida, header.Filename, dados)
	if err != nil {
		log.Printf("Erro ao criar participante: %v", err)
		msg := err.Error()
		if strings.Contains(msg, "23505") || strings.Contains(msg, "unicidade") || strings.Contains(msg, "unique") {
			msg = "Esta loja já está cadastrada como participante do Eco Flamboyant."
		}
		pageData := c.montarPaginaErro(r, "lojas", time.Now().Format("2006-01-02"))
		pageData.ErroForm = msg
		tmpl := tmplEcoFlamboyant()
		tmpl.ExecuteTemplate(w, "eco-flamboyant", pageData)
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
		msg := err.Error()
		if strings.Contains(msg, "23505") || strings.Contains(msg, "unicidade") || strings.Contains(msg, "unique") {
			msg = "Este resíduo já foi registrado."
		}
		pageData := c.montarPaginaErro(r, "residuos", time.Now().Format("2006-01-02"))
		pageData.ErroForm = msg
		tmpl := tmplEcoFlamboyant()
		tmpl.ExecuteTemplate(w, "eco-flamboyant", pageData)
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
		msg := err.Error()
		if strings.Contains(msg, "23505") || strings.Contains(msg, "unicidade") || strings.Contains(msg, "unique") {
			msg = "Este kit já foi registrado."
		}
		pageData := c.montarPaginaErro(r, "kits", time.Now().Format("2006-01-02"))
		pageData.ErroForm = msg
		tmpl := tmplEcoFlamboyant()
		tmpl.ExecuteTemplate(w, "eco-flamboyant", pageData)
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

func (c *EcoflamboyantController) EditarParticipanteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}
	lojaID := r.FormValue("loja_id")
	dataEntrada := r.FormValue("data_entrada")
	dataSaida := r.FormValue("data_saida")

	err := s.AtualizarParticipante(utils.DB, lojaID, dataEntrada, dataSaida)
	if err != nil {
		log.Printf("Erro ao editar participante: %v", err)
		pageData := c.montarPaginaErro(r, "lojas", time.Now().Format("2006-01-02"))
		pageData.ErroForm = err.Error()
		tmpl := tmplEcoFlamboyant()
		tmpl.ExecuteTemplate(w, "eco-flamboyant", pageData)
		return
	}
	http.Redirect(w, r, "/conservacao/eco-flamboyant?aba=lojas", http.StatusSeeOther)
}

func (c *EcoflamboyantController) RemoverParticipanteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}
	lojaID := r.FormValue("loja_id")
	if err := s.InativarParticipante(utils.DB, lojaID); err != nil {
		log.Printf("Erro ao remover participante: %v", err)
		pageData := c.montarPaginaErro(r, "lojas", time.Now().Format("2006-01-02"))
		pageData.ErroForm = err.Error()
		tmpl := tmplEcoFlamboyant()
		tmpl.ExecuteTemplate(w, "eco-flamboyant", pageData)
		return
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

func (c *EcoflamboyantController) BuscarLojasDisponiveis(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	lojas, err := s.ListarLojasBusca(utils.DB, q)
	if err != nil {
		log.Printf("Erro ao buscar lojas: %v", err)
		json.NewEncoder(w).Encode([]models.LojaBusca{})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lojas)
}

func (c *EcoflamboyantController) EmitirRelatorioPDF(w http.ResponseWriter, r *http.Request) {
	dataInicio := r.URL.Query().Get("filtro_data_inicio")
	dataFim := r.URL.Query().Get("filtro_data_fim")
	lojaID := r.URL.Query().Get("filtro_loja_id")

	pdfBytes, err := s.GerarRelatorioPDF(utils.DB, dataInicio, dataFim, lojaID)
	if err != nil {
		log.Printf("Erro ao gerar PDF: %v", err)
		http.Error(w, "Erro ao gerar relatório", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=relatorio-eco-flamboyant.pdf")
	w.Write(pdfBytes)
}

func (c *EcoflamboyantController) montarPaginaErro(r *http.Request, abaAtiva, hojeStr string) models.EcoFlamboyantPageData {
	participantes, err := s.ListarParticipantes(utils.DB, "", "", 0, 0)
	if err != nil {
		log.Printf("Erro ao listar participantes: %v", err)
	}
	lojas, err := s.ListarLojas(utils.DB)
	if err != nil {
		log.Printf("Erro ao listar lojas: %v", err)
	}
	lojasParticipantes, err := s.ObterLojasParticipantes(utils.DB)
	if err != nil {
		log.Printf("Erro ao listar lojas participantes: %v", err)
	}

	filtroDataInicio := r.URL.Query().Get("filtro_data_inicio")
	filtroDataFim := r.URL.Query().Get("filtro_data_fim")
	filtroLojaID := r.URL.Query().Get("filtro_loja_id")

	nomeLojaFiltrada := ""
	if filtroLojaID != "" {
		for _, l := range lojas {
			if l.ID == filtroLojaID {
				nomeLojaFiltrada = l.Nome + " (LUC " + l.ID + ")"
				break
			}
		}
	}

	residuos, err := s.ObterResiduos(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID, 0, 0)
	if err != nil {
		log.Printf("Erro ao listar resíduos: %v", err)
	}
	kits, err := s.ObterKits(utils.DB, "", "", "", 0, 0)
	if err != nil {
		log.Printf("Erro ao listar kits: %v", err)
	}
	totalKits, err := s.ObterTotalKits(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID)
	if err != nil {
		log.Printf("Erro ao somar kits: %v", err)
	}
	totalLojasParticipantes, crescimentoLojas, err := s.ObterDadosLojas(utils.DB, filtroDataInicio, filtroDataFim)
	if err != nil {
		log.Printf("Erro ao obter dados de lojas: %v", err)
	}
	fluxoKits, err := s.ObterFluxoKits(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID)
	if err != nil {
		log.Printf("Erro ao obter fluxo de kits: %v", err)
	}
	volumeTotalGeral, totalAdubo, totalDescartado, taxaAproveitamento, fluxoResiduos, _ := s.ObterResumoResiduos(utils.DB, filtroDataInicio, filtroDataFim, filtroLojaID)
	registros, _ := s.ListarAuditorias(utils.DB, "", "", "", 0, 0)

	return models.EcoFlamboyantPageData{
		Participantes:           participantes,
		Lojas:                   lojas,
		Residuos:                residuos,
		TotalResiduos:           len(residuos),
		Kits:                    kits,
		TotalKits:               totalKits,
		TotalLojasParticipantes: totalLojasParticipantes,
		CrescimentoLojas:        crescimentoLojas,
		FluxoKits:               fluxoKits,
		VolumeTotalGeral:        volumeTotalGeral,
		TotalAdubo:              totalAdubo,
		TotalDescartado:         totalDescartado,
		TaxaAproveitamento:      taxaAproveitamento,
		FluxoResiduos:           fluxoResiduos,
		Registros:               registros,
		AbaAtiva:                abaAtiva,
		TodasLojas:              lojasParticipantes,
		FiltroDataInicio:        filtroDataInicio,
		FiltroDataFim:           filtroDataFim,
		FiltroLojaID:            filtroLojaID,
		NomeLojaFiltrada:        nomeLojaFiltrada,
		HojeStr:                 hojeStr,
	}
}
