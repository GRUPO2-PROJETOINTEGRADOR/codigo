package controllers

import (
	"fmt"
	"net/http"
)

type SegurancaAlimentarController struct{}

func (c SegurancaAlimentarController) ListarPaginaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Página Segurança Alimentar")
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
