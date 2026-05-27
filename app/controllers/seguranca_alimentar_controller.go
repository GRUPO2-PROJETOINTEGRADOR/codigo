package controllers

import (
	"fmt"
	"html/template"
	"net/http"
)

type SegurancaAlimentarController struct{}

func (c SegurancaAlimentarController) ListarPaginaHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.ParseFiles("templates/conservacao/relatorio-seguranca-alimentar.html"),
	)

	tmpl.Execute(w, nil)
}

func (c SegurancaAlimentarController) SalvarHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Salvando relatório de Segurança Alimentar")
}

func (c SegurancaAlimentarController) EditarHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Editando inspeção de Segurança Alimentar")
}

func (c SegurancaAlimentarController) ExcluirHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Excluindo inspeção de Segurança Alimentar")
}
