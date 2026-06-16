package controllers

import (
	"codigo/app/models"
	utils "codigo/app/repository"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type SegurancaAlimentarController struct{}

func (c SegurancaAlimentarController) ListarPaginaHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.ParseFiles("templates/conservacao/relatorio-seguranca-alimentar.html"),
	)

	tmpl.Execute(w, nil)
}

func (c SegurancaAlimentarController) ListarHandler(w http.ResponseWriter, r *http.Request) {

	auditorias, err := utils.ListarAuditorias()

	if err != nil {
		http.Error(w, "Erro ao listar auditorias", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auditorias)
}

func (c SegurancaAlimentarController) SalvarHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// aceita multipart/form-data
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	// =========================
	// PDF
	// =========================
	file, header, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "O arquivo PDF é obrigatório", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validação de tamanho (máximo 5MB)
	if header.Size > 5*1024*1024 {
		http.Error(w, "O tamanho do arquivo PDF excede o limite de 5MB", http.StatusBadRequest)
		return
	}

	// Validação de tipo MIME / Extensão
	contentType := header.Header.Get("Content-Type")
	if contentType != "application/pdf" && !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
		http.Error(w, "Apenas arquivos PDF são permitidos", http.StatusBadRequest)
		return
	}

	// Ler bytes
	pdfBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Erro ao ler o arquivo PDF", http.StatusInternalServerError)
		return
	}

	// Validação da assinatura do arquivo PDF (%PDF)
	if len(pdfBytes) > 4 && string(pdfBytes[:4]) != "%PDF" {
		http.Error(w, "O arquivo enviado não é um PDF válido", http.StatusBadRequest)
		return
	}

	pdfNome := header.Filename
	pdfTipo := contentType
	if pdfTipo == "" {
		pdfTipo = "application/pdf"
	}
	pdfTamanho := int64(len(pdfBytes))

	// =========================
	// CAMPOS FORM
	// =========================
	nota, _ := strconv.Atoi(r.FormValue("nota"))

	auditoria := models.SegurancaAlimentar{
		LojaID:           r.FormValue("loja_id"),
		DataAuditoria:    r.FormValue("data_auditoria"),
		ResponsavelLoja:  r.FormValue("responsavel_loja"),
		CargoResponsavel: r.FormValue("cargo_responsavel"),
		Nota:             nota,
		Classificacao:    r.FormValue("classificacao"),
		AnexoTiller:      pdfNome,
		TipoInspecao:     r.FormValue("tipo_inspecao"),
		NCGrave:          r.FormValue("nc_grave") == "true",
		PDFNome:          pdfNome,
		PDFTipo:          pdfTipo,
		PDFTamanho:       pdfTamanho,
		PDFArquivo:       pdfBytes,
	}

	// =========================
	// SALVA NO POSTGRES
	// =========================
	err = utils.CriarAuditoria(auditoria)

	if err != nil {
		http.Error(w, "Erro ao salvar auditoria", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"pdf_url": pdfNome,
	})
}

func (c SegurancaAlimentarController) ExcluirHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "ID não informado", http.StatusBadRequest)
		return
	}

	auditoriaID, err := strconv.Atoi(id)

	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Deletar auditoria do banco (deleta registro e o binário junto)
	err = utils.DeletarAuditoria(auditoriaID)

	if err != nil {
		http.Error(w, "Erro ao deletar auditoria", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (c SegurancaAlimentarController) EditarHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID não informado", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Erro ao processar formulário", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("pdf")
	var pdfNome string
	var pdfTipo string
	var pdfTamanho int64
	var pdfArquivo []byte
	var pdfURL string

	if err == nil {
		defer file.Close()

		// Validação de tamanho (máximo 5MB)
		if header.Size > 5*1024*1024 {
			http.Error(w, "O tamanho do arquivo PDF excede o limite de 5MB", http.StatusBadRequest)
			return
		}

		// Validação de tipo MIME / Extensão
		contentType := header.Header.Get("Content-Type")
		if contentType != "application/pdf" && !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
			http.Error(w, "Apenas arquivos PDF são permitidos", http.StatusBadRequest)
			return
		}

		// Ler bytes
		pdfBytes, readErr := io.ReadAll(file)
		if readErr != nil {
			http.Error(w, "Erro ao ler o arquivo PDF", http.StatusInternalServerError)
			return
		}

		// Validação da assinatura do arquivo PDF (%PDF)
		if len(pdfBytes) > 4 && string(pdfBytes[:4]) != "%PDF" {
			http.Error(w, "O arquivo enviado não é um PDF válido", http.StatusBadRequest)
			return
		}

		pdfNome = header.Filename
		pdfTipo = contentType
		if pdfTipo == "" {
			pdfTipo = "application/pdf"
		}
		pdfTamanho = int64(len(pdfBytes))
		pdfArquivo = pdfBytes
		pdfURL = pdfNome
	} else {
		pdfURL = r.FormValue("anexo_tiller")
	}

	nota, _ := strconv.Atoi(r.FormValue("nota"))
	ncGrave := r.FormValue("nc_grave") == "true"
	classificacao := r.FormValue("classificacao")
	if ncGrave {
		classificacao = "Crítica"
	}

	tipoInspecao := r.FormValue("tipo_inspecao")
	if tipoInspecao == "" {
		tipoInspecao = r.FormValue("tipo")
	}

	auditoria := models.SegurancaAlimentar{
		ID:               id,
		LojaID:           r.FormValue("loja_id"),
		DataAuditoria:    r.FormValue("data_auditoria"),
		ResponsavelLoja:  r.FormValue("responsavel_loja"),
		CargoResponsavel: r.FormValue("cargo_responsavel"),
		Nota:             nota,
		Classificacao:    classificacao,
		AnexoTiller:      pdfURL,
		TipoInspecao:     tipoInspecao,
		NCGrave:          ncGrave,
		PDFNome:          pdfNome,
		PDFTipo:          pdfTipo,
		PDFTamanho:       pdfTamanho,
		PDFArquivo:       pdfArquivo,
	}

	err = utils.AtualizarAuditoria(auditoria)
	if err != nil {
		http.Error(w, "Erro ao atualizar auditoria", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (c SegurancaAlimentarController) AbrirPDFHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID não informado", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	pdf, err := utils.BuscarPDFAuditoria(id)
	if err != nil {
		http.Error(w, "PDF não encontrado", http.StatusNotFound)
		return
	}

	// Se pdf_arquivo estiver vazio (como em registros antigos sem dados binários)
	if len(pdf.Arquivo) == 0 {
		http.Error(w, "PDF não disponível no banco de dados", http.StatusNotFound)
		return
	}

	contentType := pdf.Tipo
	if contentType == "" {
		contentType = "application/pdf"
	}

	filename := pdf.Nome
	if filename == "" {
		filename = "laudo.pdf"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
	w.Header().Set("Content-Length", strconv.FormatInt(pdf.Tamanho, 10))

	w.Write(pdf.Arquivo)
}
