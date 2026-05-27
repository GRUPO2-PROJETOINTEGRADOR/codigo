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

	var auditoria models.SegurancaAlimentar

	err := json.NewDecoder(r.Body).Decode(&auditoria)

	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	err = utils.CriarAuditoria(auditoria)

	if err != nil {
		http.Error(w, "Erro ao salvar auditoria", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auditoria)
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
	w.Write([]byte("Editando inspeção de Segurança Alimentar"))
}
