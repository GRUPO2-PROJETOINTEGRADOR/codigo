package controllers

import (
	"fmt"
	"net/http"
)

type OrientacaoController struct{}

func (c OrientacaoController) ListarPaginaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Página Orientação Educativa")
}

func (c OrientacaoController) SalvarHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Salvando orientação")
}
