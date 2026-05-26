package routes

import (
	"codigo/app/controllers"
	"net/http"
)

func Rotas() {
	var orientacaoController controllers.OrientacaoController

	http.HandleFunc("/conservacao/orientacao-educativa", orientacaoController.ListarPaginaHandler)

	http.HandleFunc("/conservacao/orientacao-educativa/salvar", orientacaoController.SalvarHandler)
}
