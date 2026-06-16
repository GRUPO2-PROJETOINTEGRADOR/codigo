package controllers

import (
	"html/template"
	"net/http"
)

type RelatoriosController struct{}

func (c RelatoriosController) ListarPaginaHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(
		template.ParseFiles("templates/conservacao/relatorios.html"),
	)

	tmpl.Execute(w, nil)
}
