package controllers

import (
	"codigo/app/models"
	utils "codigo/app/repository"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
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

	var pdfURL string

	if err == nil {

		defer file.Close()

		pdfURL, err = utils.UploadPDF(
			file,
			header.Filename,
			header.Size,
		)

		if err != nil {
			http.Error(w, "Erro ao enviar PDF para MinIO", http.StatusInternalServerError)
			return
		}
	}

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
		AnexoTiller:      pdfURL,
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
		"pdf_url": pdfURL,
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

	var auditoria models.SegurancaAlimentar

	err := json.NewDecoder(r.Body).Decode(&auditoria)

	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "ID inválido", http.StatusBadRequest)
			return
		}
		auditoria.ID = id
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
