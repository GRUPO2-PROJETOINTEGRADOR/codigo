package controllers

import (
	"net/http"
	"os"
)

type DashboardController struct{}

func (c DashboardController) ListarPaginaHandler(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("templates/conservacao/dashboard.html")
	if err != nil {
		http.Error(w, "Página não encontrada", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}
