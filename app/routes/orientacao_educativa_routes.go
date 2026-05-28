package routes

import (
	"codigo/app/controllers"
	"net/http"
)

func RotasOrientacao() {
	var orientacaoController controllers.OrientacaoController

	http.HandleFunc("/conservacao/orientacao-educativa", orientacaoController.ListarPaginaHandler)
	http.HandleFunc("/conservacao/orientacao-educativa/salvar", orientacaoController.SalvarHandler)
	//http.HandleFunc("/conservacao/orientacao-educativa/editar", orientacaoController.Editar) FUNÇÃO SERÁ IMPLEMENTADA BREVEMENTE

	http.HandleFunc("/teste", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("rotas funcionando"))
	})
}
